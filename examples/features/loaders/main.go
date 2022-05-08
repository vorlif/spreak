package main

import (
	"fmt"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak"
)

const (
	Domain = "mydomain"
)

// A loader allows you to load data from any source.
// For example, you can use it to load the translation from the database, a MinIO instance,
// or from your Facebook page.
func main() {
	bundle, err := spreak.NewBundle(
		spreak.WithSourceLanguage(language.English),
		spreak.WithDefaultDomain(Domain),
		spreak.WithDomainLoader(Domain, &myLoader{}),
		spreak.WithLanguage(language.Spanish),
	)
	if err != nil {
		panic(err)
	}

	t := spreak.NewLocalizer(bundle, language.Spanish)
	fmt.Println(t.Get("Hello world"))
	// Output: Hola Mundo
}

type myLoader struct {}

func (myLoader) Load(lang language.Tag, domain string) (spreak.Catalog, error) {
	if lang == language.Spanish && domain == Domain {
		decoder := spreak.NewPoDecoder()
		return decoder.Decode(lang, domain, esTranslations)
	}

	return nil, spreak.NewErrNotFound(lang, "code", "domain=%q", domain)
}

var _ spreak.Loader = (*myLoader)(nil)

var esTranslations = []byte(`
msgid "Hello world"
msgstr "Hola Mundo"
`)