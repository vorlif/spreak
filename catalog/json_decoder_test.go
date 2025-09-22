package catalog

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

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
		data, err := os.ReadFile(enJSONTestFile)
		assert.NoError(t, err)
		require.NotNil(t, data)

		cat, errC := NewJSONDecoder().Decode(language.English, "a", data)
		assert.NoError(t, errC)
		require.NotNil(t, cat)
	})
}
