package el

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("el"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          "d/m/Y",
			TimeFormat:          "P",
			DateTimeFormat:      "d/m/Y P",
			YearMonthFormat:     "F Y",
			MonthDayFormat:      "j F",
			ShortDateFormat:     "d/m/Y",
			ShortDatetimeFormat: "d/m/Y P",
			FirstDayOfWeek:      0,
		},
	}
}
