package uz

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("uz"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          `j-E, Y-\y\i\l`,
			TimeFormat:          "G:i",
			DateTimeFormat:      `j-E, Y-\y\i\l G:i`,
			YearMonthFormat:     `F Y-\y\i\l`,
			MonthDayFormat:      "j-E",
			ShortDateFormat:     "d.m.Y",
			ShortDatetimeFormat: "d.m.Y H:i",
			FirstDayOfWeek:      1,
		},
	}
}
