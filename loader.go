package spreak

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/internal/util"
)

const (
	UnknownFile = "unknown"
	PoFile      = ".po"
	MoFile      = ".mo"
)

// Loader is responsible for loading Catalogs for a language and a domain.
// A bundle loads each domain through its own loader.
type Loader interface {
	Load(lang language.Tag, domain string) (Catalog, error)
}

// A Reducer is used by the FilesystemLoader to find the appropriate path for a file.
// If a file was not found, os.ErrNotFound should be returned.
// All other errors cause the loaders search to stop.
type Reducer interface {
	Reduce(fsys fs.FS, extensions string, lang language.Tag, domain string) (string, error)
}

// FsOption is an option which can be used when creating the FilesystemLoader.
type FsOption func(l *FilesystemLoader) error

// ReducerOption is an option which can be used when creating the DefaultReducer.
type ReducerOption func(r *defaultReducer)

// FilesystemLoader is a Loader which loads and decodes files from a file system.
// A file system here means an implementation of fs.FS.
type FilesystemLoader struct {
	fsys       fs.FS
	reducer    Reducer
	extensions []string
	decoder    []Decoder
}

var _ Loader = (*FilesystemLoader)(nil)

// NewFilesystemLoader creates a new FileSystemLoader.
// If no file system was stored during the creation, an error is returned.
// If no decoder has been stored, the Po and Mo decoders are automatically used.
// Otherwise, only the stored decoders are used.
func NewFilesystemLoader(opts ...FsOption) (*FilesystemLoader, error) {
	l := &FilesystemLoader{
		decoder:    make([]Decoder, 0),
		extensions: make([]string, 0),
	}

	for _, opt := range opts {
		if opt == nil {
			return nil, errors.New("spreak.Loader: try to create an FilesystemLoader with a nil option")
		}
		if err := opt(l); err != nil {
			return nil, err
		}
	}

	if len(l.decoder) == 0 {
		l.addDecoder(PoFile, NewPoDecoder())
		l.addDecoder(MoFile, NewMoDecoder())
	}

	if l.fsys == nil {
		return nil, errors.New("spreak.Loader: try to create an FilesystemLoader without an filesystem")
	}

	if l.reducer == nil {
		reducer, err := NewDefaultReducer()
		if err != nil {
			return nil, err
		}
		l.reducer = reducer
	}

	return l, nil
}

func (l *FilesystemLoader) Load(lang language.Tag, domain string) (Catalog, error) {

	for i, extension := range l.extensions {
		path, errP := l.reducer.Reduce(l.fsys, extension, lang, domain)
		if errP != nil {
			if errors.Is(errP, os.ErrNotExist) {
				continue
			}
			return nil, errP
		}

		f, errF := l.fsys.Open(path)
		if errF != nil {
			if f != nil {
				_ = f.Close()
			}
			return nil, errF
		}
		defer f.Close()

		data, errD := io.ReadAll(f)
		if errD != nil {
			return nil, errD
		}

		catalog, errC := l.decoder[i].Decode(lang, domain, data)
		if errC != nil {
			return nil, errC
		}
		return catalog, nil
	}

	return nil, NewErrNotFound(lang, "file", "domain=%q", domain)
}

func (l *FilesystemLoader) addDecoder(ext string, decoder Decoder) {
	l.extensions = append(l.extensions, ext)
	l.decoder = append(l.decoder, decoder)
}

// WithFs stores a fs.FS as filesystem.
// Lets the creation of the FilesystemLoader fail, if a filesystem was already deposited.
func WithFs(fsys fs.FS) FsOption {
	return func(l *FilesystemLoader) error {
		if l.fsys != nil {
			return errors.New("spreak.Loader: filesystem for FilesystemLoader already set")
		}
		l.fsys = fsys
		return nil
	}
}

// WithPath stores a path as filesystem.
// Lets the creation of the FilesystemLoader fail, if a filesystem was already deposited.
func WithPath(path string) FsOption {
	return func(l *FilesystemLoader) error {
		if l.fsys != nil {
			return errors.New("spreak.Loader: filesystem for FilesystemLoader already set")
		}
		l.fsys = util.DirFS(path)
		return nil
	}
}

// WithSystemFs stores the root path as filesystem.
// Lets the creation of the FilesystemLoader fail, if a filesystem was already deposited.
func WithSystemFs() FsOption { return WithPath("") }

// WithReducer stores the reducer of a FilesystemLoader.
// Lets the creation of the FilesystemLoader fail, if a Reducer was already deposited.
func WithReducer(reducer Reducer) FsOption {
	return func(l *FilesystemLoader) error {
		if l.reducer != nil {
			return errors.New("spreak.Loader: Reducer for FilesystemLoader already set")
		}
		l.reducer = reducer
		return nil
	}
}

// WithDecoder stores a decoder for an extension.
func WithDecoder(ext string, decoder Decoder) FsOption {
	return func(r *FilesystemLoader) error {
		r.addDecoder(ext, decoder)
		return nil
	}
}

// WithMoDecoder stores the mo file decoder.
func WithMoDecoder() FsOption { return WithDecoder(MoFile, &moDecoder{}) }

// WithPoDecoder stores the mo file decoder.
func WithPoDecoder() FsOption { return WithDecoder(PoFile, &poDecoder{}) }

type defaultReducer struct {
	search   bool
	category string
}

// NewDefaultReducer create a reducer which can be used for a FilesystemLoader.
func NewDefaultReducer(opts ...ReducerOption) (Reducer, error) {
	l := &defaultReducer{
		search:   true,
		category: "",
	}

	for _, opt := range opts {
		opt(l)
	}

	return l, nil
}

func WithDisabledSearch() ReducerOption { return func(r *defaultReducer) { r.search = false } }

func WithCategory(category string) ReducerOption {
	return func(l *defaultReducer) {
		l.category = category
	}
}

func (r *defaultReducer) Reduce(fsys fs.FS, extension string, tag language.Tag, domain string) (string, error) {
	for _, lang := range ExpandLanguage(tag) {
		path, err := r.searchFileForLanguageName(fsys, lang, domain, extension)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}

		return path, nil
	}

	return "", os.ErrNotExist
}

func (r *defaultReducer) searchFileForLanguageName(fsys fs.FS, locale, domain, ext string) (string, error) {

	if domain != "" {
		// .../locale/category/domain.mo
		path := filepath.Join(locale, r.category, domain+ext)
		if _, err := fs.Stat(fsys, path); err == nil {
			return path, nil
		}
	}

	if r.search {
		if domain != "" {
			// .../locale/LC_MESSAGES/domain.mo
			path := filepath.Join(locale, "LC_MESSAGES", domain+ext)
			if _, err := fs.Stat(fsys, path); err == nil {
				return path, nil
			}

			// .../locale/domain.mo
			path = filepath.Join(locale, domain+ext)
			if _, err := fs.Stat(fsys, path); err == nil {
				return path, nil
			}

			// .../domain/locale.mo
			path = filepath.Join(domain, locale+ext)
			if _, err := fs.Stat(fsys, path); err == nil {
				return path, nil
			}
		}

		// .../locale.mo
		path := filepath.Join(locale + ext)
		if _, err := fs.Stat(fsys, path); err == nil {
			return path, nil
		}

		if r.category != "" {
			// .../category/locale.mo
			path = filepath.Join(r.category, locale+ext)
			if _, err := fs.Stat(fsys, path); err == nil {
				return path, nil
			}
		}

		if r.category != "LC_MESSAGES" {
			// .../LC_MESSAGES/locale.mo
			path = filepath.Join("LC_MESSAGES", locale+ext)
			if _, err := fs.Stat(fsys, path); err == nil {
				return path, nil
			}
		}
	}

	return "", os.ErrNotExist
}
