package spreak

import (
	"golang.org/x/text/language"

	"github.com/vorlif/spreak/internal/plural"
	"github.com/vorlif/spreak/localize"
)

// A Locale holds the catalogs of all domains for a language and provides an interface for their use.
// It has a default domain, which can differ from the bundle default domain and is used
// if no domain is specified for translations.
// A locale needs at least one catalog and therefore there can be no locale,
// if no catalog is found for at least one domain.
type Locale struct {
	bundle         *Bundle
	language       language.Tag
	defaultDomain  string
	domainCatalogs map[string]Catalog
	printFunc      PrintFunc

	pluralFunc       pluralFunction
	isSourceLanguage bool
}

// NewLocale creates a new locale for a language and the default domain of the bundle.
// If a locale is not found, an error is returned.
func NewLocale(bundle *Bundle, lang language.Tag) (*Locale, error) {
	return NewLocaleWithDomain(bundle, lang, bundle.defaultDomain)
}

// NewLocaleWithDomain creates a new locale for a language and a default domain.
// If no locale is found, an error is returned.
func NewLocaleWithDomain(bundle *Bundle, lang language.Tag, domain string) (*Locale, error) {
	if locale, hasLocale := bundle.locales[lang]; hasLocale {
		if domain == bundle.defaultDomain {
			return locale, nil
		}
		clone := cloneLocale(bundle.locales[lang], domain)
		return clone, nil
	}

	return nil, newMissingLanguageError(lang)
}

func buildLocale(bundle *Bundle, lang language.Tag, catalogs map[string]Catalog) *Locale {
	l := &Locale{
		bundle:         bundle,
		language:       lang,
		defaultDomain:  bundle.defaultDomain,
		domainCatalogs: catalogs,
		printFunc:      bundle.printer.GetPrintFunc(lang),
	}

	l.pluralFunc, _ = plural.ForLanguage(lang)

	return l
}

func buildSourceCodeLocale(bundle *Bundle) *Locale {
	l := &Locale{
		bundle:           bundle,
		language:         bundle.sourceLanguage,
		defaultDomain:    bundle.defaultDomain,
		printFunc:        bundle.printer.GetPrintFunc(bundle.sourceLanguage),
		isSourceLanguage: true,
	}
	l.pluralFunc, _ = plural.ForLanguage(bundle.sourceLanguage)
	return l
}

func cloneLocale(orig *Locale, defaultDomain string) *Locale {
	return &Locale{
		bundle:           orig.bundle,
		language:         orig.language,
		defaultDomain:    defaultDomain,
		domainCatalogs:   orig.domainCatalogs,
		printFunc:        orig.printFunc,
		isSourceLanguage: orig.isSourceLanguage,
	}
}

// Language returns the language for which the locale is used.
func (l *Locale) Language() language.Tag { return l.language }

// DefaultDomain returns the default domain.
// The default domain is used if a domain is not explicitly specified for a requested translation.
// If no default domain is specified, the default domain of the bundle is used.
func (l *Locale) DefaultDomain() string { return l.defaultDomain }

// HasDomain checks whether a catalog has been loaded for a specified domain.
func (l *Locale) HasDomain(domain string) bool {
	_, hasDomain := l.domainCatalogs[domain]
	return hasDomain
}

// Domains returns a list of all domains for which a catalog was found.
func (l *Locale) Domains() []string {
	domains := make([]string, 0, len(l.domainCatalogs))
	for domain := range l.domainCatalogs {
		domains = append(domains, domain)
	}
	return domains
}

// WithDomain creates a copy of the Locale with a different default domain.
// The original Locale will not be modified.
func (l *Locale) WithDomain(domain string) *Locale {
	return cloneLocale(l, domain)
}

// The Get function return the localized translation of message, based on the current default domain and
// language of the Locale.
// The message argument identifies the message to be translated.
// If no suitable translation exists, the original message is returned.
func (l *Locale) Get(message localize.Singular) string {
	t, _ := l.dpGettextErr(l.defaultDomain, NoCtx, message)
	return t
}

// Getf operates like Get, but formats the message according to a format identifier and returns the resulting string.
func (l *Locale) Getf(message localize.Singular, vars ...interface{}) string {
	t, _ := l.dpGettextErr(l.defaultDomain, NoCtx, message, vars...)
	return t
}

// DGet operates like Get, but look the message up in the specified domain.
func (l *Locale) DGet(domain localize.Domain, message localize.Singular) string {
	t, _ := l.dpGettextErr(domain, NoCtx, message)
	return t
}

// DGetf operates like Get, but look the message up in the specified domain and
// formats the message according to a format identifier and returns the resulting string.
func (l *Locale) DGetf(domain localize.Domain, message localize.Singular, vars ...interface{}) string {
	t, _ := l.dpGettextErr(domain, NoCtx, message, vars...)
	return t
}

// NGet acts like Get, but consider plural forms.
// The plural formula is applied to n and return the resulting message (some languages have more than two plurals).
func (l *Locale) NGet(singular localize.Singular, plural localize.Plural, n interface{}) string {
	t, _ := l.dnpGettextErr(l.defaultDomain, NoCtx, singular, plural, n)
	return t
}

// NGetf operates like NGet, but formats the message according to a format identifier and returns the resulting string.
func (l *Locale) NGetf(singular localize.Singular, plural localize.Plural, n interface{}, vars ...interface{}) string {
	t, _ := l.dnpGettextErr(l.defaultDomain, NoCtx, singular, plural, n, vars...)
	return t
}

// DNGet operates like NGet, but look the message up in the specified domain.
func (l *Locale) DNGet(domain localize.Domain, singular localize.Singular, plural localize.Plural, n interface{}) string {
	t, _ := l.dnpGettextErr(domain, NoCtx, singular, plural, n)
	return t
}

// DNGetf operates like DNGet, but formats the message according to a format identifier and returns the resulting string.
func (l *Locale) DNGetf(domain localize.Domain, singular localize.Singular, plural localize.Plural, n interface{}, vars ...interface{}) string {
	t, _ := l.dnpGettextErr(domain, NoCtx, singular, plural, n, vars...)
	return t
}

// PGet operates like Get, but restricted to the specified context.
func (l *Locale) PGet(context localize.Context, message localize.Singular) string {
	t, _ := l.dpGettextErr(l.defaultDomain, context, message)
	return t
}

// PGetf operates like PGet, but formats the message according to a format identifier and returns the resulting string.
func (l *Locale) PGetf(context localize.Context, message localize.Singular, vars ...interface{}) string {
	t, _ := l.dpGettextErr(l.defaultDomain, context, message, vars...)
	return t
}

// DPGet operates like Get, but look the message up in the specified domain and with the specified context.
func (l *Locale) DPGet(domain localize.Domain, context localize.Context, message localize.Singular) string {
	t, _ := l.dpGettextErr(domain, context, message)
	return t
}

// DPGetf operates like DPGet, but formats the message according to a format identifier and returns the resulting string.
func (l *Locale) DPGetf(domain localize.Domain, context localize.Context, message localize.Singular, vars ...interface{}) string {
	t, _ := l.dpGettextErr(domain, context, message, vars...)
	return t
}

// NPGet operates like NGet, but restricted to the specified context.
func (l *Locale) NPGet(context localize.Context, singular localize.Singular, plural localize.Plural, n interface{}) string {
	t, _ := l.dnpGettextErr(l.defaultDomain, context, singular, plural, n)
	return t
}

// NPGetf operates like NPGet, but formats the message according to a format identifier and returns the resulting string.
func (l *Locale) NPGetf(context localize.Context, singular localize.Singular, plural localize.Plural, n interface{}, vars ...interface{}) string {
	t, _ := l.dnpGettextErr(l.defaultDomain, context, singular, plural, n, vars...)
	return t
}

// DNPGet operates like NGet, but look the message up in the specified domain and with the specified context.
func (l *Locale) DNPGet(domain localize.Domain, context localize.Context, singular localize.Singular, plural localize.Plural, n interface{}) string {
	t, _ := l.dnpGettextErr(domain, context, singular, plural, n)
	return t
}

// DNPGetf operates like DNPGet, but formats the message according to a format identifier and returns the resulting string.
func (l *Locale) DNPGetf(domain localize.Domain, context localize.Context, singular localize.Singular, plural localize.Plural, n interface{}, vars ...interface{}) string {
	t, _ := l.dnpGettextErr(domain, context, singular, plural, n, vars...)
	return t
}

// LocalizeWithError translates structs that implement the interface localize.Localizable.
// If a suitable translation is found, it will be returned.
// If no matching translation is found, the original string with the matching plural form and an error are returned.
func (l *Locale) LocalizeWithError(t localize.Localizable) (string, error) {
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
func (l *Locale) Localize(t localize.Localizable) string {
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
func (l *Locale) LocalizeError(err error) error {
	switch v := err.(type) {
	case localize.Localizable:
		translation, errT := l.LocalizeWithError(v)
		if errT != nil {
			return err
		}

		return &localize.Error{Translation: translation, Wrapped: err}
	default:
		translation, errT := l.dpGettextErr(l.defaultDomain, l.bundle.errContext, err.Error())
		if errT != nil {
			return err
		}

		return &localize.Error{Translation: translation, Wrapped: err}
	}
}

func (l *Locale) dpGettextErr(domain localize.Domain, ctx localize.Context, msgID localize.Singular, vars ...interface{}) (string, error) {
	if l.isSourceLanguage {
		return l.printFunc(msgID, vars...), nil
	}

	catl, err := l.getCatalog(domain)
	if err != nil {
		if l.bundle.missingCallback != nil {
			l.bundle.missingCallback(err)
		}

		return l.printFunc(msgID, vars...), err
	}

	translation, errT := catl.GetTranslation(ctx, msgID)
	if errT != nil {
		if l.bundle.missingCallback != nil {
			l.bundle.missingCallback(errT)
		}

		return l.printFunc(msgID, vars...), errT
	}

	return l.printFunc(translation, vars...), nil
}

func (l *Locale) dnpGettextErr(domain string, ctx localize.Context, singular localize.Singular, plural localize.Plural, n interface{}, vars ...interface{}) (string, error) {
	if l.isSourceLanguage {
		return l.printSourceMessage(singular, plural, n, vars), nil
	}

	catl, err := l.getCatalog(domain)
	if err != nil {
		if l.bundle.missingCallback != nil {
			l.bundle.missingCallback(err)
		}

		return l.printSourceMessage(singular, plural, n, vars), err
	}

	translation, errT := catl.GetPluralTranslation(ctx, singular, n)
	if errT != nil {
		if l.bundle.missingCallback != nil {
			l.bundle.missingCallback(errT)
		}

		return l.printSourceMessage(singular, plural, n, vars), errT
	}

	return l.printFunc(translation, vars...), nil
}

func (l *Locale) printSourceMessage(singular, plural string, n interface{}, vars []interface{}) string {
	idx := l.pluralFunc(n)
	if idx == 0 || plural == "" {
		return l.printFunc(singular, vars...)
	}

	return l.printFunc(plural, vars...)
}

func (l *Locale) getCatalog(domain string) (Catalog, error) {
	if _, hasDomain := l.domainCatalogs[domain]; !hasDomain {
		err := &ErrMissingDomain{
			Language: l.language,
			Domain:   domain,
		}
		return nil, err
	}

	return l.domainCatalogs[domain], nil
}
