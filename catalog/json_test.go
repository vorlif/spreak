package catalog

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/vorlif/spreak/catalog/cldrplural"
)

var enJSONTestFile = filepath.FromSlash("../testdata/translation-test/json/en.json")
var deJSONTestFile = filepath.FromSlash("../testdata/translation-test/json/de.json")

func TestJsonMessage_UnmarshalJSON(t *testing.T) {
	var msg JSONMessage
	err := json.Unmarshal([]byte(`"test"`), &msg)
	assert.NoError(t, err)
	assert.Equal(t, msg.Translations[cldrplural.Other], "test")

	msg = JSONMessage{}
	err = json.Unmarshal([]byte(`{"context":"ctx","other":"test"}`), &msg)
	assert.NoError(t, err)
	assert.Equal(t, msg.Translations[cldrplural.Other], "test")
	assert.Equal(t, msg.Context, "ctx")

	t.Run("Unmarshals all values", func(t *testing.T) {

		msg = JSONMessage{}
		err = json.Unmarshal([]byte(`
{
	"context":"ctx",
	"comment":"comment",
	"zero": "zero",
	"one": "one",
	"two": "two",
	"few": "few",
	"many": "many",
	"other":"other"
}`), &msg)
		assert.NoError(t, err)
		assert.Equal(t, msg.Context, "ctx")
		assert.Equal(t, msg.Comment, "comment")
		assert.Len(t, msg.Translations, 6)
		for cat, value := range msg.Translations {
			assert.Equal(t, strings.ToLower(cat.String()), value)
		}
	})
}

func TestJsonMessage_MarshalJSON(t *testing.T) {
	t.Run("Marshals all values", func(t *testing.T) {
		msg := JSONMessage{
			Context: "ctx",
			Comment: "comment",
			Translations: map[cldrplural.Category]string{
				cldrplural.Zero:  "zero",
				cldrplural.One:   "one",
				cldrplural.Two:   "two",
				cldrplural.Few:   "few",
				cldrplural.Many:  "many",
				cldrplural.Other: "other",
			},
		}
		res, err := json.Marshal(msg)
		require.NoError(t, err)
		want := `{
	"context":"ctx",
	"comment":"comment",
	"zero": "zero",
	"one": "one",
	"two": "two",
	"few": "few",
	"many": "many",
	"other":"other"
}`
		assert.JSONEq(t, want, string(res))
	})
}

func TestJSONCatalog(t *testing.T) {
	t.Run("test translation lookup", func(t *testing.T) {
		data, err := os.ReadFile(enJSONTestFile)
		assert.NoError(t, err)
		require.NotNil(t, data)

		cat, errC := NewJSONDecoder().Decode(language.English, "a", data)
		assert.NoError(t, errC)
		require.NotNil(t, cat)

		assert.Equal(t, language.English, cat.Language())

		tr, err := cat.Lookup("", "app.name")
		assert.NoError(t, err)
		assert.Equal(t, "TODO List", tr)

		tr, err = cat.LookupPlural("", "animal.cat", 1)
		assert.NoError(t, err)
		assert.Equal(t, "I do not have a cat", tr)

		tr, err = cat.LookupPlural("", "animal.cat", 2)
		assert.NoError(t, err)
		assert.Equal(t, "I do not have cats", tr)

		tr, err = cat.LookupPlural("my-animals", "animal.dog", 2)
		assert.NoError(t, err)
		assert.Equal(t, "I have dogs", tr)
	})

	t.Run("test translation lookup errors", func(t *testing.T) {
		data, err := os.ReadFile(deJSONTestFile)
		assert.NoError(t, err)
		require.NotNil(t, data)

		cat, errC := NewJSONDecoder().Decode(language.German, "a", data)
		assert.NoError(t, errC)
		require.NotNil(t, cat)

		assert.Equal(t, language.German, cat.Language())

		tr, err := cat.Lookup("", "animal.cat")
		assert.Error(t, err)
		assert.Equal(t, "animal.cat", tr)

		tr, err = cat.LookupPlural("", "animal.cat", 1)
		assert.Error(t, err)
		assert.Equal(t, "animal.cat", tr)

		tr, err = cat.Lookup("unknown", "animal.dog")
		assert.Error(t, err)
		assert.Equal(t, "animal.dog", tr)

		tr, err = cat.LookupPlural("", "missing.plural", 1)
		assert.Error(t, err)
		assert.Equal(t, "missing.plural", tr)

		tr, err = cat.LookupPlural("my-animals", "animal.dog", 2)
		assert.NoError(t, err)
		assert.Equal(t, "Ich habe Hunde", tr)
	})
}

func Test_jsonCatalog_MarshalJSON(t *testing.T) {
	t.Run("Marshals singular as singel string", func(t *testing.T) {
		catl := NewJSONCatalog(language.English, "").(*jsonCatalog)
		catl.mustSetMessage("car", &JSONMessage{
			Translations: map[cldrplural.Category]string{
				cldrplural.Other: "Car",
			},
		})

		res := `{"car": "Car"}`

		buf := &bytes.Buffer{}
		err := NewJsonEncoder(buf).Encode(catl)
		assert.NoError(t, err)
		assert.JSONEq(t, res, buf.String())
	})

	t.Run("Marshals context as object", func(t *testing.T) {
		catl := NewJSONCatalog(language.English, "").(*jsonCatalog)
		catl.mustSetMessage("car_ctx", &JSONMessage{
			Context: "ctx",
			Translations: map[cldrplural.Category]string{
				cldrplural.Other: "Car",
			},
		})

		res := `{"car_ctx": {"context": "ctx","other": "Car"}}`
		buf := &bytes.Buffer{}
		err := NewJsonEncoder(buf).Encode(catl)
		assert.NoError(t, err)
		assert.JSONEq(t, res, buf.String())
	})

	t.Run("Marshals plural as object", func(t *testing.T) {
		catl := NewJSONCatalog(language.English, "").(*jsonCatalog)
		catl.mustSetMessage("car", &JSONMessage{
			Translations: map[cldrplural.Category]string{
				cldrplural.One:   "Car",
				cldrplural.Other: "Cars",
			},
		})

		res := `{"car": {"one": "Car","other": "Cars"}}`
		buf := &bytes.Buffer{}
		err := NewJsonEncoder(buf).Encode(catl)
		assert.NoError(t, err)
		assert.JSONEq(t, res, buf.String())
	})
}

func Test_jsonCatalog_UnmarshalJSON(t *testing.T) {

	t.Run("Unmarshal adds missing plural categories", func(t *testing.T) {
		catl := NewJSONCatalog(language.Polish, "").(*jsonCatalog)
		catl.mustSetMessage("car", &JSONMessage{
			Translations: map[cldrplural.Category]string{
				cldrplural.One:   "Car",
				cldrplural.Other: "Cars",
			},
		})

		translations := catl.lookupMap[""]["car"].Translations
		assert.Contains(t, translations, cldrplural.One)
		assert.Contains(t, translations, cldrplural.Few)
		assert.Contains(t, translations, cldrplural.Many)
		assert.Contains(t, translations, cldrplural.Other)
	})

	t.Run("Unmarshal removes missing plural categries", func(t *testing.T) {
		catl := NewJSONCatalog(language.English, "").(*jsonCatalog)
		catl.mustSetMessage("car", &JSONMessage{
			Translations: map[cldrplural.Category]string{
				cldrplural.One:   "Car",
				cldrplural.Few:   "Few cars",
				cldrplural.Other: "Cars",
			},
		})

		translations := catl.lookupMap[""]["car"].Translations
		assert.Contains(t, translations, cldrplural.One)
		assert.NotContains(t, translations, cldrplural.Few)
		assert.NotContains(t, translations, cldrplural.Many)
		assert.Contains(t, translations, cldrplural.Other)
	})
}

func TestNewJsonCatalogWithTranslations(t *testing.T) {
	t.Run("returns correct value", func(t *testing.T) {
		catl, err := NewJSONCatalogWithMessages(language.Polish, "", nil)
		assert.NoError(t, err)
		assert.NotNil(t, catl)

		assert.Equal(t, language.Polish, catl.Language())
		assert.Equal(t, "", catl.Domain())

		catl, err = NewJSONCatalogWithMessages(language.Polish, "domain", nil)
		assert.NoError(t, err)
		assert.NotNil(t, catl)
		assert.Equal(t, "domain", catl.Domain())
	})

	t.Run("empty map is valid", func(t *testing.T) {
		_, err := NewJSONCatalogWithMessages(language.Polish, "", nil)
		assert.NoError(t, err)

		_, err = NewJSONCatalogWithMessages(language.Polish, "", make(JSONMessages))
		assert.NoError(t, err)
	})

	t.Run("Updates plural rules", func(t *testing.T) {
		srcCatl := NewJSONCatalog(language.English, "").(*jsonCatalog)
		srcCatl.mustSetMessage("car", &JSONMessage{
			Translations: map[cldrplural.Category]string{
				cldrplural.One:   "Car",
				cldrplural.Other: "Cars",
			},
		})

		cpy, err := NewJSONCatalogWithMessages(language.Polish, "", srcCatl.Messages())
		require.NoError(t, err)
		catl := cpy.(*jsonCatalog)

		srcTranslations := srcCatl.lookupMap[""]["car"].Translations
		translations := catl.lookupMap[""]["car"].Translations

		assert.Contains(t, srcTranslations, cldrplural.One)
		assert.NotContains(t, srcTranslations, cldrplural.Few)
		assert.NotContains(t, srcTranslations, cldrplural.Many)
		assert.Contains(t, srcTranslations, cldrplural.Other)

		assert.Contains(t, translations, cldrplural.One)
		assert.Contains(t, translations, cldrplural.Few)
		assert.Contains(t, translations, cldrplural.Many)
		assert.Contains(t, translations, cldrplural.Other)

		cpy, err = NewJSONCatalogWithMessages(language.English, "", catl.Messages())
		require.NoError(t, err)
		catl = cpy.(*jsonCatalog)

		translations = catl.lookupMap[""]["car"].Translations
		assert.Contains(t, translations, cldrplural.One)
		assert.Contains(t, translations, cldrplural.Other)
	})

	t.Run("Returns error in invalid data", func(t *testing.T) {
		messages := make(JSONMessages)
		messages["key"] = nil

		_, err := NewJSONCatalogWithMessages(language.English, "", messages)
		assert.Error(t, err)

		delete(messages, "key")
		messages[""] = &JSONMessage{Translations: map[cldrplural.Category]string{cldrplural.Other: "Car"}}
		_, err = NewJSONCatalogWithMessages(language.English, "", messages)
		assert.Error(t, err)
	})
}

func Test_jsonCatalog_setMessage(t *testing.T) {
	catl := NewJSONCatalog(language.English, "").(*jsonCatalog)
	validMsg := &JSONMessage{
		Translations: map[cldrplural.Category]string{
			cldrplural.One:   "Car",
			cldrplural.Other: "Cars",
		},
	}

	t.Run("empty message id", func(t *testing.T) {
		assert.Error(t, catl.setMessage("", validMsg))
	})

	t.Run("nil message", func(t *testing.T) {
		assert.Error(t, catl.setMessage("", nil))
	})

	t.Run("missiong plural rule other", func(t *testing.T) {
		msg := &JSONMessage{
			Translations: map[cldrplural.Category]string{
				cldrplural.One: "Car",
			},
		}
		assert.Error(t, catl.setMessage("valid_key", msg))
	})

	t.Run("invalid context", func(t *testing.T) {
		msg := &JSONMessage{
			Context: "ctx",
			Translations: map[cldrplural.Category]string{
				cldrplural.One:   "Car",
				cldrplural.Other: "Cars",
			},
		}

		assert.Error(t, catl.setMessage("valid_key", msg))
		assert.NoError(t, catl.setMessage("valid_ctx", msg))
	})

	t.Run("dont edit message", func(t *testing.T) {
		msg := &JSONMessage{
			Translations: map[cldrplural.Category]string{
				cldrplural.Other: "Car",
			},
		}
		assert.NoError(t, catl.setMessage("car", msg))

		msg.Translations[cldrplural.One] = "One car"
		catMsg := catl.lookupMap[""]["car"]
		assert.NotContains(t, catMsg.Translations, cldrplural.One)

		catMsg.Translations[cldrplural.Few] = "Few cars"
		assert.NotContains(t, msg.Translations, cldrplural.Few)

		msg.Context = "ctx"
		assert.Empty(t, catMsg.Context)
	})
}

func TestApplyPluralCategoriesToJSONMessage(t *testing.T) {
	t.Run("dont panic with invalid input", func(t *testing.T) {
		assert.NotPanics(t, func() {
			ApplyPluralCategoriesToJSONMessage(language.English, nil)
			ApplyPluralCategoriesToJSONMessage(language.English, &JSONMessage{})
		})
	})

	t.Run("add categories", func(t *testing.T) {
		msg := &JSONMessage{
			Translations: map[cldrplural.Category]string{
				cldrplural.One:   "Car",
				cldrplural.Other: "Cars",
			},
		}

		ApplyPluralCategoriesToJSONMessage(language.Polish, msg)
		assert.Contains(t, msg.Translations, cldrplural.One)
		assert.Contains(t, msg.Translations, cldrplural.Few)
		assert.Contains(t, msg.Translations, cldrplural.Many)
		assert.Contains(t, msg.Translations, cldrplural.Other)
	})
}
