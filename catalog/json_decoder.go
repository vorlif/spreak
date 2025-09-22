package catalog

import (
	"encoding/json"

	"golang.org/x/text/language"
)

// JSONDecoder is a Decoder for reading JSON catalog files.
type JSONDecoder struct{}

var _ Decoder = (*JSONDecoder)(nil)

// NewJSONDecoder returns a new Decoder for reading JSON catalog files.
func NewJSONDecoder() *JSONDecoder { return &JSONDecoder{} }

func (JSONDecoder) Decode(lang language.Tag, domain string, data []byte) (Catalog, error) {
	catl := NewJSONCatalog(lang, domain)
	if err := json.Unmarshal(data, &catl); err != nil {
		return nil, err
	}

	return catl, nil
}
