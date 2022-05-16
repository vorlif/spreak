package spreak

import (
	"errors"

	"golang.org/x/text/language"
)

// NoDomain is the domain which is used if no default domain is stored.
const NoDomain = ""

// NoCtx is the context which is used if no context is stored.
const NoCtx = ""

// ErrorsCtx ist the context under which translations for extracted errors are searched.
// Can be changed when creating a bundle with WithErrorContext.
const ErrorsCtx = "errors"

// MissingTranslationCallback is a callback which can be stored with WithMissingTranslationCallback for a bundle.
// Called when translations, domains, or other are missing.
// The call is not goroutine safe.
type MissingTranslationCallback func(err error)

// LanguageMatcherBuilder is a builder which creates a language matcher.
// It is an abstraction of language.NewMatcher of the language package and should return the same values.
// Can be set when creating a bundle with WithLanguageMatcherBuilder for a bundle.
// The matcher is used, for example, when a new Localizer is created to determine the best matching language.
type LanguageMatcherBuilder func(t []language.Tag, options ...language.MatchOption) language.Matcher

type pluralFunction func(n interface{}) int
type setupAction func(options *bundleBuilder) error

type bundleBuilder struct {
	*Bundle

	languageMatcherBuilder LanguageMatcherBuilder
	domainLoaders          map[string]Loader
	setupActions           []setupAction
}

// A Bundle is the central place to load and manage translations.
// It holds all catalogs for all domains and all languages.
// The bundle cannot be edited after creation and is goroutine safe.
// Typically, an application contains a bundle as a singleton.
// The catalog of the specified domains and languages will be loaded during the creation.
type Bundle struct {
	missingCallback MissingTranslationCallback
	printer         Printer
	defaultDomain   string
	errContext      string

	sourceLanguage   language.Tag
	fallbackLanguage language.Tag
	languageMatcher  language.Matcher

	languages []language.Tag
	locales   map[language.Tag]*Locale
	domains   map[string]bool
}

// NewBundle creates a new bundle and returns it.
// It returns an error message if something failed during creation.
// The catalog of the specified domains and languages will be loaded during the creation.
func NewBundle(opts ...BundleOption) (*Bundle, error) {
	builder := &bundleBuilder{
		Bundle: &Bundle{
			printer:          NewDefaultPrinter(),
			defaultDomain:    NoDomain,
			sourceLanguage:   language.Und,
			fallbackLanguage: language.Und,
			errContext:       ErrorsCtx,

			languages: make([]language.Tag, 0),
			locales:   make(map[language.Tag]*Locale),
			domains:   make(map[string]bool),
		},
		languageMatcherBuilder: language.NewMatcher,

		domainLoaders: make(map[string]Loader),
		setupActions:  make([]setupAction, 0),
	}

	for _, opt := range opts {
		if opt == nil {
			return nil, errors.New("spreak.Bundle: try to create an bundle with a nil option")
		}
		if err := opt(builder); err != nil {
			return nil, err
		}
	}

	for domain := range builder.domainLoaders {
		builder.domains[domain] = false
	}

	for _, action := range builder.setupActions {
		if action == nil {
			return nil, errors.New("spreak.Bundle: try to create an bundle with a nil action")
		}
		if err := action(builder); err != nil {
			return nil, err
		}
	}

	builder.languageMatcher = builder.languageMatcherBuilder(builder.languages)
	builder.printer.Init(builder.languages)
	return builder.Bundle, nil
}

// Domains returns a list of loaded domains.
// A domain is only loaded if at least one catalog is found in one language.
func (b *Bundle) Domains() []string {
	domains := make([]string, 0, len(b.domains))
	for domain, loaded := range b.domains {
		if loaded {
			domains = append(domains, domain)
		}
	}
	return domains
}

// CanLocalize indicates whether locales and domains have been loaded for translation.
func (b *Bundle) CanLocalize() bool {
	return len(b.locales) > 0 && len(b.Domains()) > 0
}

// SupportedLanguages returns all languages for which a catalog was found for at least one domain.
func (b *Bundle) SupportedLanguages() []language.Tag {
	languages := make([]language.Tag, 0, len(b.locales))
	for lang := range b.locales {
		languages = append(languages, lang)
	}
	return languages
}

// IsLanguageSupported indicates whether a language can be translated.
// The check is done by the bundle's matcher and therefore languages that are not returned by
// SupportedLanguages can be supported.
func (b *Bundle) IsLanguageSupported(lang language.Tag) bool {
	_, _, confidence := b.languageMatcher.Match(lang)
	return confidence > language.No
}

func (b *bundleBuilder) preloadLanguages(optional bool, languages ...interface{}) error {
	for _, accept := range languages {
		tag, errT := languageInterfaceToTag(accept)
		if errT != nil {
			return errT
		}

		_, err := b.createLocale(optional, tag)
		if err == nil {
			continue
		}

		if !optional {
			return err
		}

		var missErr *ErrMissingLanguage
		if errors.As(err, &missErr) {
			if b.missingCallback != nil {
				b.missingCallback(missErr)
			}

			continue
		}

		return err

	}

	return nil
}

func (b *bundleBuilder) createLocale(optional bool, lang language.Tag) (*Locale, error) {
	if lang == language.Und {
		return nil, newMissingLanguageError(lang)
	}

	if lang == b.sourceLanguage {
		locale := buildSourceCodeLocale(b.Bundle)
		b.locales[lang] = locale
		b.languages = append(b.languages, lang)
		return locale, nil
	}

	if locale, isCached := b.locales[lang]; isCached {
		return locale, nil
	}

	catalogs := make(map[string]Catalog, len(b.domainLoaders))

	for domain, domainLoader := range b.domainLoaders {
		catl, errD := domainLoader.Load(lang, domain)
		if errD != nil {
			var notFoundErr *ErrNotFound
			if errors.As(errD, &notFoundErr) {
				if b.missingCallback != nil {
					b.missingCallback(notFoundErr)
				}

				if optional {
					continue
				}
			}
			return nil, errD
		}

		catalogs[domain] = catl
		b.domains[domain] = true
	}

	if len(catalogs) == 0 {
		return nil, newMissingLanguageError(lang)
	}

	locale := buildLocale(b.Bundle, lang, catalogs)
	b.locales[lang] = locale
	b.languages = append(b.languages, lang)
	return locale, nil
}
