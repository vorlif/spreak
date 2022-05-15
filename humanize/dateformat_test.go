package humanize

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func formatTime(h *Humanizer, t time.Time, format string) string {
	tf := newTimeFormatter(h, t)
	return tf.format(format)
}

func init() {
	_ = os.Setenv("TZ", "CET")
}

func TestTimeFormatter(t *testing.T) {
	h := createSourceHumanizer(t)

	now := time.Now()
	myBirthday := time.Date(1979, 7, 8, 22, 00, 0, 0, time.Local)
	summertime := time.Date(2005, 10, 30, 1, 00, 0, 0, time.Local)
	wintertime := time.Date(2005, 10, 30, 4, 00, 0, 0, time.Local)
	noon := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())

	tz := time.FixedZone("myzone", -210*60)
	awareTz := time.Date(2009, 5, 16, 5, 30, 30, 0, tz)

	t.Run("test_date", func(t *testing.T) {
		d := time.Date(2009, 5, 16, 5, 30, 30, 0, time.Local)
		unixEpochStr := formatTime(h, d, "U")
		e, err := strconv.ParseInt(unixEpochStr, 10, 64)
		require.NoError(t, err)
		assert.Equal(t, time.Unix(e, 0), d)
	})

	t.Run("test_am_pm", func(t *testing.T) {
		d := time.Date(2009, 5, 16, 7, 0, 0, 0, time.Local)
		assert.Equal(t, formatTime(h, d, "a"), "a.m.")
		assert.Equal(t, formatTime(h, d, "A"), "AM")

		d = time.Date(2009, 5, 16, 19, 0, 0, 0, time.Local)
		assert.Equal(t, formatTime(h, d, "a"), "p.m.")
		assert.Equal(t, formatTime(h, d, "A"), "PM")
	})

	t.Run("test epoch", func(t *testing.T) {
		assert.Equal(t, formatTime(h, time.UnixMilli(0), "U"), "0")
	})

	t.Run("test microseconds", func(t *testing.T) {
		d := time.UnixMicro(123)
		assert.Equal(t, "000123", formatTime(h, d, "u"))
	})

	t.Run("test date formats", func(t *testing.T) {
		tests := []struct {
			specifier string
			expected  string
		}{
			{"b", "jul"},
			{"d", "08"},
			{"D", "Sun"},
			{"E", "July"},
			{"F", "July"},
			{"j", "8"},
			{"l", "Sunday"},
			{"L", "false"},
			{"m", "07"},
			{"M", "Jul"},
			{"n", "7"},
			{"N", "July"},
			{"o", "1979"},
			{"S", "th"},
			{"t", "31"},
			{"w", "0"},
			{"W", "27"},
			{"y", "79"},
			{"Y", "1979"},
			{"z", "189"},
		}

		for _, tt := range tests {
			t.Run(tt.specifier, func(t *testing.T) {
				assert.Equal(t, tt.expected, formatTime(h, myBirthday, tt.specifier))
			})
		}
	})

	t.Run("test_empty_format", func(t *testing.T) {
		assert.Equal(t, formatTime(h, myBirthday, ""), "")
	})

	t.Run("test date formats", func(t *testing.T) {
		tests := []struct {
			specifier string
			expected  string
		}{
			{"a", "p.m."},
			{"A", "PM"},
			{"f", "10"},
			{"g", "10"},
			{"G", "22"},
			{"h", "10"},
			{"H", "22"},
			{"i", "00"},
			{"P", "10 p.m."},
			{"s", "00"},
			{"u", "000000"},
		}
		for _, tt := range tests {
			t.Run(tt.specifier, func(t *testing.T) {
				assert.Equal(t, tt.expected, formatTime(h, myBirthday, tt.specifier))
			})
		}
	})

	t.Run("test future date", func(t *testing.T) {
		theFuture := time.Date(2100, 10, 25, 0, 0, 0, 0, time.Local)
		assert.Equal(t, "2100", formatTime(h, theFuture, "Y"))
	})

	t.Run("test day of year leap", func(t *testing.T) {
		d := time.Date(2000, 12, 31, 0, 0, 0, 0, time.Local)
		assert.Equal(t, "366", formatTime(h, d, "z"))
	})

	t.Run("test timezones", func(t *testing.T) {
		tests := []struct {
			specifier string
			expected  string
		}{
			{"e", "CEST"},
			{"O", "+0200"},
			{"r", "Sun, 08 Jul 1979 22:00:00 +0200"},
			{"T", "CET"},
			{"U", "300312000"},
			{"Z", "7200"},
		}
		for _, tt := range tests {
			t.Run(tt.specifier, func(t *testing.T) {
				assert.Equal(t, tt.expected, formatTime(h, myBirthday, tt.specifier))
			})
		}

		assert.Equal(t, "myzone", formatTime(h, awareTz, "e"))
		assert.Equal(t, "Sat, 16 May 2009 05:30:30 -0330", formatTime(h, awareTz, "r"))

		assert.Equal(t, "1", formatTime(h, summertime, "I"))
		assert.Equal(t, "+0200", formatTime(h, summertime, "O"))

		assert.Equal(t, "0", formatTime(h, wintertime, "I"))
		assert.Equal(t, "+0100", formatTime(h, wintertime, "O"))

		for _, specifier := range []string{"e", "O", "T", "Z"} {
			t.Run(specifier, func(t *testing.T) {
				assert.NotEmpty(t, formatTime(h, noon, specifier))
			})
		}

		assert.Equal(t, "-0330", formatTime(h, awareTz, "O"))
	})

	t.Run("test_e_format_with_named_time_zone", func(t *testing.T) {
		d := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		assert.Equal(t, "UTC", formatTime(h, d, "e"))
	})

	t.Run("test_e_format_with_time_zone_with_unimplemented_tzname", func(t *testing.T) {
		tz := time.FixedZone("", 0)
		d := time.Date(1970, 1, 1, 0, 0, 0, 0, tz)
		assert.Equal(t, "", formatTime(h, d, "e"))
	})

	t.Run("test P format", func(t *testing.T) {
		tests := []struct {
			expected string
			hour     int
			minute   int
		}{
			{"midnight", 0, 0},
			{"noon", 12, 0},
			{"4 a.m.", 4, 0},
			{"8:30 a.m.", 8, 30},
			{"4 p.m.", 16, 0},
			{"8:30 p.m.", 20, 30},
		}

		for _, tt := range tests {
			t.Run(tt.expected, func(t *testing.T) {
				d := time.Date(1970, 1, 1, tt.hour, tt.minute, 0, 0, tz)
				assert.Equal(t, tt.expected, formatTime(h, d, "P"))
			})
		}
	})

	t.Run("test S format", func(t *testing.T) {
		var thDays []int
		for i := 4; i < 31; i++ {
			if i != 21 && i != 22 && i != 23 {
				thDays = append(thDays, i)
			}
		}

		tests := []struct {
			expected string
			days     []int
		}{
			{"st", []int{1, 21, 31}},
			{"nd", []int{2, 22}},
			{"rd", []int{3, 23}},
			{"th", thDays},
		}
		for _, tt := range tests {
			t.Run(tt.expected, func(t *testing.T) {
				for _, day := range tt.days {
					d := time.Date(1970, 1, day, 0, 0, 0, 0, tz)
					assert.Equal(t, tt.expected, formatTime(h, d, "S"))
				}
			})
		}
	})

	t.Run("test y format year before 1000", func(t *testing.T) {
		tests := []struct {
			expected string
			year     int
		}{
			{"76", 476},
			{"42", 42},
			{"04", 4},
		}

		for _, tt := range tests {
			t.Run(tt.expected, func(t *testing.T) {
				d := time.Date(tt.year, 9, 8, 5, 0, 0, 0, time.Local)
				assert.Equal(t, tt.expected, formatTime(h, d, "y"))
			})
		}
	})

	t.Run("test twelve hour format", func(t *testing.T) {
		tests := []struct {
			gExpected string
			hExpected string
			hour      int
		}{
			{"12", "12", 0},
			{"1", "01", 1},
			{"11", "11", 11},
			{"12", "12", 12},
			{"1", "01", 13},
			{"11", "11", 23},
		}

		for _, tt := range tests {
			t.Run(tt.gExpected, func(t *testing.T) {
				d := time.Date(2000, 1, 1, tt.hour, 0, 0, 0, tz)
				assert.Equal(t, tt.gExpected, formatTime(h, d, "g"))
				assert.Equal(t, tt.hExpected, formatTime(h, d, "h"))
			})
		}
	})

	t.Run("test escaping", func(t *testing.T) {
		d := time.Date(2001, 03, 10, 17, 16, 18, 0, time.UTC)
		loc := time.Local
		time.Local = time.UTC
		defer func() {
			time.Local = loc
		}()
		tests := []struct {
			format   string
			expected string
		}{
			{"F j, Y, g:i a", "March 10, 2001, 5:16 p.m."},
			{"m.d.y", "03.10.01"},
			{"j, n, Y", "10, 3, 2001"},
			{"Ymd", "20010310"},
			{`h-i-s, j-m-y, it is w Day`, "05-16-18, 10-03-01, 1631 1618 6 Satp.m.01"},
			{`\i\t \i\s \t\h\e jS \d\a\y.`, "it is the 10th day."},
			{"D M j G:i:s T Y", "Sat Mar 10 17:16:18 UTC 2001"},
			{`H:m:s \m \i\s\ \m\o\n\t\h`, "17:03:18 m is month"},
			{"H:i:s", "17:16:18"},
			{"Y-m-d H:i:s", "2001-03-10 17:16:18"},
		}

		for _, tt := range tests {
			t.Run(tt.format+" => "+tt.expected, func(t *testing.T) {
				assert.Equal(t, tt.expected, formatTime(h, d, tt.format))
			})
		}
	})

	t.Run("test iso", func(t *testing.T) {
		assert.Equal(t, "43", formatTime(h, summertime, "W"))
	})
}

func TestTimeFormatterWithI18n(t *testing.T) {
	h := createGermanHumanizer(t)

	t.Run("test translations", func(t *testing.T) {
		d := time.Date(2001, 03, 10, 17, 16, 18, 0, time.UTC)

		tests := []struct {
			format   string
			expected string
			time     time.Time
		}{
			{"b", "Jan", d.AddDate(0, -2, 0)},
			{"b", "Feb", d.AddDate(0, -1, 0)},
			{"b", "Mär", d},
			{"E", "März", d},
			{"E", "Januar", d.AddDate(0, -2, 0)},
			{"l", "Samstag", d},
			{"l", "Sonntag", d.AddDate(0, 0, 1)},
			{"l", "Montag", d.AddDate(0, 0, 2)},
			{"L", "false", d},
			{"L", "true", d.AddDate(-1, 0, 0)},
			{"N", "Jan.", d.AddDate(0, -2, 0)},
			{"N", "Feb.", d.AddDate(0, -1, 0)},
			{"N", "März", d},
		}

		for _, tt := range tests {
			t.Run(tt.format+" => "+tt.expected, func(t *testing.T) {
				assert.Equal(t, tt.expected, formatTime(h, tt.time, tt.format))
			})
		}
	})

	t.Run("test date formats c format", func(t *testing.T) {
		d := time.Date(2008, 5, 19, 11, 45, 23, 123456, time.Local)
		assert.Equal(t, "2008-05-19T11:45:23.123456", formatTime(h, d, "c"))
	})

	t.Run("test escaping", func(t *testing.T) {
		d := time.Date(2001, 03, 10, 17, 16, 18, 0, time.UTC)
		loc := time.Local
		time.Local = time.UTC
		defer func() {
			time.Local = loc
		}()

		tests := []struct {
			format   string
			expected string
		}{
			{"F j, Y, g:i a", "März 10, 2001, 5:16 nachm."},
			{"m.d.y", "03.10.01"},
			{"j, n, Y", "10, 3, 2001"},
			{"Ymd", "20010310"},
			{`h-i-s, j-m-y, it is w Day\`, "05-16-18, 10-03-01, 1631 1618 6 Sanachm.01\\"},
			{`\i\t \i\s \t\h\e jS \d\a\y.`, "it is the 10th day."},
			{"D M j G:i:s T Y", "Sa Mär 10 17:16:18 UTC 2001"},
			{`H:m:s \m \i\s\ \m\o\n\t\h`, "17:03:18 m is month"},
			{"H:i:s", "17:16:18"},
			{"Y-m-d H:i:s", "2001-03-10 17:16:18"},
		}

		for _, tt := range tests {
			t.Run(tt.format+" => "+tt.expected, func(t *testing.T) {
				assert.Equal(t, tt.expected, formatTime(h, d, tt.format))
			})
		}
	})
}
