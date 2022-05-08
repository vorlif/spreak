package spreak

import (
	"errors"
	"fmt"
	"io/fs"

	"golang.org/x/text/language"
)

// BundleOption is an option which can be passed when creating a bundle to customize its configuration.
type BundleOption func(opts *bundleBuilder) error

// WithFallbackLanguage sets the fallback language to be used when creating Localizer if no suitable language is available.
// Should be used only if the fallback language is different from source language.
// Otherwise it should not be set.
func WithFallbackLanguage(lang interface{}) BundleOption {
	return func(opts *bundleBuilder) error {
		tag, err := languageInterfaceToTag(lang)
		if err != nil {
			return err
		}

		opts.fallbackLanguage = tag
		opts.setupActions = append(opts.setupActions, func(builder *bundleBuilder) error {
			return builder.preloadLanguages(true, tag)
		})
		return nil
	}
}

// WithSourceLanguage sets the source language used for programming.
// If it is set, it will be considered as a matching language when creating a Localizer.
// Also, it will try to use the appropriate plural form and will not trigger any missing callbacks for the language.
// It is recommended to always set the source language.
func WithSourceLanguage(tag language.Tag) BundleOption {
	return func(opts *bundleBuilder) error {
		opts.sourceLanguage = tag
		opts.setupActions = append(opts.setupActions, func(builder *bundleBuilder) error {
			return builder.preloadLanguages(false, tag)
		})
		return nil
	}
}

// WithMissingTranslationCallback stores a MissingTranslationCallback which is called when a translation,
// domain or something else is missing.
// The call is not goroutine safe.
func WithMissingTranslationCallback(cb MissingTranslationCallback) BundleOption {
	return func(opts *bundleBuilder) error {
		opts.missingCallback = cb
		return nil
	}
}

// WithDefaultDomain sets the default domain which will be used if no domain is specified.
// By default, NoDomain (the empty string) is used.
func WithDefaultDomain(domain string) BundleOption {
	return func(opts *bundleBuilder) error {
		opts.defaultDomain = domain
		return nil
	}
}

// WithDomainLoader loads a domain via a specified loader.
func WithDomainLoader(domain string, l Loader) BundleOption {
	return func(opts *bundleBuilder) error {
		if _, found := opts.domainLoaders[domain]; found {
			return fmt.Errorf("spreak.Bundle: loader for domain %s already set", domain)
		}
		if l == nil {
			return errors.New("spreak.Bundle: loader of WithDomainLoader(..., loader) is nil")
		}
		opts.domainLoaders[domain] = l
		return nil
	}
}

// WithFilesystemLoader Loads a domain via a FilesystemLoader.
// The loader can be customized with options.
func WithFilesystemLoader(domain string, fsOpts ...FsOption) BundleOption {
	return func(opts *bundleBuilder) error {
		l, err := NewFilesystemLoader(fsOpts...)
		if err != nil {
			return err
		}

		if _, found := opts.domainLoaders[domain]; found {
			return fmt.Errorf("%w: loader for domain %s already set", errSpreak, domain)
		}

		opts.domainLoaders[domain] = l
		return nil
	}
}

// WithDomainPath loads a domain from a specified path.
func WithDomainPath(domain string, path string) BundleOption {
	return WithFilesystemLoader(domain, WithLoadPath(path))
}

// WithDomainFs loads a domain from a fs.FS.
func WithDomainFs(domain string, fsys fs.FS) BundleOption {
	if fsys == nil {
		return func(opts *bundleBuilder) error {
			return errors.New("spreak.Bundle: fsys of WithDomainFs(..., fsys) is nil")
		}
	}

	return WithFilesystemLoader(domain, WithLoaderFs(fsys))
}

// WithLanguage loads the catalogs of the domains for one or more languages.
// The passed languages must be of type string or language.Tag,
// all other values will abort the initialization of the bundle with an error.
// If a catalog file for a domain is not found for a language, it will be ignored.
// If a catlaog file for a domain is found but cannot be loaded, the bundle creation will fail and return errors.
//
// If you want to use a Localizer, you should pay attention to the order in which the languages are specified,
// otherwise unexpected behavior may occur.
// This is because the matching algorithm of the language.matcher can give unexpected results.
// See https://github.com/golang/go/issues/49176
func WithLanguage(languages ...interface{}) BundleOption {
	loadFunc := func(builder *bundleBuilder) error {
		return builder.preloadLanguages(true, languages...)
	}

	return func(opts *bundleBuilder) error {
		opts.setupActions = append(opts.setupActions, loadFunc)
		return nil
	}
}

// WithRequiredLanguage works like WithLanguage except that the creation of the bundle fails
// if a catalog for a language could not be found.
func WithRequiredLanguage(languages ...interface{}) BundleOption {
	loadFunc := func(builder *bundleBuilder) error {
		return builder.preloadLanguages(false, languages...)
	}

	return func(opts *bundleBuilder) error {
		opts.setupActions = append(opts.setupActions, loadFunc)
		return nil
	}
}

// WithPrinter sets a printer which creates a function for a language which converts a formatted string
// and variables into a string. (Like fmt.Sprintf).
func WithPrinter(p Printer) BundleOption {
	return func(opts *bundleBuilder) error {
		if p == nil {
			return errors.New("spreak.Bundle: printer of WithPrinter(...) is nil")
		}
		opts.printer = p
		return nil
	}
}

// WithPrintFunction sets a PrintFunc which converts a formatted string and variables to a string. (Like fmt.Sprintf).
func WithPrintFunction(printFunc PrintFunc) BundleOption {
	if printFunc != nil {
		printer := &printFunctionWrapper{f: printFunc}
		return WithPrinter(printer)
	}

	return func(opts *bundleBuilder) error {
		return errors.New("spreak.Bundle: parameter of WithPrintFunction(...) is nil")
	}
}

// WithLanguageMatcherBuilder sets a LanguageMatcherBuilder.
func WithLanguageMatcherBuilder(mc LanguageMatcherBuilder) BundleOption {
	return func(opts *bundleBuilder) error {
		if mc == nil {
			return errors.New("spreak.Bundle: MatchCreator of WithMatchCreator(...) is nil")
		}
		opts.languageMatcherBuilder = mc
		return nil
	}
}

// WithErrorContext set a context, which is used for the translation of errors.
// If no context is set, ErrorsCtx is used.
func WithErrorContext(ctx string) BundleOption {
	return func(opts *bundleBuilder) error {
		opts.errContext = ctx
		return nil
	}
}
