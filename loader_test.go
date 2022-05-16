package spreak

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/vorlif/spreak/internal/util"
)

func TestNewFilesystemLoader(t *testing.T) {
	t.Run("error is returned when a nil option is passed", func(t *testing.T) {
		fsLoader, err := NewFilesystemLoader(WithPath(testdataStructureDir), nil)
		require.Error(t, err)
		require.Nil(t, fsLoader)
	})

	t.Run("error is returned when no filesystem is passed", func(t *testing.T) {
		fsLoader, err := NewFilesystemLoader()
		require.Error(t, err)
		require.Nil(t, fsLoader)
	})

	t.Run("error is returned if multiple file systems are to be stored", func(t *testing.T) {
		fsLoader, err := NewFilesystemLoader(WithPath(testdataStructureDir), WithPath(testdataStructureDir))
		require.Error(t, err)
		require.Nil(t, fsLoader)
	})
}

func TestFilesystemLoader_Load(t *testing.T) {
	t.Run("failure when opening a file returns an error", func(t *testing.T) {
		reducer := &testResolver{
			f: func(fsys fs.FS, extension string, lang language.Tag, domain string) (string, error) {
				return "/failure.po", nil
			},
		}
		fsys := &testFs{f: func(name string) (fs.File, error) {
			return nil, os.ErrPermission
		}}
		fsLoader, err := NewFilesystemLoader(
			WithFs(fsys),
			WithResolver(reducer),
		)
		require.NoError(t, err)
		require.NotNil(t, fsLoader)

		data, errLoad := fsLoader.Load(language.English, NoDomain)
		if assert.Error(t, errLoad) {
			assert.Equal(t, os.ErrPermission, errLoad)
		}
		assert.Nil(t, data)
	})
}

func TestWithCategory(t *testing.T) {
	t.Run("when a category is passed, it is used", func(t *testing.T) {
		testCategory := "my_category"

		reducer, errResolver := NewDefaultResolver(WithDisabledSearch(), WithCategory(testCategory))
		require.NoError(t, errResolver)

		fsLoader, err := NewFilesystemLoader(
			WithFs(util.DirFS(testdataStructureDir)),
			WithResolver(reducer),
		)
		require.NoError(t, err)
		require.NotNil(t, fsLoader)

		data, errLoad := fsLoader.Load(language.German, "domain_test")
		assert.Error(t, errLoad)
		assert.Nil(t, data)

		data, errLoad = fsLoader.Load(language.German, "b")
		assert.NoError(t, errLoad)
		assert.NotNil(t, data)
	})
}

func TestLoadPo(t *testing.T) {
	fsys := util.DirFS(testdataStructureDir)
	tests := []struct {
		lang      language.Tag
		domain    string
		category  string
		wantErr   bool
		wantPath  string
		extension string
	}{
		{
			language.German, "b", "my_category",
			false, filepath.Join("de", "my_category", "b.po"), PoFile,
		},
		{
			language.German, "a", "",
			false, filepath.Join("de", "LC_MESSAGES", "a.po"), PoFile,
		},
		{
			language.French, "a", "",
			true, "", UnknownFile,
		},
		{
			language.English, "domain", "cat",
			true, "", UnknownFile,
		},
		{
			language.MustParse("de_AT"), "", "",
			false, "de_AT.po", PoFile,
		},
	}

	for _, tt := range tests {
		reducer, errR := NewDefaultResolver(WithCategory(tt.category))
		require.NoError(t, errR)
		require.NotNil(t, reducer)
		path, err := reducer.Resolve(fsys, tt.extension, tt.lang, tt.domain)
		if tt.wantErr {
			assert.Error(t, err)
			continue
		}

		if assert.NoError(t, err, "Resolve(... %v %v %v %v...", tt.lang, tt.domain, tt.category, tt.extension) {
			assert.Equal(t, tt.wantPath, path)
		}
	}

}

func TestReduceMoFiles(t *testing.T) {
	reducer, errR := NewDefaultResolver()

	require.NoError(t, errR)
	require.NotNil(t, reducer)

	fsys := util.DirFS(testdataStructureDir)
	tests := []struct {
		lang      language.Tag
		domain    string
		category  string
		wantErr   bool
		wantPath  string
		extension string
	}{
		{
			language.German, "b", "my_category",
			true, "", UnknownFile,
		},
		{
			language.French, "a", "",
			true, "", UnknownFile,
		},
		{
			language.English, "domain", "cat",
			false, "en.mo", MoFile,
		},
	}

	for _, tt := range tests {
		path, err := reducer.Resolve(fsys, tt.extension, tt.lang, tt.domain)
		if tt.wantErr {
			assert.Error(t, err)
			continue
		}

		if assert.NoError(t, err) {
			assert.Equal(t, tt.wantPath, path)
		}
	}
}

func TestDisableSearch(t *testing.T) {
	fsys := util.DirFS(testdataStructureDir)
	tests := []struct {
		lang      language.Tag
		domain    string
		category  string
		wantErr   bool
		wantPath  string
		extension string
	}{
		{
			language.German, "other", "my_category",
			true, "", UnknownFile,
		},
		{
			language.German, "a", "LC_MESSAGES",
			false, filepath.Join("de", "LC_MESSAGES", "a.po"), PoFile,
		},
		{
			language.German, "b", "my_category",
			false, filepath.Join("de", "my_category", "b.po"), PoFile,
		},
		{
			language.English, "domain", "cat",
			true, "", UnknownFile,
		},
	}

	for idx, tt := range tests {
		reducer, errR := NewDefaultResolver(WithDisabledSearch(), WithCategory(tt.category))
		require.NoError(t, errR)
		require.NotNil(t, reducer)

		path, err := reducer.Resolve(fsys, tt.extension, tt.lang, tt.domain)
		if tt.wantErr {
			assert.Error(t, err, idx)
			continue
		}

		if assert.NoError(t, err, idx) {
			assert.Equal(t, tt.wantPath, path)
		}
	}
}

func TestNewDefaultResolver(t *testing.T) {
	reducer, err := NewDefaultResolver(WithDisabledSearch())
	if assert.NoError(t, err) {
		require.NotNil(t, reducer)
	}

	reducer, err = NewDefaultResolver()
	if assert.NoError(t, err) {
		require.NotNil(t, reducer)
	}
}

func TestWithDecoder(t *testing.T) {
	t.Run("fallback to default decoders", func(t *testing.T) {
		fl, err := NewFilesystemLoader(WithSystemFs())
		assert.NoError(t, err)
		require.NotNil(t, fl)
		if assert.Len(t, fl.decoder, 2) {
			assert.IsType(t, (*poDecoder)(nil), fl.decoder[0])
			assert.IsType(t, (*moDecoder)(nil), fl.decoder[1])
		}

	})

	t.Run("WithPoDecoder disables fallback", func(t *testing.T) {
		fl, err := NewFilesystemLoader(
			WithPoDecoder(),
			WithSystemFs(),
		)
		assert.NoError(t, err)
		require.NotNil(t, fl)
		assert.Contains(t, fl.extensions, PoFile)
		if assert.Len(t, fl.decoder, 1) {
			assert.IsType(t, (*poDecoder)(nil), fl.decoder[0])
		}
	})

	t.Run("WithMoDecoder disables fallback", func(t *testing.T) {
		fl, err := NewFilesystemLoader(
			WithMoDecoder(),
			WithSystemFs(),
		)
		assert.NoError(t, err)
		require.NotNil(t, fl)
		assert.Contains(t, fl.extensions, MoFile)
		if assert.Len(t, fl.decoder, 1) {
			assert.IsType(t, (*moDecoder)(nil), fl.decoder[0])
		}
	})

	t.Run("WithDecoder sets decoder", func(t *testing.T) {
		ext := ".json"
		dec := &testDecoder{}

		fl, err := NewFilesystemLoader(WithDecoder(ext, dec), WithSystemFs())
		assert.NoError(t, err)
		require.NotNil(t, fl)

		assert.Contains(t, fl.extensions, ext)
		assert.Contains(t, fl.decoder, dec)
	})
}

func TestWithFs(t *testing.T) {
	t.Run("WithFs sets fs", func(t *testing.T) {
		fsys := &testFs{}
		fl, err := NewFilesystemLoader(WithFs(fsys))
		assert.NoError(t, err)
		require.NotNil(t, fl)
		assert.Equal(t, fl.fsys, fsys)
	})

	t.Run("multiple fs returns error", func(t *testing.T) {
		fsys := &testFs{}
		fl, err := NewFilesystemLoader(WithFs(fsys), WithFs(fsys))
		assert.Error(t, err)
		require.Nil(t, fl)
	})
}

func TestWithResolver(t *testing.T) {
	t.Run("WithResolver sets resolver", func(t *testing.T) {
		resolver := &testResolver{}
		fl, err := NewFilesystemLoader(WithSystemFs(), WithResolver(resolver))
		assert.NoError(t, err)
		require.NotNil(t, fl)
		assert.Equal(t, fl.resolver, resolver)
	})

	t.Run("multiple resolver returns error", func(t *testing.T) {
		resolver := &testResolver{}
		fl, err := NewFilesystemLoader(WithSystemFs(),
			WithResolver(resolver),
			WithResolver(resolver),
		)
		assert.Error(t, err)
		require.Nil(t, fl)
	})
}
