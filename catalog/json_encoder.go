package catalog

import (
	"encoding/json"
	"io"
)

// JSONEncoder is an encoder for writing JSON catalog files.
type JSONEncoder struct {
	w io.Writer
}

// NewJSONEncoder returns a new encoder that writes a JSON catalog to w.
func NewJSONEncoder(w io.Writer) *JSONEncoder { return &JSONEncoder{w: w} }

// Encode converts a JSON messages map into the content of a JSON catalog file.
//
// If nil or an empty map is passed, an empty JSON object is written.
// If the map contains an empty message ID, an error ist returned.
func (enc JSONEncoder) Encode(catl JSONCatalog) error {
	if catl == nil {
		_, err := enc.w.Write([]byte("{}"))
		return err
	}

	jsonEnc := json.NewEncoder(enc.w)
	jsonEnc.SetIndent("", "  ")
	if err := jsonEnc.Encode(catl); err != nil {
		return err
	}

	return nil
}
