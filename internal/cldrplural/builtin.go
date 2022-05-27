package cldrplural

import "golang.org/x/text/language"

var builtInRuleSets = make(map[string]*RuleSet, 40)

func addRuleSet(langs []string, set *RuleSet) {
	for _, lang := range langs {
		tag := language.MustParse(lang)
		builtInRuleSets[tag.String()] = set
	}
}

func newCategories(categories ...Category) []Category {
	return categories
}
