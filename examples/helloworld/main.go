package main

import (
	"fmt"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak"
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

	t = spreak.NewLocalizer(bundle, language.English)
	fmt.Println(t.Get("Hello world"))

	// Output:
	// Hola Mundo
	// Hello world
}
