package locale

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vorlif/spreak/catalog/po"
)

func TestParsePo(t *testing.T) {

	errW := filepath.WalkDir("./", func(path string, d fs.DirEntry, err error) error {
		assert.NoError(t, err)
		if filepath.Ext(path) != ".po" {
			return nil
		}

		data, errR := os.ReadFile(path)
		if assert.NoError(t, errR) {
			assert.NotNil(t, data)
			assert.NotEmpty(t, data)

			poFile, errP := po.Parse(data)
			assert.NoError(t, errP)
			assert.NotNil(t, poFile)
			assert.NotEmpty(t, poFile.Messages)
		}

		return nil
	})

	require.NoError(t, errW)

}
