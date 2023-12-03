package spreak

import (
	"golang.org/x/text/language"

	"github.com/vorlif/spreak/catalog"
	"github.com/vorlif/spreak/catalog/poplural"
	"github.com/vorlif/spreak/localize"
)

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

func (l *locale) lookupSingularTranslation(domain localize.Domain, ctx localize.Context, msgID localize.Singular, vars ...any) (string, error) {
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

func (l *locale) lookupPluralTranslation(domain string, ctx localize.Context, singular localize.Singular, plural localize.Plural, n any, vars ...any) (string, error) {
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

func (l *locale) printSourceMessage(singular, plural string, n any, vars []any) string {
	idx, errPlural := l.pluralFunc(n)
	if errPlural != nil && l.bundle.missingCallback != nil {
		l.bundle.missingCallback(errPlural)
	}

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
