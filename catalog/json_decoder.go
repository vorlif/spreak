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
	cat := NewJSONCatalog(lang, domain)
	if err := json.Unmarshal(data, &cat); err != nil {
		return nil, err
	}

	return cat, nil
}
