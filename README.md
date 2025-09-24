# Spreak ![Test status](https://github.com/vorlif/spreak/workflows/Test/badge.svg) [![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE) [![PkgGoDev](https://pkg.go.dev/badge/github.com/vorlif/spreak)](https://pkg.go.dev/github.com/vorlif/spreak) [![Go Report Card](https://goreportcard.com/badge/github.com/vorlif/spreak)](https://goreportcard.com/report/github.com/vorlif/spreak) [![codecov](https://codecov.io/gh/vorlif/spreak/branch/main/graph/badge.svg?token=N1O0ZE1OFW)](https://codecov.io/gh/vorlif/spreak) ![MinVersion](https://img.shields.io/badge/Go-1.23+-blue)

Flexible translation and humanization library for Go, based on the concepts behind gettext. Requires Go 1.23+.

### Why another library?

There are already many good libraries for Go, which allow localizing an application.
However, I always came to a point where I was dissatisfied.
Either they use a self defined format, which could not be edited with common tools.
Some libraries only support one language at a time or are using a lot of mutexes.
And no tool could easily extract the strings to be translated.
I wanted to solve these problems for myself, and so spreak was born.

### Features

* Built-in support for `po`, `mo` and [`json`](./examples/features/jhttptempl) files
* Support for `fs.FS` (e.g. `embed`)
* Goroutine-safe and lock free through immutability
* [Powerful extractor](https://github.com/vorlif/xspreak#xspreak) for strings to simplify the localization process
  (with **support for templates**)
* [Support](https://pkg.go.dev/github.com/vorlif/spreak#hdr-Plurals)
  for [gettext](https://www.gnu.org/software/gettext/manual/html_node/Plural-forms.html)
and [CLDR v47](https://cldr.unicode.org/index/cldr-spec/plural-rules) plural rules.
* Support of bilingual and monolingual formats

### Usage

Using spreak is easy. First, use go get to install the latest version of the library.

```shell
go get -u github.com/vorlif/spreak
```

After that, spreak offers you a comprehensive interface to load and query your translations.

```go
package main

import (
	"fmt"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak"
	"github.com/vorlif/spreak/localize"
)

func main() {
	// Create a bundle that loads the translations for the required languages.
	// Typically, you only need one bundle in an application.
	bundle, err := spreak.NewBundle(
		// Set the language used in the program code/templates
		spreak.WithSourceLanguage(language.English),
		// Set the path from which the translations should be loaded
		spreak.WithDomainPath(spreak.NoDomain, "./locale"),
		// Specify the languages you want to load
		spreak.WithLanguage(language.German, language.Spanish, language.Chinese),
	)
	if err != nil {
		panic(err)
	}

	// Create a Localizer to select the language to translate.
	t := spreak.NewLocalizer(bundle, language.Spanish)

	// Translate
	fmt.Println(t.Get("Hello world"))
	fmt.Println(t.NGetf("I have %d dog", "I have %d dogs", 2, 2))
	fmt.Println(t.Localize(GetPlanet()))

	// Output: 
	// Hola Mundo
	// Tengo 2 perros
	// No conozco ningún planeta
}

func GetPlanet() *localize.Message {
	return &localize.Message{
		Singular: "I do not know any planet",
		Plural:   "I do not know any planets",
	}
}
```

### Extract strings

Strings for the translations can be extracted via the [command line program xspreak](https://github.com/vorlif/xspreak).
Use a [pre-built binary](https://github.com/vorlif/xspreak/releases/latest) or install it from source:

```bash
go install github.com/vorlif/xspreak@latest
```

Tests installation with:

```bash
xspreak --help
```

xspreak extracts the strings from the program code and creates a template which can be used for new translations.
Before you extract your strings, you have to decide on a format.

It can either create a `.pot` (PO Templates) file in po format or a `.json` file.
If you are not sure which format to use, I would recommend you to use `po` format,
because it is supported by almost all translation programs, which makes your life and the life of your translators much
easier.

With `-D` you specify the path to your source code and with `-p` the one where the translation template should be saved.

```bash
# for .pot
xspreak -D path/to/source/ -p path/to/source/locale
# for .json
xspreak -D path/to/source/ -p path/to/source/locale -f json
```

This creates a new `.pot` or `.json` file representing the translation template.

### Translate

**po files**

The generated POT files can be easily imported by most translation tools.
If you are dealing with Po files for the first time,
I recommend the application [Poedit](https://poedit.net/) for a quick start.

After translation `.po` or `.mo` files are generated, which are used by spreak for looking up translations.
Attention, do not translate the `.pot` file directly, as this is only a template.

**json files**

You can open and edit the JSON files with any text editor.
The extracted template always uses only the plural categories `one` and `other`. 
To create a template with the plural categories suitable for your language, you can use xspreak.

```bash
xspreak merge -i locale/template.json -o locale/de.json -l de
```

### Update translation files

> :warning: It would be wise before making any move to **keep a backup**.

When you add, edit or delete text in your program code, you should also update the translation files.
To achieve this, you must first update the template:

```bash
# for .pot
xspreak -D path/to/source/ -p path/to/source/locale
# for .json
xspreak -D path/to/source/ -p path/to/source/locale -f json
```

This creates a new `.pot` or `.json` file representing the *new* translation template.

**po files**

For PO files, almost every translation tool offers the possibility to update the files from a POT file.
With Poedit you can do it via `Translation -> Update from POT file`.
If you use the gettext utilities, you can use `msgmerge -U es.po template.pot`.
For all other tools, it is worth taking a look at the documentation.

**json files**

For JSON files, xspreak offers a simple way to update the translation files.

```bash
xspreak merge -i locale/template.json -o locale/de.json -l de
```

This copies new keys and existing translations from `template.json` to `de.json` and deletes keys from `de.json`
that are not present in `template.json`.

### Structure translations

How you structure the files with the translations is up to you.
Assuming you load the domain `"helloworld"` from the path `"./locale"` and the language `language.Spanish`

```go
spreak.WithDomainPath("helloworld", "./locale"),
spreak.WithLanguage(language.Spanish),
```

Then spreak searches for the following files by default

```text
./locale/es/helloworld.po
./locale/helloworld/es.po
./locale/es.po
./locale/LC_MESSAGES/es/helloworld.po
```

Where `es` is an example from the list `[es-ES, es_ES, spa, es]` and the file extension `.po` is an example from the
list `[po, mo, json]`.
If you don't like this behavior, you can implement your own [`Resolver`](examples/features/resolver/main.go).
For special cases you can also implement your own [`Loader`](examples/features/loaders/main.go).

### How to use in tests?

Just create a `Bundle` without options.
This will never return an error and can be used to create `Localizer` which then simply return the input.

```go
bundle, _ := spreak.NewBundle()
t := spreak.NewLocalizer(bundle, language.English)
```

### What's next

* Read what you can extract with [xspreak](https://github.com/vorlif/xspreak#xspreak)
* Take a look in the [examples folder](./examples) for more examples of using spreak.
* Use it!

## Package humanize

**Attention**: For a clean separation and easier maintenance, the `humanize` package has been separated
into a separate library and can now be found under https://github.com/vorlif/humanize.

The package in this repository is marked as deprecated and will be removed with the release of v2.

To change the library, simply
follow [the migration guide](https://github.com/vorlif/humanize?tab=readme-ov-file#migrate-from-spreak).

### Add translations

If you would like to add a translation or add a new language, **do not do so in this repository**.
The translations in this repository are automatically generated from the Django translations and additions should also
be made there.
Use the following link to do so: https://www.transifex.com/django/django/.
For all non-translation related errors, this repository must be used.

## License

spreak is available under the MIT license. See the [LICENSE](LICENSE) file for more info.
The translations of the `humanize` packages are licensed
under [Django's BSD license](https://raw.githubusercontent.com/django/django/main/LICENSE).
