package catalog

import (
	"encoding/json"
	"io"
)

// JsonEncoder is an encoder for writing JSON catalog files.
type JsonEncoder struct {
	w io.Writer
}

// NewJsonEncoder returns a new encoder that writes a JSON catalog to w.
func NewJsonEncoder(w io.Writer) *JsonEncoder { return &JsonEncoder{w: w} }

// Encode converts a JSON messages map into the content of a JSON catalog file.
//
// If nil or an empty map is passed, an empty JSON object is written.
// If the map contains an empty message ID, an error ist returned.
func (enc JsonEncoder) Encode(catl JSONCatalog) error {
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
