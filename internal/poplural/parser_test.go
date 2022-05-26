package poplural

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fixture struct {
	PluralForm string
	Fixture    []int
}

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

	t.Run("test fixtures", func(t *testing.T) {
		f, err := os.Open("../../testdata/pluralfixtures.json")
		if err != nil {
			t.Fatal(err)
		}
		dec := json.NewDecoder(f)
		var fixtures []fixture
		err = dec.Decode(&fixtures)
		if err != nil {
			t.Fatal(err)
		}

		for _, data := range fixtures {
			forms, err := Parse(data.PluralForm)
			if err != nil {
				t.Errorf("'%s' triggered error: %s", data.PluralForm, err)
			} else if forms == nil || forms.node == nil {
				t.Logf("'%s' compiled to nil", data.PluralForm)
				t.Fail()
			} else {
				for n, e := range data.Fixture {
					i := forms.IndexForN(n)
					if i != e {
						t.Logf("'%s' with n = %d, expected %d, got %d, compiled to", data.PluralForm, n, e, i)
						t.Fail()
					}
					if i == -1 {
						break
					}
				}
			}
		}
	})
}

func TestForms_IndexForN(t *testing.T) {
	t.Run("test nil return zero", func(t *testing.T) {
		f := MustParse("nplurals=2; plural=(n != 1);")
		assert.NotNil(t, f)

		f.node = nil
		assert.Zero(t, f.IndexForN(-1))
		assert.Zero(t, f.IndexForN(0))
		assert.Zero(t, f.IndexForN(1))
		assert.Zero(t, f.IndexForN(10))
	})

	t.Run("test invalid input", func(t *testing.T) {
		f := MustParse("nplurals=2; plural=(n != 1);")
		assert.NotNil(t, f)

		f.node = nil
		assert.Zero(t, f.IndexForN("test"))
		assert.Zero(t, f.IndexForN([]string{}))
		assert.Zero(t, f.IndexForN(nil))
	})

	t.Run("nplurals=1", func(t *testing.T) {
		f := MustParse("nplurals=1; plural=0;")
		require.NotNil(t, f)

		for i := -100; i <= 100; i++ {
			assert.Zero(t, f.IndexForN(i))
		}
	})
}

func TestMustParse(t *testing.T) {
	f := func() {
		MustParse("invalid plural formst")
	}
	assert.Panics(t, f)
}
