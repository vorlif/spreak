package poplural

import (
	"golang.org/x/text/language"
)

// Plural rule for english texts.
const fallbackRule = "nplurals=2; plural=n != 1;"

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
		if form := forLanguage(n.String()); form != nil {
			return form, true
		}

		base, confidence := n.Base()
		if confidence >= language.High {
			if form := forLanguage(base.String()); form != nil {
				return form, true
			}
		}

		n = n.Parent()
	}

	return forRawRule(fallbackRule), false
}
