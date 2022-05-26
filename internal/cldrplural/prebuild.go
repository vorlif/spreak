package cldrplural

import "golang.org/x/text/language"

var prebuildRuleSets = make(map[language.Tag]*RuleSet)

func addRuleSet(langs []string, set *RuleSet) {
	for _, lang := range langs {
		tag := language.MustParse(lang)
		prebuildRuleSets[tag] = set
	}
}

func newCategories(categories ...Category) []Category {
	return categories
}
