package main

import (
	"embed"
	"fmt"
	"io/fs"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak"
)

const (
	HelloWorldDomain = "helloworld"
)

//go:embed locale/*
var locales embed.FS

func main() {
	fsys, _ := fs.Sub(locales, "locale")
	bundle, err := spreak.NewBundle(
		spreak.WithSourceLanguage(language.English),
		spreak.WithDefaultDomain(HelloWorldDomain),
		spreak.WithDomainFs(HelloWorldDomain, fsys),
		spreak.WithLanguage(language.German, language.Spanish, language.French),
	)
	if err != nil {
		panic(err)
	}

	t := spreak.NewLocalizer(bundle, language.Spanish)

	fmt.Println(t.Get("Hello world"))
	// Output: Hola Mundo
}
