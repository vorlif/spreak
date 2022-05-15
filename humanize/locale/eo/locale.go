package eo

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("eo"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          `j\-\a \d\e F Y`,
			TimeFormat:          "H:i",
			DateTimeFormat:      `j\-\a \d\e F Y\, \j\e H:i`,
			YearMonthFormat:     `F \d\e Y`,
			MonthDayFormat:      `j\-\a \d\e F`,
			ShortDateFormat:     `Y-m-d`,
			ShortDatetimeFormat: `Y-m-d H:i`,
			FirstDayOfWeek:      1,
		},
	}
}
