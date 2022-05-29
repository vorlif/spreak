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

func TestJsonDecoder(t *testing.T) {
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
}
