package ka

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("ka"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          "l, j F, Y",
			TimeFormat:          "h:i a",
			DateTimeFormat:      "j F, Y h:i a",
			YearMonthFormat:     "F, Y",
			MonthDayFormat:      "j F",
			ShortDateFormat:     "j.M.Y",
			ShortDatetimeFormat: "j.M.Y H:i",
			FirstDayOfWeek:      1,
		},
	}
}
