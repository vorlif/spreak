package ca

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("ca"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          `j E \d\e Y`,
			TimeFormat:          "G:i",
			DateTimeFormat:      `j E \d\e Y \a \l\e\s G:i`,
			YearMonthFormat:     `F \d\e\l Y`,
			MonthDayFormat:      "j E",
			ShortDateFormat:     "d/m/Y",
			ShortDatetimeFormat: "d/m/Y G:i",
			FirstDayOfWeek:      1,
		},
	}
}
