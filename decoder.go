package spreak

import (
	"fmt"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/internal/mo"
	"github.com/vorlif/spreak/internal/plural"
	"github.com/vorlif/spreak/pkg/po"
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

	return buildGettextCatalog(poFile, lang, domain)
}

func (moDecoder) Decode(lang language.Tag, domain string, data []byte) (Catalog, error) {
	moFile, errParse := mo.ParseBytes(data)
	if errParse != nil {
		return nil, errParse
	}

	// We could check here if the language of the file matches the target language,
	// but leave it off to make loading more flexible.

	return buildGettextCatalog(moFile, lang, domain)
}

func buildGettextCatalog(file *po.File, lang language.Tag, domain string) (Catalog, error) {
	messages := make(messageLookupMap, len(file.Messages))

	for ctx := range file.Messages {
		if len(file.Messages[ctx]) == 0 {
			continue
		}

		if _, hasContext := messages[ctx]; !hasContext {
			messages[ctx] = make(map[string]*gettextMessage)
		}

		for msgID, poMsg := range file.Messages[ctx] {
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

	if file.Header != nil && file.Header.PluralForms != "" {
		forms, err := plural.Parse(file.Header.PluralForms)
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
