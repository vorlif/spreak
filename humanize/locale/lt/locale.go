package lv

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("lv"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          `Y. \g\a\d\a j. F`,
			TimeFormat:          "H:i",
			DateTimeFormat:      `Y. \g\a\d\a j. F, H:i`,
			YearMonthFormat:     `Y. \g. F`,
			MonthDayFormat:      "j. F",
			ShortDateFormat:     `j.m.Y`,
			ShortDatetimeFormat: "j.m.Y H:i",
			FirstDayOfWeek:      1,
		},
	}
}
