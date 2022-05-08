package spreak_test

import (
	"fmt"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak"
	"github.com/vorlif/spreak/localize"
)

func ExampleNewBundle() {
	bundle, err := spreak.NewBundle(
		spreak.WithSourceLanguage(language.English),
		spreak.WithDefaultDomain("helloworld"),
		spreak.WithDomainPath("helloworld", "./examples/locale/"),
		spreak.WithLanguage(language.German, language.Spanish, language.French),
	)
	if err != nil {
		panic(err)
	}

	t := spreak.NewLocalizer(bundle, language.Spanish)

	fmt.Println(t.Get("Hello world"))
	// Output:
	// Hola Mundo
}

func ExampleLocalizer_Get() {
	bundle, err := spreak.NewBundle(
		spreak.WithSourceLanguage(language.English),
		spreak.WithDefaultDomain("helloworld"),
		spreak.WithDomainPath("helloworld", "./examples/locale/"),
		spreak.WithLanguage(language.German, language.Spanish, language.French),
	)
	if err != nil {
		panic(err)
	}

	t := spreak.NewLocalizer(bundle, language.Spanish)

	fmt.Println(t.Get("Hello world"))
	// Output:
	// Hola Mundo
}

func ExampleLocalizer_Localize() {
	bundle, err := spreak.NewBundle(
		spreak.WithDefaultDomain("helloworld"),
		spreak.WithDomainPath("helloworld", "./examples/locale/"),
		spreak.WithLanguage(language.German, language.Spanish, language.French),
	)
	if err != nil {
		panic(err)
	}

	t := spreak.NewLocalizer(bundle, language.Spanish)

	msg := &localize.Message{Singular: "Hello world"}
	fmt.Println(t.Localize(msg))
	// Output:
	// Hola Mundo
}

func ExampleLocale_Localize() {
	bundle, err := spreak.NewBundle(
		spreak.WithDefaultDomain("helloworld"),
		spreak.WithDomainPath("helloworld", "./examples/locale/"),
		spreak.WithLanguage(language.German, language.Spanish, language.French),
	)
	if err != nil {
		panic(err)
	}

	t := spreak.NewLocalizer(bundle, language.Spanish)

	msg := &localize.Message{Singular: "Hello world"}
	fmt.Println(t.Localize(msg))
	// Output:
	// Hola Mundo
}