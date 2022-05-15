package eu

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("eu"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          `Y\k\o N j\a`,
			TimeFormat:          "H:i",
			DateTimeFormat:      `Y\k\o N j\a, H:i`,
			YearMonthFormat:     `Y\k\o F`,
			MonthDayFormat:      `F\r\e\n j\a`,
			ShortDateFormat:     "Y-m-d",
			ShortDatetimeFormat: "Y-m-d H:i",
			FirstDayOfWeek:      1,
		},
	}
}
