package poplural

import (
	"golang.org/x/text/language"
)

// Language that is used for the rule if no rule can be found for a language.
var fallbackLanguage = "en"

// PluralFunc is a function that returns the appropriate plural form for a value.
type PluralFunc = func(a any) (int, error)

// ForLanguage searches for the appropriate built-in plural function for a language.
// If no function is found, the English plural function is automatically used.
// The second return value indicates whether a suitable function was found or whether the fallback is used.
func ForLanguage(lang language.Tag) (PluralFunc, bool) {
	form, found := pluralRuleForLanguage(lang)
	return form.Evaluate, found
}

func pluralRuleForLanguage(lang language.Tag) (*Rule, bool) {
	n := lang
	for !n.IsRoot() {
		if rule := getBuiltInForLanguage(n.String()); rule != nil {
			return rule, true
		}

		base, confidence := n.Base()
		if confidence >= language.High {
			if rule := getBuiltInForLanguage(base.String()); rule != nil {
				return rule, true
			}
		}

		n = n.Parent()
	}

	return getBuiltInForLanguage(fallbackLanguage), false
}
