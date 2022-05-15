package ko

import (
	"embed"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
)

//go:embed *.po
var fsys embed.FS

func New() *humanize.LocaleData {
	return &humanize.LocaleData{
		Lang: language.MustParse("ko"),
		Fs:   fsys,
		Format: &humanize.FormatData{
			DateFormat:          "Y년 n월 j일",
			TimeFormat:          "A g:i",
			DateTimeFormat:      "Y년 n월 j일 g:i A",
			YearMonthFormat:     "Y년 n월",
			MonthDayFormat:      "n월 j일",
			ShortDateFormat:     "Y-n-j.",
			ShortDatetimeFormat: "Y-n-j H:i",
		},
	}
}
