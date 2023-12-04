package cldrplural

import "golang.org/x/text/language"

// Language that is used for the rule if no rule can be found for a language.
var fallbackLanguage = "en"

// ForLanguage returns the set of rules for a language.
// If no matching language is found, the English rule set and false are returned.
func ForLanguage(lang language.Tag) (*RuleSet, bool) {
	n := lang
	for !n.IsRoot() {
		if ruleSet := getBuiltInForLanguage(n.String()); ruleSet != nil {
			return ruleSet, true
		}

		base, confidence := n.Base()
		if confidence >= language.High {
			if ruleSet := getBuiltInForLanguage(base.String()); ruleSet != nil {
				return ruleSet, true
			}
		}

		n = n.Parent()
	}

	return getBuiltInForLanguage(fallbackLanguage), false
}
