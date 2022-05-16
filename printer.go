package spreak

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// PrintFunc formats according to a format specifier and returns the resulting string.
// Like fmt.Sprintf(...)
type PrintFunc func(str string, vars ...interface{}) string

// A Printer creates a PrintFunc for a language.
// Can be stored with WithPrinter when creating a bundle.
type Printer interface {
	Init(languages []language.Tag)
	GetPrintFunc(lang language.Tag) PrintFunc
}

type defaultPrinter struct {
	printers map[language.Tag]PrintFunc
}

// NewDefaultPrinter creates a printer which will be used if no printer was defined
// with WithPrinter when creating a bundle.
func NewDefaultPrinter() Printer {
	return &defaultPrinter{}
}

func (d *defaultPrinter) Init(languages []language.Tag) {
	d.printers = make(map[language.Tag]PrintFunc, len(languages))
	for _, lang := range languages {
		d.printers[lang] = defaultPrintFunc(lang)
	}
}

func (d *defaultPrinter) GetPrintFunc(lang language.Tag) PrintFunc {
	if printFunc, ok := d.printers[lang]; ok {
		return printFunc
	}

	return defaultPrintFunc(lang)
}

func defaultPrintFunc(lang language.Tag) PrintFunc {
	printer := message.NewPrinter(lang)
	return func(str string, vars ...interface{}) string {
		if len(vars) > 0 {
			return printer.Sprintf(str, vars...)
		}

		return str
	}
}

// Simple wrapper to use a PrinterFunction as a printer.
type printFunctionWrapper struct {
	f PrintFunc
}

var _ Printer = (*printFunctionWrapper)(nil)

func (p *printFunctionWrapper) Init(_ []language.Tag)                 {}
func (p *printFunctionWrapper) GetPrintFunc(_ language.Tag) PrintFunc { return p.f }
