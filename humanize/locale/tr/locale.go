package tr

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("tr"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          "d F Y",
			TimeFormat:          "H:i",
			DateTimeFormat:      "d F Y H:i",
			YearMonthFormat:     "F Y",
			MonthDayFormat:      "d F",
			ShortDateFormat:     "d M Y",
			ShortDatetimeFormat: "d M Y H:i",
			FirstDayOfWeek:      1,
		},
	}
}
