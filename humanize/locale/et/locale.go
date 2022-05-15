package et

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("et"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:      "j. F Y",
			TimeFormat:      "G:i",
			MonthDayFormat:  "j. F",
			ShortDateFormat: "d.m.Y",
		},
	}
}
