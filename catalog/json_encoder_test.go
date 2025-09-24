package catalog

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/vorlif/spreak/v2/catalog/cldrplural"
)

func TestJsonEncoder_Encode(t *testing.T) {
	t.Run("return an empty object on invalid input", func(t *testing.T) {
		buf := &bytes.Buffer{}
		enc := NewJSONEncoder(buf)

		for _, input := range []JSONCatalog{nil, NewJSONCatalog(language.English, "")} {
			buf.Reset()
			err := enc.Encode(input)
			assert.NoError(t, err)
			assert.JSONEq(t, "{}", buf.String())
		}
	})

	t.Run("msgId is reused", func(t *testing.T) {
		cat := NewJSONCatalog(language.English, "").(*jsonCatalog)
		cat.mustSetMessage("car_ctx", &JSONMessage{
			Context: "ctx",
			Translations: map[cldrplural.Category]string{
				cldrplural.One:   "Car",
				cldrplural.Other: "Cars",
			},
		})
		cat.mustSetMessage("car", &JSONMessage{
			Translations: map[cldrplural.Category]string{
				cldrplural.One:   "Car",
				cldrplural.Other: "Cars",
			},
		})

		res := `
{
	"car": {
		"one": "Car",
		"other": "Cars"
	},
	"car_ctx": {
		"context": "ctx",
		"one": "Car",
		"other": "Cars"
	}
}
`

		buf := &bytes.Buffer{}
		err := NewJSONEncoder(buf).Encode(cat)
		assert.NoError(t, err)
		assert.JSONEq(t, res, buf.String())
	})
}
