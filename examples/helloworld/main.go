package main

import (
	"fmt"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/v2"
	"github.com/vorlif/spreak/v2/localize"
)

func main() {
	bundle, err := spreak.NewBundle(
		spreak.WithSourceLanguage(language.English),
		spreak.WithDefaultDomain("helloworld"),
		spreak.WithDomainPath("helloworld", "../locale"),
		spreak.WithRequiredLanguage(language.Spanish),
		spreak.WithLanguage(language.German, language.French),
	)
	if err != nil {
		panic(err)
	}

	t := spreak.NewLocalizer(bundle, language.Spanish)

	// TRANSLATORS: This comment is automatically extracted by xspreak
	// and can be used to leave useful hints for the translators.
	fmt.Println(t.Get("Hello world"))
	fmt.Println(t.Localize(GetPlanet()))

	t = spreak.NewLocalizer(bundle, language.English)
	fmt.Println(t.Get("Hello world"))

	// Output:
	// Hola Mundo
	// No conozco ningún planeta
	// Hello world
}

func GetPlanet() *localize.Message {
	return &localize.Message{
		Singular: "I do not know any planet",
		Plural:   "I do not know any planets",
	}
}
