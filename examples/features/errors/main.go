package main

import (
	"errors"
	"fmt"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak"
)

// When xspreak is run with the -e option, the strings for errors.New are also extracted.
// They can then be easily translated.
var (
	// TRANSLATORS: It is the name of a file on the local computer.
	//
	// This comment is not extracted because a blank line was inserted above it.
	ErrInvalidName = errors.New("name has an invalid format")
)

func main() {
	bundle, err := spreak.NewBundle(
		spreak.WithSourceLanguage(language.English),
		spreak.WithDefaultDomain("errors"),
		spreak.WithDomainPath("errors", "../../locale"),
		spreak.WithLanguage(language.German),
	)
	if err != nil {
		panic(err)
	}

	t := spreak.NewLocalizer(bundle, language.German)
	fmt.Println(t.LocalizeError(ErrInvalidName))

	// The following error is not automatically extracted
	// because it is not known what the entire string is until the program runs.
	untranslatable := fmt.Errorf("a mistake has happened: %w", ErrInvalidName)
	fmt.Println(untranslatable)

	translatable := fmt.Errorf(t.Get("another mistake has happened: %s"), t.LocalizeError(ErrInvalidName))
	fmt.Println(translatable)

	// Output:
	// der Name hat ein ungültiges Format
	// a mistake has happened: name has an invalid format
	// ein weiterer Fehler ist passiert: der Name hat ein ungültiges Format
}
