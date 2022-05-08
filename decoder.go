package spreak

import (
	"fmt"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/internal/plural"

	"github.com/vorlif/spreak/internal/mo"
	"github.com/vorlif/spreak/internal/po"
)

// A Decoder reads and decodes catalogs for a language and a domain from a byte array.
type Decoder interface {
	Decode(lang language.Tag, domain string, data []byte) (Catalog, error)
}

type poDecoder struct{}
type moDecoder struct{}

var _ Decoder = (*poDecoder)(nil)
var _ Decoder = (*moDecoder)(nil)

// NewPoDecoder returns a new Decoder for reading po files.
func NewPoDecoder() Decoder {
	return &poDecoder{}
}

// NewMoDecoder returns a new Decoder for reading mo files.
func NewMoDecoder() Decoder {
	return &moDecoder{}
}

func (poDecoder) Decode(lang language.Tag, domain string, data []byte) (Catalog, error) {
	poFile, errParse := po.Parse(data)
	if errParse != nil {
		return nil, errParse
	}

	// We could check here if the language of the file matches the target language,
	// but leave it off to make loading more flexible.

	messages := make(messageLookupMap, len(poFile.Messages))

	for ctx := range poFile.Messages {
		if len(poFile.Messages[ctx]) == 0 {
			continue
		}

		if _, hasContext := messages[ctx]; !hasContext {
			messages[ctx] = make(map[string]*gettextMessage)
		}

		for msgID, poMsg := range poFile.Messages[ctx] {
			if msgID == "" {
				continue
			}

			d := &gettextMessage{
				Context:      poMsg.Context,
				ID:           poMsg.ID,
				IDPlural:     poMsg.IDPlural,
				Translations: poMsg.Str,
			}

			messages[poMsg.Context][poMsg.ID] = d
		}
	}

	catl := &gettextCatalog{
		language:     lang,
		translations: messages,
	}

	if poFile.Header != nil && poFile.Header.PluralForms != "" {
		forms, err := plural.Parse(poFile.Header.PluralForms)
		if err != nil {
			return nil, fmt.Errorf("spreak.Decoder: plural forms for po file %v#%v could not be parsed: %w", lang, domain, err)
		}
		catl.pluralFunc = forms.IndexForN
	} else {
		forms, _ := plural.ForLanguage(lang)
		catl.pluralFunc = forms
	}

	return catl, nil
}

func (moDecoder) Decode(lang language.Tag, domain string, data []byte) (Catalog, error) {
	moFile, errParse := mo.ParseBytes(data)
	if errParse != nil {
		return nil, errParse
	}

	// We could check here if the language of the file matches the target language,
	// but leave it off to make loading more flexible.

	messages := make([]*gettextMessage, 0, len(moFile.Messages))
	for _, moMsg := range moFile.Messages {
		if moMsg.ID == "" {
			continue
		}

		d := &gettextMessage{
			Context:      moMsg.Context,
			ID:           moMsg.ID,
			IDPlural:     moMsg.IDPlural,
			Translations: make(map[int]string, len(moMsg.StrPlural)),
		}

		// singular translation
		if moMsg.IDPlural == "" {
			d.Translations[0] = moMsg.Str
			continue
		}

		// plural translation
		for idx, t := range moMsg.StrPlural {
			d.Translations[idx] = t
		}

		messages = append(messages, d)
	}

	catl := &gettextCatalog{
		language:     lang,
		translations: transformMessageArray(messages),
		domain:       domain,
	}

	if moFile.Header.PluralForms != "" {
		forms, err := plural.Parse(moFile.Header.PluralForms)
		if err != nil {
			return nil, fmt.Errorf("spreak.Decoder: plural forms for mo file %v#%v could not be parsed: %w", lang, domain, err)
		}
		catl.pluralFunc = forms.IndexForN
	} else {
		forms, _ := plural.ForLanguage(lang)
		catl.pluralFunc = forms
	}

	return catl, nil
}
