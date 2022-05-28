package main

import (
	"fmt"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak"
	"github.com/vorlif/spreak/catalog"
)

const (
	Domain = "mydomain"
)

// A loader allows you to load data from any source.
// For example, you can use it to load the translation from the database, a MinIO instance,
// or from your Facebook page.
type myLoader struct{}

var _ spreak.Loader = (*myLoader)(nil)

func (myLoader) Load(lang language.Tag, domain string) (catalog.Catalog, error) {
	if lang == language.Spanish && domain == Domain {
		decoder := catalog.NewPoDecoder()
		return decoder.Decode(lang, domain, esTranslations)
	}

	return nil, spreak.NewErrNotFound(lang, "code", "domain=%q", domain)
}

// An example of loaded data. Here we use the format of the Po files since we have a decoder.
// You can use any format you like as long as you can create a catalog for it.
var esTranslations = []byte(`
msgid "Hello world"
msgstr "Hola Mundo"
`)

func main() {
	bundle, err := spreak.NewBundle(
		spreak.WithSourceLanguage(language.English),
		spreak.WithDefaultDomain(Domain),
		// Use the loader
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
