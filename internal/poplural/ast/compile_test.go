package ast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompileToString(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"nplurals=2; plural=n     != 1;", "nplurals=2; plural=n != 1;"},
		{"nplurals=2; plural= (n == 0) ? 0 : 1;", "nplurals=2; plural=(n == 0) ? 0 : 1;"},
		{"nplurals=2; plural=n<  5;", "nplurals=2; plural=n < 5;"},
		{"nplurals=2; plural=n   >5;", "nplurals=2; plural=n > 5;"},
		{"nplurals=2; plural=n   >= 5;", "nplurals=2; plural=n >= 5;"},
		{"nplurals=2; plural=  n  >= 5  && 5 < 3 || n != 2;", "nplurals=2; plural=(n >= 5 && 5 < 3 || n != 2);"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s => %s", tt.input, tt.want), func(t *testing.T) {
			parsed := MustParse(tt.input)
			assert.Equal(t, tt.want, CompileToString(parsed))
		})
	}
}
