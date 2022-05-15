package plural

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fixture struct {
	PluralForm string
	Fixture    []int
}

func TestParser_Parse(t *testing.T) {

	t.Run("nplurals prefix", func(t *testing.T) {
		f, err := Parse("=2; plural=(n != 1);")
		assert.Error(t, err)
		assert.Nil(t, f)

		f, err = Parse("nplurals 2; plural=(n != 1);")
		assert.Error(t, err)
		assert.Nil(t, f)
	})

	t.Run("plural count prefix", func(t *testing.T) {
		f, err := Parse("nplurals= ; plural=(n != 1);")
		assert.Error(t, err)
		assert.Nil(t, f)

		f, err = Parse("nplurals=t; plural=(n != 1);")
		assert.Error(t, err)
		assert.Nil(t, f)
	})

	t.Run("semicolon prefix", func(t *testing.T) {
		f, err := Parse("nplurals=2 plural=(n != 1);")
		assert.Error(t, err)
		assert.Nil(t, f)
	})

	t.Run("test semicolon suffix", func(t *testing.T) {
		f, err := Parse("nplurals=2 plural=(n != 1)")
		assert.Error(t, err)
		assert.Nil(t, f)
	})

	t.Run("test open bracket", func(t *testing.T) {
		f, err := Parse("nplurals=2 plural=(n != 1")
		assert.Error(t, err)
		assert.Nil(t, f)
	})

	t.Run("test invalid compare sign", func(t *testing.T) {
		f, err := Parse("nplurals=2 plural=(n !!= 1")
		assert.Error(t, err)
		assert.Nil(t, f)
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
}

func TestMustParse(t *testing.T) {
	f := func() {
		MustParse("invalid plural formst")
	}
	assert.Panics(t, f)
}
