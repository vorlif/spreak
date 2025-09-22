package main

import (
	"fmt"
	"regexp"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak"
	"github.com/vorlif/spreak/localize"
)

const (
	Domain = "printer"
)

func main() {
	bundle, err := spreak.NewBundle(
		spreak.WithSourceLanguage(language.English),
		spreak.WithDefaultDomain(Domain),
		spreak.WithDomainPath(Domain, "../../locale"),
		spreak.WithLanguage(language.German, language.Spanish, language.French),
		spreak.WithPrinter(NewPrinter()),
	)
	if err != nil {
		panic(err)
	}

	t := spreak.NewLocalizer(bundle, language.German)

	var msg localize.MsgID

	// TRANSLATORS: %{name}s is the name of a person
	msg = "My name is %{name}s and I am %{age}d years old"

	fmt.Println(t.Getf(msg, "name", "Bob", "age", 8))
	// Output: Mein Name ist Bob und ich bin 8 Jahre alt
}

// A printer creates a function for each language, which is responsible for embedding the variables in the string.
// Like fmt.Sprintf(string, ...variables)
type myPrinter struct{}

var _ spreak.Printer = (*myPrinter)(nil)

func NewPrinter() spreak.Printer {
	return &myPrinter{}
}

func (m myPrinter) GetPrintFunc(lang language.Tag) spreak.PrintFunc {
	return func(str string, vars ...any) string {
		f, p := parse(str, vars...)
		return fmt.Sprintf(f, p...)
	}
}

// Information for which languages PrintFuncs probably have to be created.
func (myPrinter) Init(languages []language.Tag) {}

var re = regexp.MustCompile("%{(\\w+)}[.\\d]*[xsvTtbcdoqXUeEfFgGp]")

func parse(format string, vars ...any) (string, []any) {
	params := make(map[string]any, len(vars)/2)
	for i := 0; i < len(vars); i++ {
		key := fmt.Sprintf("%v", vars[i])
		if i+1 < len(vars) {
			params[key] = vars[i+1]
			i++
		} else {
			params[key] = key
		}
	}

	f, n := reformat(format)
	p := make([]any, len(n))
	for i, v := range n {
		p[i] = params[v]
	}
	return f, p
}

// The following code was copied from https://github.com/chonla/format

func reformat(f string) (string, []string) {
	i := re.FindAllStringSubmatchIndex(f, -1)

	var ord []string
	pair := []int{0}
	for _, v := range i {
		ord = append(ord, f[v[2]:v[3]])
		pair = append(pair, v[2]-1)
		pair = append(pair, v[3]+1)
	}
	pair = append(pair, len(f))
	plen := len(pair)

	out := ""
	for n := 0; n < plen; n += 2 {
		out += f[pair[n]:pair[n+1]]
	}

	return out, ord
}
