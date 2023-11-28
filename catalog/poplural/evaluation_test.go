package poplural

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vorlif/spreak/catalog/poplural/ast"
)

type fixture struct {
	PluralForm string
	Fixture    []int
}

func TestWithFixtures(t *testing.T) {
	t.Run("test fixtures", func(t *testing.T) {
		f, err := os.Open(filepath.FromSlash("../../testdata/pluralfixtures.json"))
		require.NoError(t, err)
		defer f.Close()

		dec := json.NewDecoder(f)
		var fixtures []fixture
		err = dec.Decode(&fixtures)
		require.NoError(t, err)

		builtInCount := len(rawToBuiltIn)

		for _, data := range fixtures {
			parsed, err := ast.Parse(data.PluralForm)
			require.NoError(t, err)
			require.NotNil(t, parsed)

			formFunc := generateFormFunc(parsed)

			require.NotNil(t, formFunc)

			builtInForm := rawToBuiltIn[ast.CompileToString(parsed)]
			builtInCount--
			require.NotNil(t, builtInForm, data.PluralForm)
			require.NotNil(t, builtInForm.FormFunc)

			for input, want := range data.Fixture {
				assert.Equalf(t, want, formFunc(int64(input)), "%s form.FormFunc(%d) = %d", data.PluralForm, input, want)
				assert.Equalf(t, want, builtInForm.FormFunc(int64(input)), "%s builtInForm.FormFunc(%d) = %d", data.PluralForm, input, want)
			}
		}

		assert.Zero(t, builtInCount)
	})
}

func TestMustParse(t *testing.T) {
	t.Run("panics on invalid plural form", func(t *testing.T) {
		f := func() { MustParse("zzz") }
		assert.Panics(t, f)
	})
}

func TestParse(t *testing.T) {
	t.Run("error on invalid plural form", func(t *testing.T) {
		f, err := Parse("zzz")
		assert.Error(t, err)
		assert.Nil(t, f)
	})

	t.Run("unknown rule evaluates at runtime", func(t *testing.T) {
		f, err := Parse("nplurals=2; plural=n > 12;")
		assert.NoError(t, err)
		require.NotNil(t, f)
		assert.Equal(t, 2, f.NPlurals)
		assert.Equal(t, 0, f.FormFunc(1))
		assert.Equal(t, 1, f.FormFunc(13))
	})
}

func TestEvaluate(t *testing.T) {
	t.Run("test nil return zero", func(t *testing.T) {
		f := ast.MustParse("nplurals=2; plural=(n != 1);")
		assert.NotNil(t, f)

		f.Root = nil
		assert.Zero(t, generateFormFunc(f)(-1))
		assert.Zero(t, generateFormFunc(f)(0))
		assert.Zero(t, generateFormFunc(f)(1))
		assert.Zero(t, generateFormFunc(f)(10))
	})

	t.Run("test invalid input", func(t *testing.T) {
		f := MustParse("nplurals=2; plural=(n != 1);")
		assert.NotNil(t, f)

		assert.Zero(t, f.Evaluate("test"))
		assert.Zero(t, f.Evaluate([]string{}))
		assert.Zero(t, f.Evaluate(nil))
	})

	t.Run("nplurals=1", func(t *testing.T) {
		f := MustParse("nplurals=1; plural=0;")
		require.NotNil(t, f)

		for i := -100; i <= 100; i++ {
			assert.Zero(t, f.Evaluate(i))
		}
	})
}
