> Package is in development

# Spreak ![Build status](https://github.com/vorlif/spreak/workflows/Build/badge.svg) [![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE) [![PkgGoDev](https://pkg.go.dev/badge/github.com/vorlif/spreak)](https://pkg.go.dev/github.com/vorlif/spreak) [![Go Report Card](https://goreportcard.com/badge/github.com/vorlif/spreak)](https://goreportcard.com/report/github.com/vorlif/spreak)

Flexible translation library for Go based on the concepts of gettext.

## Features

* Support for `fs.FS` and `embed`
* Goroutine Safe and lock free through immutability
* Easily extendable
* [Powerful extractor](xspreak/README.md) for strings to simplify the localization process

## Usage

Spreak loads your translations and provides a simple interface for querying them.

```go
package main

import (
	"fmt"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak"
)

func main() {
	// Create a bundle that loads the translations for the required languages.
	// Typically, you only need one bundle in an application.
	bundle, err := spreak.NewBundle(
		spreak.WithSourceLanguage(language.English),
		// Set the path from which the translations should be loaded
		spreak.WithDomainPath(spreak.NoDomain, "../locale"),
		// Specify the languages you want to load
		spreak.WithLanguage(language.German, language.Spanish, language.Chinese),
	)
	if err != nil {
		panic(err)
	}

	// Create a Localiser to select the language to translate.
	t := spreak.NewLocalizer(bundle, language.Spanish)

	// Translate
	fmt.Println(t.Get("Hello world"))
	fmt.Println(t.NGetf("I have %d dog", "I have %d dogs", 2, 2))

	// Output:
	//  Hola Mundo
	//  Tengo 2 perros
}
```

### Extract strings

Strings for the translations can be extracted via the command line program xspreak.

```
go install github.com/vorlif/spreak/xspreak@main
xspreak -help
```

xspreak creates `.pot` (PO Templates) files for this purpose which can be easily imported directly by many translation
tools.

```bash
xspreak -D path/to/source/ -p path/to/source/locale
```

How you structure the files with the translations is up to you.
Common structures are:

```text
{path}/{language}/{domain}.po
{path}/{language}.po
{path}/{domain}/{language}po
{path}/{language}/{category}/{domain}.po

Example:
.../locale/es/helloworld.po
.../locale/es.po
.../locale/helloworld/es.po
.../locale/LC_MESSAGES/domain.po
```

By default spreak searches for translations in all these places.
If you don't like this behavior you can implement your own Reducer or Loader.

### Translate

The pot files generated by xspreak can be easily imported by most translation tools.
Most tools also have the ability to update existing translation via the pot file.

If you are dealing with Po files for the first time,
I recommend the application [poedit](https://poedit.net/) for a quick start.

After translation `.po` or `.mo` files are generated, which are used by spreak for looking up translations.
Attention, do not translate the `.pot` file directly, as this is only a template.

## What's next

* Read what you can extract with [xspreak](xspreak/README.md)
* Take a look in the [examples folder](./examples) for more examples of using spreak.

## License

spreak is available under the MIT license. See the [LICENSE](LICENSE) file for more info.