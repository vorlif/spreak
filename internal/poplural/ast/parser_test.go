package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {

	t.Run("test invalid forms plural", func(t *testing.T) {
		tests := []struct {
			name string
			form string
		}{
			{"missing nplurals prefix", "=2; plural=(n != 1);"},
			{"missing nplurals assign", "nplurals 2; plural=(n != 1);"},
			{"first semicolon missing", "nplurals=2 plural=(n != 1);"},
			{"missing semicolon missing", "nplurals=2 plural=(n != 1)"},
			{"missing nplurals count", "nplurals= ; plural=(n != 1);"},
			{"invalid nplurals count", "nplurals=t; plural=(n != 1);"},
			{"open bracket right", "nplurals=2 plural=(n != 1;"},
			{"open bracket left", "nplurals=2; plural=n != 1);"},
			{"invalid neq sign", "nplurals=2; plural=(n !!= 1);"},
			{"missing plural", "nplurals=2; =(n != 1);"},
			{"missing plural assign", "nplurals=2; plural(n != 1);"},
			{"invalid variable", "nplurals=2; plural=(a != 1);"},
			{"invalid number", "nplurals=2; plural=(n != i);"},
			{"invalid sign", "nplurals=2; plural=(n != 1)#;"},
			{"invalid suffix", "nplurals=2; plural=(n != 1); more"},
			{"missing colon", "nplurals=3; plural=(n==1 ? 0 n==2 ? 1 : 2);"},
		}
		//
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				f, err := Parse(tt.form)
				assert.Error(t, err)
				assert.Nil(t, f)
			})
		}

	})
}

func TestMustParse(t *testing.T) {
	f := func() {
		MustParse("invalid plural formst")
	}
	assert.Panics(t, f)
}
