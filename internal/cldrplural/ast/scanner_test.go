package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func collectScan(s *scanner) ([]Token, []string) {
	var tokens []Token
	var literals []string
	for tok, lit := s.Scan(); tok != eof; tok, lit = s.Scan() {
		tokens = append(tokens, tok)
		literals = append(literals, lit)
	}
	return tokens, literals
}

func TestScanner(t *testing.T) {
	t.Run("base cases", func(t *testing.T) {
		s := newScanner("n = 1 @integer 1 @decimal 1.0, 1.00, 1.000, 1.0000")
		toks, lits := collectScan(s)
		wantToks := []Token{Operand, Equal, Value, sample, Value, sample, sampleValue, sampleValue, sampleValue, sampleValue}
		wantLits := []string{"n", "=", "1", "@integer", "1", "@decimal", "1.0", "1.00", "1.000", "1.0000"}
		assert.Equal(t, wantToks, toks)
		assert.Equal(t, wantLits, lits)
	})

	t.Run("value range", func(t *testing.T) {
		tests := []string{"n = 0..1 "}
		for _, tt := range tests {
			t.Run(tt, func(t *testing.T) {
				s := newScanner(tt)
				toks, lits := collectScan(s)
				wantToks := []Token{Operand, Equal, ValueRange}
				wantLits := []string{"n", "=", "0..1"}
				assert.Equal(t, wantToks, toks)
				assert.Equal(t, wantLits, lits)
			})
		}
	})

	t.Run("value range list", func(t *testing.T) {
		s := newScanner("100 != 10..19,70..79,90..99,9")
		toks, lits := collectScan(s)
		wantToks := []Token{Value, NotEqual, ValueRange, ValueRange, ValueRange, Value}
		wantLits := []string{"100", "!=", "10..19", "70..79", "90..99", "9"}
		assert.Equal(t, wantToks, toks)
		assert.Equal(t, wantLits, lits)
	})

	t.Run("and", func(t *testing.T) {
		s := newScanner("n % 10 = 1 and n % 100 != 11")
		toks, lits := collectScan(s)
		wantToks := []Token{Operand, Remainder, Value, Equal, Value, And, Operand, Remainder, Value, NotEqual, Value}
		wantLits := []string{"n", "%", "10", "=", "1", "and", "n", "%", "100", "!=", "11"}
		assert.Equal(t, wantToks, toks)
		assert.Equal(t, wantLits, lits)
	})

	t.Run("or", func(t *testing.T) {
		s := newScanner("i = 0 or n = 1")
		toks, lits := collectScan(s)
		wantToks := []Token{Operand, Equal, Value, Or, Operand, Equal, Value}
		wantLits := []string{"i", "=", "0", "or", "n", "=", "1"}
		assert.Equal(t, wantToks, toks)
		assert.Equal(t, wantLits, lits)
	})

	t.Run("sample range", func(t *testing.T) {
		s := newScanner("@integer 2~17, 100 … @decimal 1.1~2.6, 10.0 …")
		toks, lits := collectScan(s)
		wantToks := []Token{sample, sampleRange, Value, sample, sampleRange, sampleValue}
		wantLits := []string{"@integer", "2~17", "100", "@decimal", "1.1~2.6", "10.0"}
		assert.Equal(t, wantToks, toks)
		assert.Equal(t, wantLits, lits)
	})
}

func TestRegex(t *testing.T) {
	t.Run("sample value", func(t *testing.T) {
		valid := []string{"1", "123.45", "3.340934", "0.45", "0", "1c6", "1.1c6"}
		for _, tt := range valid {
			assert.True(t, reSampleValue.MatchString(tt), tt)
		}

		invalid := []string{".45", "-3", "-3.45", "", "1cc6", "1c6c", "c6", "1.1c6c"}
		for _, tt := range invalid {
			assert.False(t, reSampleValue.MatchString(tt), tt)
		}
	})

	t.Run("range", func(t *testing.T) {
		tests := []struct {
			literal       string
			assertionFunc assert.BoolAssertionFunc
		}{
			{"123..45", assert.True},
			{"0..1", assert.True},
			{"123.45..45", assert.False},
			{"123.45..", assert.False},
			{"123..0", assert.False},
			{"0..123", assert.True},
			{"0", assert.False},
			{"", assert.False},
		}
		for _, tt := range tests {
			tt.assertionFunc(t, reRange.MatchString(tt.literal), tt.literal)
		}
	})

	t.Run("value", func(t *testing.T) {
		tests := []struct {
			literal       string
			assertionFunc assert.BoolAssertionFunc
		}{
			{"1234", assert.True},
			{"0", assert.True},
			{"-1", assert.False},
			{"123.45", assert.False},
			{"0.1", assert.False},
			{"0123", assert.True},
			{"", assert.False},
		}
		for _, tt := range tests {
			tt.assertionFunc(t, reValue.MatchString(tt.literal), tt.literal)
		}
	})

	t.Run("sample range", func(t *testing.T) {
		tests := []struct {
			literal       string
			assertionFunc assert.BoolAssertionFunc
		}{
			{"123~45", assert.True},
			{"0~1", assert.True},
			{"123.45~45", assert.True},
			{"123.45~", assert.False},
			{"123~0", assert.True},
			{"0~123", assert.True},
			{"0.1~123", assert.True},
			{"0~123.5", assert.True},
			{"0.1~123.5", assert.True},
			{"0", assert.False},
			{"", assert.False},
		}
		for _, tt := range tests {
			tt.assertionFunc(t, reSampleRange.MatchString(tt.literal), tt.literal)
		}
	})
}
