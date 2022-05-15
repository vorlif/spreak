package gd

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("gd"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          "j F Y",
			TimeFormat:          "h:ia",
			DateTimeFormat:      "j F Y h:ia",
			MonthDayFormat:      "j F",
			ShortDateFormat:     "j M Y",
			ShortDatetimeFormat: "j M Y h:ia",
			FirstDayOfWeek:      1,
		},
	}
}
