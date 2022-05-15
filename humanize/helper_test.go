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
