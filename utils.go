package spreak

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/text/language"
)

// ExpandLanguage returns possible filenames for a language tag without extension.
func ExpandLanguage(lang language.Tag) []string {
	expansions := make(map[string]bool, 4)
	expansions[lang.String()] = true

	base, baseConf := lang.Base()
	if baseConf > language.No {
		expansions[base.ISO3()] = true
		expansions[base.String()] = true
	}

	region, regionConf := lang.Region()
	if regionConf > language.No && baseConf > language.No {
		key := fmt.Sprintf("%s_%s", base.String(), region.String())
		expansions[key] = true

		key = fmt.Sprintf("%s-%s", base.String(), region.String())
		expansions[key] = true
	}

	script, scriptConf := lang.Script()
	if scriptConf > language.No && baseConf > language.No {
		key := fmt.Sprintf("%s_%s", base.String(), script.String())
		expansions[key] = true

		key = fmt.Sprintf("%s-%s", base.String(), script.String())
		expansions[key] = true
	}

	return stringMapKeys(expansions)
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

func stringMapKeys(m map[string]bool) []string {
	keys := make([]string, len(m))
	i := 0
	for key := range m {
		keys[i] = key
		i++
	}
	// Longest first
	sort.SliceStable(keys, func(i, j int) bool {
		if x, y := len(keys[i]), len(keys[j]); x != y {
			return x > y
		}

		return keys[i] < keys[j]
	})
	return keys
}
