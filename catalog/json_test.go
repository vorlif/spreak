package catalog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

var jsonTestFile = filepath.FromSlash("../testdata/translation-test/json/en.json")
var deJsonTestFile = filepath.FromSlash("../testdata/translation-test/json/de.json")

func TestJsonMessage_MarshalJSON(t *testing.T) {
	msg := &JSONMessage{Other: "test"}

	data, err := json.Marshal(msg)
	assert.NoError(t, err)
	require.NotNil(t, data)
	assert.Equal(t, `"test"`, string(data))

	msg.Context = "ctx"
	data, err = json.Marshal(msg)
	assert.NoError(t, err)
	require.NotNil(t, data)
	assert.Equal(t, `{"context":"ctx","other":"test"}`, string(data))

	data, err = json.MarshalIndent(msg, "", "\t")
	assert.NoError(t, err)
	require.NotNil(t, data)
	assert.Equal(t, `{
	"context": "ctx",
	"other": "test"
}`, string(data))
}

func TestJsonMessage_UnmarshalJSON(t *testing.T) {
	var msg JSONMessage
	err := json.Unmarshal([]byte(`"test"`), &msg)
	assert.NoError(t, err)
	assert.Equal(t, msg.Other, "test")

	msg = JSONMessage{}
	err = json.Unmarshal([]byte(`{"context":"ctx","other":"test"}`), &msg)
	assert.NoError(t, err)
	assert.Equal(t, msg.Other, "test")
	assert.Equal(t, msg.Context, "ctx")
}

func TestJsonDecoder_Decode(t *testing.T) {
	t.Run("returns error on invalid json", func(t *testing.T) {
		cat, errC := NewJSONDecoder().Decode(language.English, "a", []byte("invalid json"))
		assert.Error(t, errC)
		require.Nil(t, cat)
	})

	t.Run("returns error on empty file", func(t *testing.T) {
		cat, errC := NewJSONDecoder().Decode(language.English, "a", []byte(""))
		assert.Error(t, errC)
		require.Nil(t, cat)

		cat, errC = NewJSONDecoder().Decode(language.English, "a", []byte("{}"))
		assert.Error(t, errC)
		require.Nil(t, cat)
	})

	t.Run("returns no error on valid input", func(t *testing.T) {
		data, err := os.ReadFile(jsonTestFile)
		assert.NoError(t, err)
		require.NotNil(t, data)

		cat, errC := NewJSONDecoder().Decode(language.English, "a", data)
		assert.NoError(t, errC)
		require.NotNil(t, cat)
	})
}

func TestJSONCatalog(t *testing.T) {
	t.Run("test translation lookup", func(t *testing.T) {
		data, err := os.ReadFile(jsonTestFile)
		assert.NoError(t, err)
		require.NotNil(t, data)

		cat, errC := NewJSONDecoder().Decode(language.English, "a", data)
		assert.NoError(t, errC)
		require.NotNil(t, cat)

		assert.Equal(t, language.English, cat.Language())

		tr, err := cat.GetTranslation("", "app.name")
		assert.NoError(t, err)
		assert.Equal(t, "TODO List", tr)

		tr, err = cat.GetPluralTranslation("", "animal.cat", 1)
		assert.NoError(t, err)
		assert.Equal(t, "I do not have a cat", tr)

		tr, err = cat.GetPluralTranslation("", "animal.cat", 2)
		assert.NoError(t, err)
		assert.Equal(t, "I do not have cats", tr)

		tr, err = cat.GetPluralTranslation("my-animals", "animal.dog", 2)
		assert.NoError(t, err)
		assert.Equal(t, "I have dogs", tr)
	})

	t.Run("test translation lookup errors", func(t *testing.T) {
		data, err := os.ReadFile(deJsonTestFile)
		assert.NoError(t, err)
		require.NotNil(t, data)

		cat, errC := NewJSONDecoder().Decode(language.German, "a", data)
		assert.NoError(t, errC)
		require.NotNil(t, cat)

		assert.Equal(t, language.German, cat.Language())

		tr, err := cat.GetTranslation("", "animal.cat")
		assert.Error(t, err)
		assert.Equal(t, "animal.cat", tr)

		tr, err = cat.GetPluralTranslation("", "animal.cat", 1)
		assert.Error(t, err)
		assert.Equal(t, "animal.cat", tr)

		tr, err = cat.GetTranslation("unknown", "animal.dog")
		assert.Error(t, err)
		assert.Equal(t, "animal.dog", tr)

		tr, err = cat.GetPluralTranslation("", "missing.plural", 1)
		assert.Error(t, err)
		assert.Equal(t, "missing.plural", tr)
	})
}
