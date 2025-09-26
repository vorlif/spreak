package catalog

import (
	"fmt"
	"strings"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/catalog/cldrplural"
	"github.com/vorlif/spreak/catalog/po"
	"github.com/vorlif/spreak/catalog/poplural"
	"github.com/vorlif/spreak/internal/mo"
)

const poCLDRHeader = "X-spreak-use-CLDR"

type poDecoder struct {
	useCLDRPlural bool
}

type moDecoder struct {
	useCLDRPlural bool
}

// NewPoDecoder returns a new Decoder for reading po files.
// If a plural forms header is set, it will be used.
// Otherwise, the CLDR plural rules are used to set the plural form.
// If there is no CLDR plural rule, the English plural rules will be used.
func NewPoDecoder() Decoder { return &poDecoder{} }

// NewPoCLDRDecoder creates a decoder for reading po files,
// which always uses the CLDR plural rules for determining the plural form.
// If no matching CLDR rule exists, the Po header rule is used. If no header exists,
// the english plural rules (1 is singular, otherwise plural) are used.
// Attention: The "Plural-Forms" header inside the Po file is ignored when using the CLDR rules.
// To ensure optimal compatibility with other applications, care should be taken to ensure that the Po header is compatible with the CLDR rules.
func NewPoCLDRDecoder() Decoder { return &poDecoder{useCLDRPlural: true} }

// NewMoDecoder returns a new Decoder for reading mo files.
// If a plural forms header is set, it will be used.
// Otherwise, the CLDR plural rules are used to set the plural form.
// If there is no CLDR plural rule, the English plural rules will be used.
func NewMoDecoder() Decoder { return &moDecoder{useCLDRPlural: false} }

// NewMoCLDRDecoder creates a decoder for reading mo files,
// which always uses the CLDR plural rules for determining the plural form.
// If no matching CLDR rule exists, the Mo header rule is used. If no header exists,
// the english plural rules (1 is singular, otherwise plural) are used.
// Attention: The "Plural-Forms" header inside the Mo file is ignored when using the CLDR rules.
// To ensure optimal compatibility with other applications, care should be taken to ensure that the Mo header is compatible with the CLDR rules.
func NewMoCLDRDecoder() Decoder { return &moDecoder{useCLDRPlural: true} }

func (d poDecoder) Decode(lang language.Tag, domain string, data []byte) (Catalog, error) {
	parser := po.NewParser()
	parser.SetIgnoreComments(true)
	poFile, errParse := parser.Parse(string(data))
	if errParse != nil {
		return nil, errParse
	}

	// We could check here if the language of the file matches the target language,
	// but leave it off to make loading more flexible.

	return buildGettextCatalog(poFile, lang, domain, d.useCLDRPlural)
}

func (d moDecoder) Decode(lang language.Tag, domain string, data []byte) (Catalog, error) {
	moFile, errParse := mo.ParseBytes(data)
	if errParse != nil {
		return nil, errParse
	}

	// We could check here if the language of the file matches the target language,
	// but leave it off to make loading more flexible.

	return buildGettextCatalog(moFile, lang, domain, d.useCLDRPlural)
}

func buildGettextCatalog(file *po.File, lang language.Tag, domain string, useCLDRPlural bool) (Catalog, error) {
	lookupMap := make(PoLookupMap, len(file.Messages))

	for ctx := range file.Messages {
		if len(file.Messages[ctx]) == 0 {
			continue
		}

		if _, hasContext := lookupMap[ctx]; !hasContext {
			lookupMap[ctx] = make(map[string]*GettextMessage, len(file.Messages[ctx]))
		}

		for msgID, poMsg := range file.Messages[ctx] {
			if msgID == "" {
				continue
			}

			if poMsg.Comment != nil && poMsg.Comment.HasFlag("fuzzy") {
				continue
			}

			d := &GettextMessage{
				translations: poMsg.Str,
			}

			lookupMap[poMsg.Context][poMsg.ID] = d
		}
	}

	cat := &GettextCatalog{
		language:  lang,
		lookupMap: lookupMap,
	}

	hasPluralFormHeader := file.Header != nil && file.Header.PluralForms != ""
	if useCLDRPlural || !hasPluralFormHeader {
		cat.pluralFunc = getCLDRPluralFunction(lang)
		return cat, nil
	}

	if val := file.Header.Get(poCLDRHeader); strings.EqualFold(val, "true") {
		cat.pluralFunc = getCLDRPluralFunction(lang)
		return cat, nil
	}

	rule, err := poplural.Parse(file.Header.PluralForms)
	if err != nil {
		return nil, fmt.Errorf("spreak.Decoder: plural rule for po file %v#%v could not be parsed: %w", lang, domain, err)
	}
	cat.pluralFunc = rule.Evaluate
	return cat, nil
}

func getCLDRPluralFunction(lang language.Tag) func(a any) (int, error) {
	ruleSet, _ := cldrplural.ForLanguage(lang)

	catToForm := make(map[cldrplural.Category]int, len(ruleSet.Categories))
	for idx, cat := range ruleSet.Categories {
		catToForm[cat] = idx
	}

	return func(a any) (int, error) {
		cat, err := ruleSet.Evaluate(a)
		if err != nil {
			return 0, err
		}

		if form, ok := catToForm[cat]; ok {
			return form, nil
		}

		return 0, nil
	}
}
