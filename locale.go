package spreak

import (
	"golang.org/x/text/language"

	"github.com/vorlif/spreak/catalog"
	"github.com/vorlif/spreak/internal/poplural"
	"github.com/vorlif/spreak/localize"
)

// A locale holds the catalogs of all domains for a language and provides an interface for their use.
// It has a default domain, which can differ from the bundle default domain and is used
// if no domain is specified for translations.
// A locale needs at least one catalog and therefore there can be no locale,
// if no catalog is found for at least one domain.
type locale struct {
	bundle           *Bundle
	language         language.Tag
	domainCatalogs   map[string]catalog.Catalog
	printFunc        PrintFunc
	pluralFunc       poplural.PluralFunc
	isSourceLanguage bool
}

func buildLocale(bundle *Bundle, lang language.Tag, catalogs map[string]catalog.Catalog) *locale {
	l := &locale{
		bundle:         bundle,
		language:       lang,
		domainCatalogs: catalogs,
		printFunc:      bundle.printer.GetPrintFunc(lang),
	}

	l.pluralFunc, _ = poplural.ForLanguage(lang)

	return l
}

func buildSourceLocale(bundle *Bundle, sourceLang language.Tag) *locale {
	l := &locale{
		bundle:           bundle,
		language:         sourceLang,
		printFunc:        bundle.printer.GetPrintFunc(sourceLang),
		isSourceLanguage: true,
	}
	l.pluralFunc, _ = poplural.ForLanguage(sourceLang)
	return l
}

func (l *locale) lookupSingularTranslation(domain localize.Domain, ctx localize.Context, msgID localize.Singular, vars ...interface{}) (string, error) {
	if l.isSourceLanguage {
		return l.printFunc(msgID, vars...), nil
	}

	catl, err := l.getCatalog(domain)
	if err != nil {
		if l.bundle.missingCallback != nil {
			l.bundle.missingCallback(err)
		}

		return "", err
	}

	translation, errT := catl.GetTranslation(ctx, msgID)
	if errT != nil {
		if l.bundle.missingCallback != nil {
			l.bundle.missingCallback(errT)
		}

		return "", errT
	}

	return l.printFunc(translation, vars...), nil
}

func (l *locale) lookupPluralTranslation(domain string, ctx localize.Context, singular localize.Singular, plural localize.Plural, n interface{}, vars ...interface{}) (string, error) {
	if l.isSourceLanguage {
		return l.printSourceMessage(singular, plural, n, vars), nil
	}

	catl, err := l.getCatalog(domain)
	if err != nil {
		if l.bundle.missingCallback != nil {
			l.bundle.missingCallback(err)
		}

		return "", err
	}

	translation, errT := catl.GetPluralTranslation(ctx, singular, n)
	if errT != nil {
		if l.bundle.missingCallback != nil {
			l.bundle.missingCallback(errT)
		}

		return "", errT
	}

	return l.printFunc(translation, vars...), nil
}

func (l *locale) printSourceMessage(singular, plural string, n interface{}, vars []interface{}) string {
	idx := l.pluralFunc(n)
	if idx == 0 || plural == "" {
		return l.printFunc(singular, vars...)
	}

	return l.printFunc(plural, vars...)
}

func (l *locale) getCatalog(domain string) (catalog.Catalog, error) {
	if _, hasDomain := l.domainCatalogs[domain]; !hasDomain {
		err := &ErrMissingDomain{
			Language: l.language,
			Domain:   domain,
		}
		return nil, err
	}

	return l.domainCatalogs[domain], nil
}
