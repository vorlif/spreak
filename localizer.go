package spreak

import (
	"golang.org/x/text/language"

	"github.com/vorlif/spreak/localize"
)

// A Localizer is similar to a Locale except that it is not bound to a language.
// A number of supported languages can be specified at creation time
// where the language matcher of the bundle decides which language fits best.
// For this language the Localizer then offers the possibility to translate.
// If no language fits, the fallback language is used.
// If no fallback language is specified, the source language is used.
type Localizer struct {
	bundle        *Bundle
	locale        *Locale
	defaultDomain string
}

// NewLocalizerForDomain creates a new Localizer for a language and a default domain,
// which is used if no domain is specified.
// Multiple languages can be passed and the best matching language is searched for.
// If no matching language is found, a Localizer is created which returns the original messages.
// Valid languages are strings or language.Tag. All other inputs are dropped.
func NewLocalizerForDomain(bundle *Bundle, domain string, lang ...interface{}) *Localizer {
	var tags []language.Tag

	for _, accept := range lang {
		switch val := accept.(type) {
		case language.Tag:
			tags = append(tags, val)
		case string:
			desired, _, err := language.ParseAcceptLanguage(val)
			if err != nil {
				continue
			}
			tags = append(tags, desired...)
		default:
			continue
		}
	}

	if _, index, conf := bundle.languageMatcher.Match(tags...); conf > language.No {
		tag := bundle.languages[index]
		return newLocalizerFromTag(bundle, domain, tag)
	}

	if bundle.fallbackLanguage != language.Und {
		return newLocalizerFromTag(bundle, domain, bundle.fallbackLanguage)
	} else if bundle.sourceLanguage != language.Und {
		return newLocalizerFromTag(bundle, domain, bundle.sourceLanguage)
	} else {
		return newLocalizerFromTag(bundle, domain)
	}
}

// NewLocalizer operates like NewLocalizerForDomain, with the default domain of the bundle as domain.
func NewLocalizer(bundle *Bundle, langs ...interface{}) *Localizer {
	return NewLocalizerForDomain(bundle, bundle.defaultDomain, langs...)
}

func newLocalizerFromTag(bundle *Bundle, domain string, tag ...language.Tag) *Localizer {
	l := &Localizer{
		bundle:        bundle,
		defaultDomain: domain,
	}

	for _, accept := range tag {
		if locale, err := NewLocale(bundle, accept); err == nil {
			l.locale = locale
			break
		}
	}

	return l
}

// HasLocale returns whether a matching locale has been found and message translation can take place.
func (l *Localizer) HasLocale() bool { return l.locale != nil }

// GetLocale returns the found Locale.
// If no locale was found, nil is returned.
func (l *Localizer) GetLocale() *Locale { return l.locale }

// DefaultDomain returns the default domain.
// The default domain is used if a domain is not explicitly specified for a requested translation.
// If no default domain is specified, the default domain of the bundle is used.
func (l *Localizer) DefaultDomain() string { return l.defaultDomain }

// Language returns the language into which the translation of messages is performed.
// If no language is present, language.Und is returned.
func (l *Localizer) Language() language.Tag {
	if l.locale != nil {
		return l.locale.Language()
	}

	return language.Und
}

// The Get function return the localized translation of message, based on the used Locale
// current default domain and language of the Locale.
// The message argument identifies the message to be translated.
// If no suitable translation exists, the original message is returned.
func (l *Localizer) Get(message localize.Singular) string {
	t, _ := l.dpGettextErr(l.defaultDomain, NoCtx, message)
	return t
}

// Getf operates like Get, but formats the message according to a format identifier and returns the resulting string.
func (l *Localizer) Getf(message localize.Singular, vars ...interface{}) string {
	t, _ := l.dpGettextErr(l.defaultDomain, NoCtx, message, vars...)
	return t
}

// DGet operates like Get, but look the message up in the specified domain.
func (l *Localizer) DGet(domain localize.Domain, message localize.Singular) string {
	t, _ := l.dpGettextErr(domain, NoCtx, message)
	return t
}

// DGetf operates like Get, but look the message up in the specified domain and
// formats the message according to a format identifier and returns the resulting string.
func (l *Localizer) DGetf(domain localize.Domain, message localize.Singular, vars ...interface{}) string {
	t, _ := l.dpGettextErr(domain, NoCtx, message, vars...)
	return t
}

// NGet acts like Get, but consider plural forms.
// The plural formula is applied to n and return the resulting message (some languages have more than two plurals).
func (l *Localizer) NGet(singular localize.Singular, plural localize.Plural, n interface{}) string {
	t, _ := l.dnpGettextErr(l.defaultDomain, NoCtx, singular, plural, n)
	return t
}

// NGetf operates like NGet, but formats the message according to a format identifier and returns the resulting string.
func (l *Localizer) NGetf(singular localize.Singular, plural localize.Plural, n interface{}, vars ...interface{}) string {
	t, _ := l.dnpGettextErr(l.defaultDomain, NoCtx, singular, plural, n, vars...)
	return t
}

// DNGet operates like NGet, but look the message up in the specified domain.
func (l *Localizer) DNGet(domain localize.Domain, singular localize.Singular, plural localize.Plural, n interface{}) string {
	t, _ := l.dnpGettextErr(domain, NoCtx, singular, plural, n)
	return t
}

// DNGetf operates like DNGet, but formats the message according to a format identifier and returns the resulting string.
func (l *Localizer) DNGetf(domain localize.Domain, singular localize.Singular, plural localize.Plural, n interface{}, vars ...interface{}) string {
	t, _ := l.dnpGettextErr(domain, NoCtx, singular, plural, n, vars...)
	return t
}

// PGet operates like Get, but restricted to the specified context.
func (l *Localizer) PGet(context localize.Context, message localize.Singular) string {
	t, _ := l.dpGettextErr(l.defaultDomain, context, message)
	return t
}

// PGetf operates like PGet, but formats the message according to a format identifier and returns the resulting string.
func (l *Localizer) PGetf(context localize.Context, message localize.Singular, vars ...interface{}) string {
	t, _ := l.dpGettextErr(l.defaultDomain, context, message, vars...)
	return t
}

// DPGet operates like Get, but look the message up in the specified domain and with the specified context.
func (l *Localizer) DPGet(domain localize.Domain, context localize.Context, message localize.Singular) string {
	t, _ := l.dpGettextErr(domain, context, message)
	return t
}

// DPGetf operates like DPGet, but formats the message according to a format identifier and returns the resulting string.
func (l *Localizer) DPGetf(domain localize.Domain, context localize.Context, message localize.Singular, vars ...interface{}) string {
	t, _ := l.dpGettextErr(domain, context, message, vars...)
	return t
}

// NPGet operates like NGet, but restricted to the specified context.
func (l *Localizer) NPGet(context localize.Context, singular localize.Singular, plural localize.Plural, n interface{}) string {
	t, _ := l.dnpGettextErr(l.defaultDomain, context, singular, plural, n)
	return t
}

// NPGetf operates like NPGet, but formats the message according to a format identifier and returns the resulting string.
func (l *Localizer) NPGetf(context localize.Context, singular localize.Singular, plural localize.Plural, n interface{}, vars ...interface{}) string {
	t, _ := l.dnpGettextErr(l.defaultDomain, context, singular, plural, n, vars...)
	return t
}

// DNPGet operates like NGet, but look the message up in the specified domain and with the specified context.
func (l *Localizer) DNPGet(domain localize.Domain, context localize.Context, singular localize.Singular, plural localize.Plural, n interface{}) string {
	t, _ := l.dnpGettextErr(domain, context, singular, plural, n)
	return t
}

// DNPGetf operates like DNPGet, but formats the message according to a format identifier and returns the resulting string.
func (l *Localizer) DNPGetf(domain localize.Domain, context localize.Context, singular localize.Singular, plural localize.Plural, n interface{}, vars ...interface{}) string {
	t, _ := l.dnpGettextErr(domain, context, singular, plural, n, vars...)
	return t
}

// LocalizeWithError translates structs that implement the interface localize.Localizable.
// If a suitable translation is found, it will be returned.
// If no matching translation is found, the original string with the matching plural form and an error are returned.
func (l *Localizer) LocalizeWithError(t localize.Localizable) (string, error) {
	if l.locale != nil {
		return l.locale.LocalizeWithError(t)
	}

	var vars []interface{}
	if len(t.GetVars()) > 0 {
		vars = append(vars, t.GetVars()...)
	}

	domain := l.defaultDomain
	if t.HasDomain() {
		domain = t.GetDomain()
	}

	if t.GetPluralID() != "" {
		return l.dnpGettextErr(domain, t.GetContext(), t.GetMsgID(), t.GetPluralID(), t.GetCount(), vars...)
	}

	return l.dpGettextErr(domain, t.GetContext(), t.GetMsgID(), vars...)
}

// Localize acts like LocalizeWithError, but does not return an error.
func (l *Localizer) Localize(t localize.Localizable) string {
	translated, _ := l.LocalizeWithError(t)
	return translated
}

// LocalizeError translates the passed error and returns a new error of type localize.Error
// which wraps the original error.
// If no suitable translation is found, the original error is returned.
// By default, localized messages with the context "errors" are searched for.
// The query is limited to the current domain and the error context specified in the corresponding bundle.
// By default, this is the context "errors".
// Using WithErrorContext("other") during bundle creation to change the error context for a bundle.
func (l *Localizer) LocalizeError(err error) error {
	if l.locale != nil {
		return l.locale.LocalizeError(err)
	}

	return err
}

func (l *Localizer) Print(format string, vars ...interface{}) string {
	if l.locale != nil {
		return l.locale.printFunc(format, vars...)
	}

	return l.bundle.fallbackPrintFunc(format, vars...)
}

func (l *Localizer) dpGettextErr(domain localize.Domain, ctx localize.Context, msgID localize.Singular, vars ...interface{}) (string, error) {
	if l.locale != nil {
		return l.locale.dpGettextErr(domain, ctx, msgID, vars...)
	}

	return l.bundle.fallbackPrintFunc(msgID, vars...), errMissingLocale
}

func (l *Localizer) dnpGettextErr(domain string, ctx localize.Context, singular localize.Singular, plural localize.Plural, n interface{}, vars ...interface{}) (string, error) {
	if l.locale != nil {
		return l.locale.dnpGettextErr(domain, ctx, singular, plural, n, vars...)
	}

	if idx := l.bundle.fallbackPluralFunc(n); idx == 0 || plural == "" {
		return l.bundle.fallbackPrintFunc(singular, vars...), errMissingLocale
	}

	return l.bundle.fallbackPrintFunc(plural, vars...), errMissingLocale
}
