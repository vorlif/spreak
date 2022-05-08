package spreak

import (
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/vorlif/spreak/internal/util"
)

func TestNewFilesystemLoader(t *testing.T) {
	t.Run("error is returned when a nil option is passed", func(t *testing.T) {
		fsLoader, err := NewFilesystemLoader(WithLoadPath(testdataStructureDir), nil)
		require.Error(t, err)
		require.Nil(t, fsLoader)
	})

	t.Run("error is returned when no filesystem is passed", func(t *testing.T) {
		fsLoader, err := NewFilesystemLoader()
		require.Error(t, err)
		require.Nil(t, fsLoader)
	})
}

func TestFilesystemLoader_Load(t *testing.T) {
	t.Run("failure when opening a file returns an error", func(t *testing.T) {
		reducer := &testReducer{
			f: func(fsys fs.FS, extension string, lang language.Tag, domain string) (string, error) {
				return "/failure.po", nil
			},
		}
		fsys := &testFs{f: func(name string) (fs.File, error) {
			return nil, os.ErrPermission
		}}
		fsLoader, err := NewFilesystemLoader(
			WithLoaderFs(fsys),
			WithReducer(reducer),
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

		reducer, errReducer := NewDefaultReducer(WithDisabledSearch(), WithCategory(testCategory))
		require.NoError(t, errReducer)

		fsLoader, err := NewFilesystemLoader(
			WithLoaderFs(util.DirFS(testdataStructureDir)),
			WithReducer(reducer),
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
			false, "de/my_category/b.po", PoFile,
		},
		{
			language.German, "a", "",
			false, "de/LC_MESSAGES/a.po", PoFile,
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
		reducer, errR := NewDefaultReducer(WithCategory(tt.category))
		require.NoError(t, errR)
		require.NotNil(t, reducer)
		path, err := reducer.Reduce(fsys, tt.extension, tt.lang, tt.domain)
		if tt.wantErr {
			assert.Error(t, err)
			continue
		}

		if assert.NoError(t, err, "Reduce(... %v %v %v %v...", tt.lang, tt.domain, tt.category, tt.extension) {
			assert.Equal(t, tt.wantPath, path)
		}
	}

}

func TestReduceMoFiles(t *testing.T) {
	reducer, errR := NewDefaultReducer()

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
		path, err := reducer.Reduce(fsys, tt.extension, tt.lang, tt.domain)
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
			false, "de/LC_MESSAGES/a.po", PoFile,
		},
		{
			language.German, "b", "my_category",
			false, "de/my_category/b.po", PoFile,
		},
		{
			language.English, "domain", "cat",
			true, "", UnknownFile,
		},
	}

	for idx, tt := range tests {
		reducer, errR := NewDefaultReducer(WithDisabledSearch(), WithCategory(tt.category))
		require.NoError(t, errR)
		require.NotNil(t, reducer)

		path, err := reducer.Reduce(fsys, tt.extension, tt.lang, tt.domain)
		if tt.wantErr {
			assert.Error(t, err, idx)
			continue
		}

		if assert.NoError(t, err, idx) {
			assert.Equal(t, tt.wantPath, path)
		}
	}
}

func TestNewDefaultReducer(t *testing.T) {
	reducer, err := NewDefaultReducer(WithDisabledSearch())
	if assert.NoError(t, err) {
		require.NotNil(t, reducer)
	}

	reducer, err = NewDefaultReducer()
	if assert.NoError(t, err) {
		require.NotNil(t, reducer)
	}
}
