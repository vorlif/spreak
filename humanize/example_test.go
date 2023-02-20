package humanize_test

import (
	"fmt"
	"math"
	"strings"
	"time"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/humanize"
	"github.com/vorlif/spreak/humanize/locale/ar"
	"github.com/vorlif/spreak/humanize/locale/be"
	"github.com/vorlif/spreak/humanize/locale/de"
	"github.com/vorlif/spreak/humanize/locale/es"
	"github.com/vorlif/spreak/humanize/locale/it"
	"github.com/vorlif/spreak/humanize/locale/zhHans"
)

func ExampleHumanizer_Apnumber() {
	collection := humanize.MustNew(humanize.WithLocale(es.New(), de.New(), zhHans.New()))

	for _, tag := range []language.Tag{language.English, language.German, language.Spanish, language.SimplifiedChinese} {
		h := collection.CreateHumanizer(tag)
		fmt.Println(h.Apnumber(5))
	}
	// Output:
	// five
	// fünf
	// cinco
	// 五
}

func ExampleHumanizer_Intword() {
	collection := humanize.MustNew(humanize.WithLocale(es.New(), de.New(), zhHans.New()))

	for _, tag := range []language.Tag{language.English, language.Spanish, language.SimplifiedChinese} {
		h := collection.CreateHumanizer(tag)
		fmt.Println(h.Intword(1_000_000_000))
	}
	// Output:
	// 1.0 billion
	// 1,0 millardo
	// 1.0 十亿
}

func ExampleHumanizer_Intcomma() {
	collection := humanize.MustNew(humanize.WithLocale(es.New(), de.New(), zhHans.New()))

	for _, tag := range []language.Tag{language.English, language.Spanish, language.SimplifiedChinese} {
		h := collection.CreateHumanizer(tag)
		fmt.Println(h.Intcomma(1_000_000_000))
	}
	// Output:
	// 1,000,000,000
	// 1.000.000.000
	// 1,000,000,000
}

func ExampleHumanizer_Ordinal() {
	collection := humanize.MustNew(humanize.WithLocale(es.New(), de.New(), zhHans.New()))

	for _, tag := range []language.Tag{language.English, language.Spanish, language.SimplifiedChinese} {
		h := collection.CreateHumanizer(tag)
		fmt.Println(h.Ordinal(5))
	}
	// Output:
	// 5th
	// 5º
	// 第5
}

func ExampleHumanizer_LanguageName() {
	collection := humanize.MustNew(humanize.WithLocale(es.New(), de.New(), zhHans.New()))

	for _, tag := range []language.Tag{language.English, language.Spanish, language.SimplifiedChinese} {
		h := collection.CreateHumanizer(tag)
		fmt.Printf("Examples in: %s, %s", h.LanguageName("English"), h.LanguageName("Spanish"))
		fmt.Printf(" and %s\n", h.LanguageName("Simplified Chinese"))
	}
	// Output:
	// Examples in: English, Spanish and Simplified Chinese
	// Examples in: Inglés, Español and Chino simplificado
	// Examples in: 英语, 西班牙语 and 简体中文
}

func ExampleHumanizer_NaturalDay() {
	collection := humanize.MustNew(humanize.WithLocale(es.New(), de.New(), zhHans.New()))

	for _, tag := range []language.Tag{language.English, language.Spanish, language.SimplifiedChinese} {
		h := collection.CreateHumanizer(tag)
		fmt.Println(h.NaturalDay(time.Now()))
	}

	fmt.Println("---")

	d := time.Date(2022, 05, 01, 0, 0, 0, 0, time.UTC)
	for _, tag := range []language.Tag{language.English, language.Spanish, language.SimplifiedChinese} {
		h := collection.CreateHumanizer(tag)
		fmt.Println(h.NaturalDay(d))
	}
	// Output:
	// today
	// hoy
	// 今天
	// ---
	// May 1, 2022
	// 1 de mayo de 2022
	// 2022年5月1日
}

func ExampleHumanizer_NaturalTime() {
	collection := humanize.MustNew(humanize.WithLocale(es.New(), ar.New(), zhHans.New()))

	for _, tag := range []language.Tag{language.English, language.Arabic, language.SimplifiedChinese} {
		h := collection.CreateHumanizer(tag)
		t := time.Now().Add(5 * time.Minute)
		fmt.Println(h.NaturalTime(t))
	}
	// Output:
	// 5 minutes from now
	// ٥ دقائق من الآن
	// 5分钟以后
}

func ExampleHumanizer_NaturalTime_past() {
	collection := humanize.MustNew(humanize.WithLocale(es.New(), ar.New(), zhHans.New()))

	for _, tag := range []language.Tag{language.English, language.Arabic, language.SimplifiedChinese} {
		h := collection.CreateHumanizer(tag)
		t := time.Now().Add(-2*time.Hour - 30*time.Minute)
		fmt.Println(h.NaturalTime(t))
	}
	// Output:
	// 2 hours ago
	// منذ ٢ ساعة
	// 2小时之前
}

func ExampleHumanizer_TimeSince() {
	collection := humanize.MustNew(humanize.WithLocale(es.New(), be.New(), de.New()))

	t := time.Now().Add(-37 * time.Hour)
	for _, tag := range []language.Tag{language.English, language.MustParse("be"), language.German} {
		h := collection.CreateHumanizer(tag)
		fmt.Println(h.TimeSince(t))
	}
	// Output:
	// 1 day, 13 hours
	// 1 дзень, 13 гадзін
	// 1 Tag, 13 Stunden
}

func ExampleHumanizer_TimeSince_duration() {
	collection := humanize.MustNew(humanize.WithLocale(es.New(), be.New(), de.New()))

	t := -80 * time.Hour
	for _, tag := range []language.Tag{language.English, language.MustParse("be"), language.German} {
		h := collection.CreateHumanizer(tag)
		fmt.Println(h.TimeSince(t))
	}
	// Output:
	// 3 days, 8 hours
	// 3  дзён, 8 гадзін
	// 3 Tage, 8 Stunden
}

func ExampleHumanizer_TimeUntil() {
	collection := humanize.MustNew(humanize.WithLocale(es.New(), be.New(), de.New()))

	t := time.Now().Add(37 * time.Hour)
	for _, tag := range []language.Tag{language.English, language.MustParse("be"), language.German} {
		h := collection.CreateHumanizer(tag)
		fmt.Println(h.TimeUntil(t))
	}
	// Output:
	// 1 day, 13 hours
	// 1 дзень, 13 гадзін
	// 1 Tag, 13 Stunden
}

func ExampleHumanizer_FormatTime() {
	collection := humanize.MustNew(humanize.WithLocale(es.New(), be.New(), de.New()))

	now := time.Date(2022, 05, 15, 18, 0, 0, 0, time.Local)
	for _, tag := range []language.Tag{language.English, language.MustParse("be"), language.German} {
		h := collection.CreateHumanizer(tag)
		fmt.Println(h.FormatTime(now, humanize.DateTimeFormat))
	}
	// Output:
	// May 15, 2022, 6 p.m.
	// траўня 15, 2022, 6 папаўдні
	// 15. Mai 2022 18:00
}

func ExampleHumanizer_FormatTime_shortDate() {
	collection := humanize.MustNew(humanize.WithLocale(es.New(), zhHans.New(), de.New()))

	now := time.Date(2022, 05, 15, 18, 0, 0, 0, time.Local)
	for _, tag := range []language.Tag{language.English, language.SimplifiedChinese, language.German} {
		h := collection.CreateHumanizer(tag)
		fmt.Println(h.FormatTime(now, humanize.ShortDateFormat))
	}
	// Output:
	// 05/15/2022
	// 2022年5月15日
	// 15.05.2022
}

func ExampleHumanizer_FormatTime_shortDateTime() {
	collection := humanize.MustNew(humanize.WithLocale(es.New(), zhHans.New(), de.New()))

	now := time.Date(2022, 05, 15, 18, 0, 0, 0, time.Local)
	for _, tag := range []language.Tag{language.English, language.SimplifiedChinese, language.German} {
		h := collection.CreateHumanizer(tag)
		fmt.Println(h.FormatTime(now, humanize.ShortDatetimeFormat))
	}
	// Output:
	// 05/15/2022 6 p.m.
	// 2022年5月15日 18:00
	// 15.05.2022 18:00
}

func ExampleHumanizer_Date() {
	collection := humanize.MustNew(humanize.WithLocale(es.New(), it.New(), de.New()))

	for _, tag := range []language.Tag{language.English, language.Italian, language.German} {
		h := collection.CreateHumanizer(tag)
		fmt.Println(h.Date())
	}
	// May 16, 2022
	// 16 Maggio 2022
	// 16. Mai 2022
}

func ExampleHumanizer_Time() {
	collection := humanize.MustNew(humanize.WithLocale(es.New(), it.New(), de.New()))

	for _, tag := range []language.Tag{language.English, language.Italian, language.German} {
		h := collection.CreateHumanizer(tag)
		fmt.Println(h.Time())
	}
	// 12:30 a.m.
	// 00:30
	// 00:30
}

func ExampleHumanizer_Now() {
	collection := humanize.MustNew(humanize.WithLocale(es.New(), it.New(), de.New()))

	for _, tag := range []language.Tag{language.English, language.Italian, language.German} {
		h := collection.CreateHumanizer(tag)
		fmt.Println(h.Now())
	}
	// May 16, 2022, 12:34 a.m.
	// Lunedì 16 Maggio 2022 00:34
	// 16. Mai 2022 00:34
}

func ExampleHumanizer_FilesizeFormat() {
	collection := humanize.MustNew(humanize.WithLocale(es.New(), zhHans.New(), de.New()))

	for _, tag := range []language.Tag{language.English, language.SimplifiedChinese, language.German} {
		h := collection.CreateHumanizer(tag)
		fmt.Println(h.FilesizeFormat(make([]byte, 1000)))
		fmt.Println(h.FilesizeFormat(math.Pow(1024, 3)))
	}
	// Output:
	// 1,000 bytes
	// 1 GB
	// 1,000 字节
	// 1 GB
	// 1.000 Bytes
	// 1 GB
}

func ExampleWithDepth() {
	collection := humanize.MustNew(humanize.WithLocale(es.New(), be.New(), de.New()))

	t := -3080*time.Hour - 5*time.Minute
	for _, tag := range []language.Tag{language.English, language.MustParse("be"), language.German} {
		h := collection.CreateHumanizer(tag)
		fmt.Println(h.TimeSince(t, humanize.WithDepth(1), humanize.WithoutAdjacentCheck()))
		fmt.Println(h.TimeSince(t, humanize.WithoutAdjacentCheck())) // 2 is default
		fmt.Println(h.TimeSince(t, humanize.WithDepth(3), humanize.WithoutAdjacentCheck()))
		fmt.Println(h.TimeSince(t, humanize.WithDepth(4), humanize.WithoutAdjacentCheck()))
		fmt.Println(h.TimeSince(t, humanize.WithDepth(5), humanize.WithoutAdjacentCheck()))
		fmt.Println(strings.Repeat(" -", 20))
	}
	// Output:
	// 4 months
	// 4 months, 1 week
	// 4 months, 1 week, 1 day
	// 4 months, 1 week, 1 day, 8 hours
	// 4 months, 1 week, 1 day, 8 hours, 5 minutes
	//  - - - - - - - - - - - - - - - - - - - -
	// 4 месяцаў
	// 4 месяцаў, 1 тыдзень
	// 4 месяцаў, 1 тыдзень, 1 дзень
	// 4 месяцаў, 1 тыдзень, 1 дзень, 8 гадзін
	// 4 месяцаў, 1 тыдзень, 1 дзень, 8 гадзін, 5 хвілін
	//  - - - - - - - - - - - - - - - - - - - -
	// 4 Monate
	// 4 Monate, 1 Woche
	// 4 Monate, 1 Woche, 1 Tag
	// 4 Monate, 1 Woche, 1 Tag, 8 Stunden
	// 4 Monate, 1 Woche, 1 Tag, 8 Stunden, 5 Minuten
	//  - - - - - - - - - - - - - - - - - - - -
}

func ExampleWithNow() {
	collection := humanize.MustNew(humanize.WithLocale(es.New(), zhHans.New(), de.New()))

	now := time.Date(1999, 12, 24, 00, 0, 0, 0, time.Local)
	event := time.Date(2000, 01, 01, 00, 0, 0, 0, time.Local)
	for _, tag := range []language.Tag{language.English, language.SimplifiedChinese, language.German} {
		h := collection.CreateHumanizer(tag)
		fmt.Println(h.TimeUntil(event, humanize.WithNow(now)))
	}
	// Output:
	// 1 week, 1 day
	// 1 周，1 日
	// 1 Woche, 1 Tag
}
