package main

import (
	"encoding/json"
	"fmt"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak"
)

func main() {
	// We want spreak to load our .json files here, so we create a FilesystemLoader here with our own decoder.
	fsLoader, errFS := spreak.NewFilesystemLoader(
		spreak.WithDecoder(".json", jsonDecoder{}),
		spreak.WithPath("./"),
	)
	if errFS != nil {
		panic(errFS)
	}

	bundle, err := spreak.NewBundle(
		spreak.WithDomainLoader(spreak.NoDomain, fsLoader),
		spreak.WithRequiredLanguage(language.Spanish, language.English),
		spreak.WithLanguage(language.German, language.French),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(bundle.SupportedLanguages())

	t := spreak.NewLocalizer(bundle, language.Spanish)
	fmt.Println(t.Get("hello.world"))
	fmt.Println(t.NGetf("dog", "dog", 2, 2))
	fmt.Println(t.NGetf("cat", "cats", 2, 2))

	t = spreak.NewLocalizer(bundle, language.English)
	fmt.Println(t.Get("hello.world"))
	fmt.Println(t.NGetf("dog", "dogs", 2, 2))

	// Output:
	// [es en]
	// Hola Mundo
	// Tengo 2 perros
	// cats
	// Hello world
	// I have 2 dogs
}

// We create our own decoder which designs a JSON catalog from the content of the JSON files.
type jsonDecoder struct{}

var _ spreak.Decoder = (*jsonDecoder)(nil)

func (jsonDecoder) Decode(lang language.Tag, domain string, data []byte) (spreak.Catalog, error) {
	catalog := NewJsonCatalog(lang, domain)
	if err := json.Unmarshal(data, &catalog.translations); err != nil {
		return nil, err
	}
	return catalog, nil
}

// We create a catalog for our JSON files.
type jsonCatalog struct {
	language     language.Tag
	domain       string
	translations map[string]string
}

func NewJsonCatalog(lang language.Tag, domain string) *jsonCatalog {
	return &jsonCatalog{
		language:     lang,
		domain:       domain,
		translations: make(map[string]string),
	}
}

var _ spreak.Catalog = (*jsonCatalog)(nil)

func (c *jsonCatalog) GetTranslation(ctx, msgID string) (string, error) {
	return c.GetPluralTranslation(ctx, msgID, 1)
}

func (c *jsonCatalog) GetPluralTranslation(ctx, msgID string, n interface{}) (string, error) {
	if ctx != "" {
		err := &spreak.ErrMissingContext{
			Language: c.language,
			Domain:   c.domain,
			Context:  ctx,
		}
		return "", err
	}

	// Common plural form:
	// n == 1 -> singular
	// n != 1 -> plural
	idx := 0
	if n != 1 {
		idx = 1
		msgID += "_plural"
	}
	if _, hasMsg := c.translations[msgID]; !hasMsg {
		err := &spreak.ErrMissingMessageID{
			Language: c.language,
			Domain:   c.domain,
			Context:  ctx,
			MsgID:    msgID,
		}
		return "", err
	}

	translation := c.translations[msgID]
	if translation == "" {
		err := &spreak.ErrMissingTranslation{
			Language: c.language,
			Domain:   c.domain,
			Context:  ctx,
			MsgID:    msgID,
			Idx:      idx,
		}
		return "", err
	}

	return translation, nil
}

func (c *jsonCatalog) Language() language.Tag {
	return c.language
}
