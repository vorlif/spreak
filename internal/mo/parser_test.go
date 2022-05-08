package mo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_File(t *testing.T) {
	content, errRead := os.ReadFile("../../testdata/parser/poedit_en_GB.mo")
	require.NoError(t, errRead)

	file, err := ParseBytes(content)
	require.NoError(t, err)
	require.NotNil(t, file)
	require.NotNil(t, file.Header)
	require.NotNil(t, file.Messages)

	assert.Equal(t, "poedit", file.Header.ProjectIDVersion)

}
