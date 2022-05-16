package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak"
)

// A resolver resolves the file path from which the file is to be loaded.
type myResolver struct{}

var _ spreak.Resolver = (*myResolver)(nil)

func (myResolver) Resolve(fsys fs.FS, extensions string, lang language.Tag, domain string) (string, error) {
	// We create the path where the file should be located
	path := filepath.Join("locale", lang.String(), domain+extensions)

	// Verify if the file exists
	if _, err := fs.Stat(fsys, path); err == nil {
		// And pass the loader the path
		return path, nil
	}

	// If the file does not exist, but we want to continue without the file, we return os.ErrNotExist.
	return "", os.ErrNotExist
}

func main() {
	// We would like to define ourselves how to resolve the paths and filenames to the files.
	// Therefore, we need to create our own FilesystemLoader
	fsLoader, errFS := spreak.NewFilesystemLoader(
		spreak.WithResolver(&myResolver{}),
		spreak.WithPath("../../"),
	)
	if errFS != nil {
		panic(errFS)
	}

	bundle, err := spreak.NewBundle(
		spreak.WithSourceLanguage(language.English),
		spreak.WithDefaultDomain("helloworld"),
		// We now load the domains via our own FilesystemLoader, which uses our resolver.
		spreak.WithDomainLoader("helloworld", fsLoader),
		spreak.WithRequiredLanguage(language.Spanish),
		spreak.WithLanguage(language.German, language.French),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(bundle.SupportedLanguages())

	t := spreak.NewLocalizer(bundle, language.Spanish)
	fmt.Println(t.Get("Hello world"))
	// Output: Hola Mundo
}
