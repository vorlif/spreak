package ug

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("ug"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          "j F, Y",
			TimeFormat:          "G:i",
			DateTimeFormat:      "",
			YearMonthFormat:     "F Y",
			MonthDayFormat:      "j F",
			ShortDateFormat:     "Y/m/d",
			ShortDatetimeFormat: "Y/m/d G:i",
			FirstDayOfWeek:      1,
		},
	}
}
