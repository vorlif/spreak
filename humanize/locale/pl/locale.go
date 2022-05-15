package pl

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("pl"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          "j E Y",
			TimeFormat:          "H:i",
			DateTimeFormat:      "j E Y H:i",
			YearMonthFormat:     "F Y",
			MonthDayFormat:      "j E",
			ShortDateFormat:     "d-m-Y",
			ShortDatetimeFormat: "d-m-Y  H:i",
			FirstDayOfWeek:      1,
		},
	}
}
