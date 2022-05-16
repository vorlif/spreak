package humanize

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestHumanizer_LanguageName(t *testing.T) {
	h := createGermanHumanizer(t)
	assert.Equal(t, "Deutsch", h.LanguageName("German"))
}

func TestToFixed(t *testing.T) {
	tests := []struct {
		num       float64
		precision int
		expected  float64
	}{
		{2.234567890, 2, 2.23},
		{2.234567890, 3, 2.235},
		{5.2222, 2, 5.22},
		{5.2252, 2, 5.23},
		{5.2249, 2, 5.22},
		{99999999999999.234567890, 3, 99999999999999.235},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v %d", tt.num, tt.precision), func(t *testing.T) {
			assert.Equal(t, tt.expected, toFixed(tt.num, tt.precision))
		})
	}
}

func TestHumanizer_Language(t *testing.T) {
	h := createGermanHumanizer(t)
	assert.Equal(t, language.German, h.Language())

	h = createSourceHumanizer(t)
	assert.Equal(t, language.English, h.Language())

	p := createNewParcel(t)
	h = p.CreateHumanizer(language.Und)
	assert.Equal(t, language.English, h.Language())
}

func TestHumanizer_FilesizeFormat(t *testing.T) {
	t.Run("source language", func(t *testing.T) {
		tests := []struct {
			size     int64
			expected string
		}{
			{1023, "1,023 bytes"},
			{-1023, "-1,023 bytes"},
			{1024, "1 KB"},
			{10 * 1024, "10 KB"},
			{1024*1024 - 1, "1,024 KB"},
			{1024 * 1024, "1 MB"},
			{1024 * 1024 * 50, "50 MB"},
			{1024*1024*1024 - 1, "1,024 MB"},
			{1024 * 1024 * 1024, "1 GB"},
			{1024 * 1024 * 1024 * 1024, "1 TB"},
			{1024 * 1024 * 1024 * 1024 * 1024, "1 PB"},
			{1024 * 1024 * 1024 * 1024 * 1024 * 2000, "2,000 PB"},
		}

		h := createSourceHumanizer(t)
		for _, tt := range tests {
			t.Run(tt.expected, func(t *testing.T) {
				assert.Equal(t, tt.expected, h.FilesizeFormat(tt.size))
			})
		}
	})

	t.Run("translation", func(t *testing.T) {
		tests := []struct {
			size     interface{}
			expected string
		}{
			{1023, "1.023 Bytes"},
			{1024, "1 KB"},
			{10 * 1024, "10 KB"},
			{1024*1024 - 1, "1.024 KB"},
			{1024 * 1024, "1 MB"},
			{1024 * 1024 * 1.5, "1,5 MB"},
			{1024 * 1024 * 50, "50 MB"},
			{1024*1024*1024 - 1, "1.024 MB"},
			{1024 * 1024 * 1024, "1 GB"},
			{1024 * 1024 * 1024 * 1024, "1 TB"},
			{1024 * 1024 * 1024 * 1024 * 1024, "1 PB"},
			{1024 * 1024 * 1024 * 1024 * 1024 * 2000, "2.000 PB"},
			{make([]byte, 10), "10 Bytes"},
			{make([]byte, 2048), "2 KB"},
		}

		h := createGermanHumanizer(t)
		for _, tt := range tests {
			t.Run(tt.expected, func(t *testing.T) {
				assert.Equal(t, tt.expected, h.FilesizeFormat(tt.size))
			})
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		h := createSourceHumanizer(t)
		assert.Equal(t, "<nil>", h.FilesizeFormat(nil))
		assert.Equal(t, "%!(string=test)", h.FilesizeFormat("test"))
	})
}
