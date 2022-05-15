package sk

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("sk"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          "j. F Y",
			TimeFormat:          "G:i",
			DateTimeFormat:      "j. F Y G:i",
			YearMonthFormat:     "F Y",
			MonthDayFormat:      "j. F",
			ShortDateFormat:     "d.m.Y",
			ShortDatetimeFormat: "d.m.Y G:i",
			FirstDayOfWeek:      1,
		},
	}
}
