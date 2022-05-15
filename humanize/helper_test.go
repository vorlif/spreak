package humanize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHumanizer_LanguageName(t *testing.T) {
	h := createGermanHumanizer(t)
	assert.Equal(t, "Deutsch", h.LanguageName("German"))
}
