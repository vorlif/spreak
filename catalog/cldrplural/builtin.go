package cldrplural

import "golang.org/x/text/language"

// ForLanguage returns the set of rules for a language.
// If no matching language is found, the English rule set and false are returned.
func ForLanguage(lang language.Tag) (*RuleSet, bool) {
	n := lang
	for !n.IsRoot() {
		if ruleSet := forLanguage(n.String()); ruleSet != nil {
			return ruleSet, true
		}

		base, confidence := n.Base()
		if confidence >= language.High {
			if ruleSet := forLanguage(base.String()); ruleSet != nil {
				return ruleSet, true
			}
		}

		n = n.Parent()
	}

	return forLanguage(language.English.String()), false
}
