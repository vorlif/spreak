package kn

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("kn"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:      "j F Y",
			TimeFormat:      "h:i A",
			MonthDayFormat:  "j F",
			ShortDateFormat: "j M Y",
		},
	}
}
