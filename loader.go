package spreak

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/internal/cldrplural"
	"github.com/vorlif/spreak/internal/mo"
	"github.com/vorlif/spreak/internal/poplural"
	"github.com/vorlif/spreak/internal/util"
	"github.com/vorlif/spreak/pkg/po"
)

const (
	UnknownFile = "unknown"
	PoFile      = ".po"
	MoFile      = ".mo"

	poCLDRHeader = "X-spreak-use-CLDR"
)

// Loader is responsible for loading Catalogs for a language and a domain.
// A bundle loads each domain through its own loader.
//
// If a loader cannot find a matching catalog for it must return error spreak.ErrNotFound,
// otherwise the bundle creation will be aborted with the returned error.
type Loader interface {
	Load(lang language.Tag, domain string) (Catalog, error)
}

// A Decoder reads and decodes catalogs for a language and a domain from a byte array.
type Decoder interface {
	Decode(lang language.Tag, domain string, data []byte) (Catalog, error)
}

// A Resolver is used by the FilesystemLoader to resolve the appropriate path for a file.
// If a file was not found, os.ErrNotFound should be returned.
// All other errors cause the loaders search to stop.
//
// fsys represents the file system from which the FilesystemLoader wants to load the file.
// extensions is the file extension for which the file is to be resolved.
// Language and Domain indicate for which domain in which language the file is searched.
type Resolver interface {
	Resolve(fsys fs.FS, extensions string, lang language.Tag, domain string) (string, error)
}

// FsOption is an option which can be used when creating the FilesystemLoader.
type FsOption func(l *FilesystemLoader) error

// ResolverOption is an option which can be used when creating the DefaultResolver.
type ResolverOption func(r *defaultResolver)

// FilesystemLoader is a Loader which loads and decodes files from a file system.
// A file system here means an implementation of fs.FS.
type FilesystemLoader struct {
	fsys       fs.FS
	resolver   Resolver
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

	if l.resolver == nil {
		resolver, err := NewDefaultResolver()
		if err != nil {
			return nil, err
		}
		l.resolver = resolver
	}

	return l, nil
}

func (l *FilesystemLoader) Load(lang language.Tag, domain string) (Catalog, error) {

	for i, extension := range l.extensions {
		path, errP := l.resolver.Resolve(l.fsys, extension, lang, domain)
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

// WithResolver stores the resolver of a FilesystemLoader.
// Lets the creation of the FilesystemLoader fail, if a Resolver was already deposited.
func WithResolver(resolver Resolver) FsOption {
	return func(l *FilesystemLoader) error {
		if l.resolver != nil {
			return errors.New("spreak.Loader: Resolver for FilesystemLoader already set")
		}
		l.resolver = resolver
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

type defaultResolver struct {
	search   bool
	category string
}

// NewDefaultResolver create a resolver which can be used for a FilesystemLoader.
func NewDefaultResolver(opts ...ResolverOption) (Resolver, error) {
	l := &defaultResolver{
		search:   true,
		category: "",
	}

	for _, opt := range opts {
		opt(l)
	}

	return l, nil
}

func WithDisabledSearch() ResolverOption { return func(r *defaultResolver) { r.search = false } }

func WithCategory(category string) ResolverOption {
	return func(l *defaultResolver) {
		l.category = category
	}
}

func (r *defaultResolver) Resolve(fsys fs.FS, extension string, tag language.Tag, domain string) (string, error) {
	for _, lang := range ExpandLanguage(tag) {
		path, err := r.searchFileForLanguageName(fsys, lang, domain, extension)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}

		return path, nil
	}

	return "", os.ErrNotExist
}

func (r *defaultResolver) searchFileForLanguageName(fsys fs.FS, locale, domain, ext string) (string, error) {

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

type poDecoder struct {
	useCLDRPlural bool
}

type moDecoder struct {
	useCLDRPlural bool
}

// NewPoDecoder returns a new Decoder for reading po files.
// If a plural forms header is set, it will be used.
// Otherwise, the CLDR plural rules are used to set the plural form.
// If there is no CLDR plural rule, the English plural rules will be used.
func NewPoDecoder() Decoder { return &poDecoder{} }

// NewPoCLDRDecoder creates a decoder for reading po files,
// which always uses the CLDR plural rules for determining the plural form.
func NewPoCLDRDecoder() Decoder { return &poDecoder{useCLDRPlural: true} }

// NewMoDecoder returns a new Decoder for reading mo files.
// If a plural forms header is set, it will be used.
// Otherwise, the CLDR plural rules are used to set the plural form.
// If there is no CLDR plural rule, the English plural rules will be used.
func NewMoDecoder() Decoder { return &moDecoder{useCLDRPlural: false} }

// NewMoCLDRDecoder creates a decoder for reading mo files,
// which always uses the CLDR plural rules for determining the plural form.
func NewMoCLDRDecoder() Decoder { return &moDecoder{useCLDRPlural: true} }

func (d poDecoder) Decode(lang language.Tag, domain string, data []byte) (Catalog, error) {
	poFile, errParse := po.Parse(data)
	if errParse != nil {
		return nil, errParse
	}

	// We could check here if the language of the file matches the target language,
	// but leave it off to make loading more flexible.

	return buildGettextCatalog(poFile, lang, domain, d.useCLDRPlural)
}

func (d moDecoder) Decode(lang language.Tag, domain string, data []byte) (Catalog, error) {
	moFile, errParse := mo.ParseBytes(data)
	if errParse != nil {
		return nil, errParse
	}

	// We could check here if the language of the file matches the target language,
	// but leave it off to make loading more flexible.

	return buildGettextCatalog(moFile, lang, domain, d.useCLDRPlural)
}

func buildGettextCatalog(file *po.File, lang language.Tag, domain string, useCLDRPlural bool) (Catalog, error) {
	messages := make(messageLookupMap, len(file.Messages))

	for ctx := range file.Messages {
		if len(file.Messages[ctx]) == 0 {
			continue
		}

		if _, hasContext := messages[ctx]; !hasContext {
			messages[ctx] = make(map[string]*gettextMessage)
		}

		for msgID, poMsg := range file.Messages[ctx] {
			if msgID == "" {
				continue
			}

			d := &gettextMessage{
				Context:      poMsg.Context,
				ID:           poMsg.ID,
				IDPlural:     poMsg.IDPlural,
				Translations: poMsg.Str,
			}

			messages[poMsg.Context][poMsg.ID] = d
		}
	}

	catl := &gettextCatalog{
		language:     lang,
		translations: messages,
	}

	if useCLDRPlural {
		catl.pluralFunc = getCLDRPluralFunction(lang)
		return catl, nil
	}

	if file.Header != nil {
		if val := file.Header.Get(poCLDRHeader); strings.ToLower(val) == "true" {
			catl.pluralFunc = getCLDRPluralFunction(lang)
			return catl, nil
		}

		if file.Header.PluralForms != "" {
			forms, err := poplural.Parse(file.Header.PluralForms)
			if err != nil {
				return nil, fmt.Errorf("spreak.Decoder: plural forms for po file %v#%v could not be parsed: %w", lang, domain, err)
			}
			catl.pluralFunc = forms.Evaluate
			return catl, nil
		}
	}

	catl.pluralFunc = getCLDRPluralFunction(lang)
	return catl, nil
}

func getCLDRPluralFunction(lang language.Tag) func(a interface{}) int {
	ruleSet, _ := cldrplural.ForLanguage(lang)
	return func(a interface{}) int {
		cat := ruleSet.Evaluate(a)
		for i := 0; i < len(ruleSet.Categories); i++ {
			if ruleSet.Categories[i] == cat {
				return i
			}
		}
		return 0
	}
}
