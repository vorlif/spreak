package ky

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("ky"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          "j E Y ж.",
			TimeFormat:          "G:i",
			DateTimeFormat:      "j E Y ж. G:i",
			YearMonthFormat:     "F Y ж.",
			MonthDayFormat:      "j F",
			ShortDateFormat:     "d.m.Y",
			ShortDatetimeFormat: "d.m.Y H:i",
			FirstDayOfWeek:      1,
		},
	}
}
