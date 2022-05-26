package cldrplural

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractOperands(t *testing.T) {
	tests := []struct {
		src string
		n   float64
		i   int64
		v   int64
		w   int64
		f   int64
		t   int64
		c   int64
	}{
		{"1", 1, 1, 0, 0, 0, 0, 0},
		{"1.", 1, 1, 0, 0, 0, 0, 0},
		{"1.0", 1, 1, 1, 0, 0, 0, 0},
		{"1.00", 1, 1, 2, 0, 0, 0, 0},
		{"1.3", 1.3, 1, 1, 1, 3, 3, 0},
		{"1.30", 1.3, 1, 2, 1, 30, 3, 0},
		{"1.03", 1.03, 1, 2, 2, 3, 3, 0},
		{"1.230", 1.23, 1, 3, 2, 230, 23, 0},
		{"1200000", 1200000, 1200000, 0, 0, 0, 0, 0},
		{"1.2c6", 1200000, 1200000, 0, 0, 0, 0, 6},
		{"123c6", 123000000, 123000000, 0, 0, 0, 0, 6},
		{"123c5", 12300000, 12300000, 0, 0, 0, 0, 5},
		{"1200.50", 1200.5, 1200, 2, 1, 50, 5, 0},
		{"1.20050c3", 1200.5, 1200, 2, 1, 50, 5, 3},
	}

	for _, tt := range tests {
		op := NewOperands(tt.src)
		assert.Equalf(t, tt.n, op.N, "N %s", tt.src)
		assert.Equalf(t, tt.i, op.I, "I %s", tt.src)
		assert.Equalf(t, tt.v, op.V, "V %s", tt.src)
		assert.Equalf(t, tt.w, op.W, "W %s", tt.src)
		assert.Equalf(t, tt.f, op.F, "F %s", tt.src)
		assert.Equalf(t, tt.t, op.T, "T %s", tt.src)
		assert.Equalf(t, tt.c, op.C, "C %s", tt.src)
	}

	assert.Equal(t, -3, -3%5)
	assert.Equal(t, -3, -3%-5)
	assert.Equal(t, 3, 3%5)
}

func Test_formatExponent(t *testing.T) {
	tests := []struct {
		input string
		c     int
		want  string
	}{
		{"1.2", 6, "1200000"},
		{"123", 6, "123000000"},
		{"123", 5, "12300000"},
		{"1.20050", 3, "1200.50"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equalf(t, tt.want, shiftDecimalPoint(tt.input, tt.c), "shiftDecimalPoint(%q, %d) = %s", tt.input, tt.c, tt.want)
		})

	}
}
