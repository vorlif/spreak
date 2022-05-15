package mn

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("mn"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:      "d F Y",
			TimeFormat:      "g:i A",
			ShortDateFormat: "j M Y",
		},
	}
}
