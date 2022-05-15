package es

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.Spanish,
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          `j \d\e F \d\e Y`,
			TimeFormat:          "H:i",
			DateTimeFormat:      `j \d\e F \d\e Y \a \l\a\s H:i`,
			YearMonthFormat:     `F \d\e Y`,
			MonthDayFormat:      `j \d\e F`,
			ShortDateFormat:     "d/m/Y",
			ShortDatetimeFormat: "d/m/Y H:i",
			FirstDayOfWeek:      1,
		},
	}
}
