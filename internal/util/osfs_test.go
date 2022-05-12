package util

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirFS(t *testing.T) {
	dir, err := os.MkdirTemp("", "spreakTest")
	require.NoError(t, err)

	file, err := os.CreateTemp(dir, "myname.*.bat")
	require.NoError(t, err)
	defer file.Close()
	defer os.Remove(file.Name())

	fs := DirFS("")
	readDir, err := fs.ReadDir(dir)
	if assert.NoError(t, err) {
		assert.Len(t, readDir, 1)
		assert.False(t, readDir[0].IsDir())
		assert.Equal(t, readDir[0].Name(), filepath.Base(file.Name()))
	}

	stat, err := fs.Stat(file.Name())
	assert.NoError(t, err)
	assert.False(t, stat.IsDir())
	assert.Zero(t, stat.Size())

	content := "this is a spreak test"
	_, err = file.WriteString(content)
	assert.NoError(t, err)
	assert.NoError(t, file.Close())

	read, errR := fs.ReadFile(file.Name())
	if assert.NoError(t, errR) {
		assert.Equal(t, content, string(read))
	}

	fs = DirFS(dir)
	stat, err = fs.Stat("unknown")
	assert.Error(t, err)
	assert.Nil(t, stat)

	file2, err := fs.Open(filepath.Base(file.Name()))
	assert.NoError(t, err)
	assert.NoError(t, file2.Close())
}
