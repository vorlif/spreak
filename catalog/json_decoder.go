package catalog

import (
	"encoding/json"

	"golang.org/x/text/language"
)

// JSONDecoder is a Decoder for reading JSON catalog files.
type JSONDecoder struct{}

var _ Decoder = (*JSONDecoder)(nil)

// NewJSONDecoder returns a new Decoder for reading JSON catalog files.
// The structure follows a key-value structure, where the key is either an ID or the singular text of the source language.
// For singular-only texts, the value is a string with a translation.
// For plural texts it is an object with the CLDR plural forms and the matching translations.
func NewJSONDecoder() *JSONDecoder { return &JSONDecoder{} }

func (JSONDecoder) Decode(lang language.Tag, domain string, data []byte) (Catalog, error) {
	catl := NewJSONCatalog(lang, domain)
	if err := json.Unmarshal(data, &catl); err != nil {
		return nil, err
	}

	return catl, nil
}
