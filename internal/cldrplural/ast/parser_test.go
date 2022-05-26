package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testParse(rule string) (*Rule, error) {
	p := parser{s: newScanner(rule)}
	return p.Parse()
}

func depth(n Node) int {
	maxDepth := 0
	currentDepth := 0
	Inspect(n, func(node Node) bool {
		if node == nil {
			currentDepth = 0
			return true
		}

		currentDepth++
		if currentDepth > maxDepth {
			maxDepth = currentDepth
		}

		return true
	})

	return maxDepth
}

func TestParse(t *testing.T) {
	t.Run("simplest case", func(t *testing.T) {
		r, err := testParse("n = 1")
		assert.NoError(t, err)
		require.NotNil(t, r)
		require.NotNil(t, r.Root)
		assert.Equal(t, Equal, r.Root.Type())
		if assert.IsType(t, (*InRelationExpr)(nil), r.Root) {
			e := r.Root.(*InRelationExpr)
			assert.Equal(t, Operand, e.X.Type())
			assert.Equal(t, RangeList, e.Y.Type())
		}
	})

	t.Run("range", func(t *testing.T) {
		r, err := testParse("n = 5..10")
		assert.NoError(t, err)
		require.NotNil(t, r)
		require.NotNil(t, r.Root)
		assert.Equal(t, Equal, r.Root.Type())
		if assert.IsType(t, (*InRelationExpr)(nil), r.Root) {
			e := r.Root.(*InRelationExpr)
			assert.Equal(t, Operand, e.X.Type())
			require.Equal(t, RangeList, e.Y.Type())
			rangeExpr := e.Y.X.(*RangeExpr)
			assert.Nil(t, e.Y.Y)
			assert.Equal(t, int64(5), rangeExpr.From)
			assert.Equal(t, int64(10), rangeExpr.To)
		}
	})

	t.Run("or", func(t *testing.T) {
		rule, err := testParse("i = 0 or n = 1")
		assert.NoError(t, err)
		require.NotNil(t, rule)
		require.NotNil(t, rule.Root)
		assert.Equal(t, Or, rule.Root.Type())
		assert.Equal(t, 3, depth(rule.Root))
	})

	t.Run("and", func(t *testing.T) {
		rule, err := testParse("i = 1 and v = 0")
		assert.NoError(t, err)
		require.NotNil(t, rule)
		require.NotNil(t, rule.Root)
		assert.Equal(t, And, rule.Root.Type())
		assert.Equal(t, 3, depth(rule.Root))
	})

	t.Run("modulo", func(t *testing.T) {
		rule, err := testParse("n % 10 = 1")
		assert.NoError(t, err)
		require.NotNil(t, rule)
		require.NotNil(t, rule.Root)
		if assert.Equal(t, Equal, rule.Root.Type()) {
			n := rule.Root.(*InRelationExpr)
			assert.Equal(t, Remainder, n.X.Type())
			assert.Equal(t, RangeList, n.Y.Type())
		}
	})

	t.Run("not equal", func(t *testing.T) {
		rule, err := testParse("n != 1")
		assert.NoError(t, err)
		require.NotNil(t, rule)
		require.NotNil(t, rule.Root)
		if assert.Equal(t, NotEqual, rule.Root.Type()) {
			n := rule.Root.(*InRelationExpr)
			assert.Equal(t, Operand, n.X.Type())
			assert.Equal(t, RangeList, n.Y.Type())
			assert.Nil(t, n.Y.Y)
			valE := n.Y.X.(*ValueExpr)
			assert.Equal(t, int64(1), valE.Value)
		}
	})

	t.Run("must parse", func(t *testing.T) {
		tests := []string{
			"n % 10 = 2..4 and n % 100 != 12..14",
			"n % 10 = 1 and n % 100 != 11,71,91",
			"n % 10 = 3..4,9 and t % 100 != 10..19,70..79,90..99",
			"@integer 2~17, 100, 1000",
		}

		for _, tt := range tests {
			a, err := testParse(tt)
			assert.NoError(t, err, tt)
			assert.NotNil(t, a, tt)
		}
	})

	t.Run("should fail", func(t *testing.T) {
		tests := []string{
			"n % n = 10", // missing constant
		}

		for _, tt := range tests {
			a, err := testParse(tt)
			assert.Error(t, err, tt)
			assert.Nil(t, a, tt)
		}
	})

	t.Run("sample", func(t *testing.T) {
		a, err := testParse("i = 0..1 @integer 0, 1 ... @decimal 0.0~0.5, 8.0 …")
		assert.NoError(t, err)
		require.NotNil(t, a)
		require.NotNil(t, a.Samples)
		assert.Len(t, a.Samples, 9)
		assert.ElementsMatch(t, a.Samples, []string{"0", "1", "0.0", "0.1", "0.2", "0.3", "0.4", "0.5", "8.0"})
	})
}

func Test_increment(t *testing.T) {
	tests := []struct {
		value string
		want  string
	}{
		{"1.1", "1.2"},
		{"0.9", "1.0"},
		{"2", "3"},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			assert.Equalf(t, tt.want, increment(tt.value), "increment(%v)", tt.value)
		})
	}
}

func Test_parseSampleRange(t *testing.T) {
	tests := []struct {
		sampleRange string
		want        []string
	}{
		{"2~4", []string{"2", "3", "4"}},
		{"3~5", []string{"3", "4", "5"}},
		{"103~105", []string{"103", "104", "105"}},
		{"1.3~1.5", []string{"1.3", "1.4", "1.5"}},
	}
	for _, tt := range tests {
		t.Run(tt.sampleRange, func(t *testing.T) {
			assert.Equalf(t, tt.want, parseSampleRange(tt.sampleRange), "parseSampleRange(%v)", tt.sampleRange)
		})
	}
}
