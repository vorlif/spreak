package vi

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("vi"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          `\N\gà\y d \t\há\n\g n \nă\m Y`,
			TimeFormat:          "H:i",
			DateTimeFormat:      `:i \N\gà\y d \t\há\n\g n \nă\m Y`,
			YearMonthFormat:     "F Y",
			MonthDayFormat:      "j F",
			ShortDateFormat:     "d-m-Y",
			ShortDatetimeFormat: "H:i d-m-Y",
		},
	}
}
