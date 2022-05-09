package spreak

import (
	"fmt"
	"strings"

	"golang.org/x/text/language"
)

// ExpandLanguage returns possible filenames for a language tag without extension.
func ExpandLanguage(lang language.Tag) []string {
	expansions := []string{lang.String()}

	base, baseConf := lang.Base()
	if baseConf == language.No {
		return expansions
	}

	region, regionConf := lang.Region()
	if regionConf == language.No {
		return expansions
	}

	expansions = append(expansions,
		fmt.Sprintf("%s_%s", base.String(), region.String()),
		fmt.Sprintf("%s-%s", base.String(), region.String()),
		base.ISO3(),
		fmt.Sprintf("%s_%s", base.ISO3(), region.ISO3()),
		fmt.Sprintf("%s-%s", base.ISO3(), region.ISO3()),
	)

	return expansions
}

func parseLanguageName(lang string) (language.Tag, error) {
	if idx := strings.Index(lang, ":"); idx != -1 {
		lang = lang[:idx]
	}
	if idx := strings.Index(lang, "@"); idx != -1 {
		lang = lang[:idx]
	}
	if idx := strings.Index(lang, "."); idx != -1 {
		lang = lang[:idx]
	}
	lang = strings.TrimSpace(lang)
	return language.Parse(lang)
}

func languageInterfaceToTag(i interface{}) (language.Tag, error) {
	switch v := i.(type) {
	case string:
		tag, err := parseLanguageName(v)
		if err != nil {
			return language.Und, err
		}
		return tag, nil
	case language.Tag:
		return v, nil
	default:
		return language.Und, ErrRequireStringTag
	}
}
