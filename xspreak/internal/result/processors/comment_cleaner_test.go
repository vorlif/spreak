package processors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRangRe(t *testing.T) {
	tests := []struct {
		text          string
		assertionFunc assert.BoolAssertionFunc
	}{
		{"range: 1..6", assert.True},
		{"range: 100..6", assert.True},
		{"range: 1..600", assert.True},
		{"range: 1..6  ", assert.True},
		{"range:   1..6  ", assert.True},
		{"range: 1...6", assert.False},
		{"range: a..6", assert.False},
		{"range: a..6 bb", assert.False},
	}

	for _, tt := range tests {
		tt.assertionFunc(t, reRange.MatchString(tt.text))
	}
}
