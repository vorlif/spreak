package humanize

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNaturalDay(t *testing.T) {
	h := createGermanHumanizer(t)

	today := time.Now()
	yesterday := time.Now().AddDate(0, 0, -1)
	tomorrow := time.Now().AddDate(0, 0, 1)
	someday := time.Date(2022, 05, 01, 0, 0, 0, 0, time.UTC)

	assert.Equal(t, "heute", h.NaturalDay(today))
	assert.Equal(t, "gestern", h.NaturalDay(yesterday))
	assert.Equal(t, "morgen", h.NaturalDay(tomorrow))
	assert.Equal(t, "1. Mai 2022", h.NaturalDay(someday))

	someday = time.Date(today.Year()+4, today.Month(), today.Day(), 0, 0, 0, today.Nanosecond(), today.Location())
	assert.Equal(t, "16. Mai 2026", h.NaturalDay(someday))

	month := (today.Month() + 1) % 12
	someday = time.Date(today.Year(), month, today.Day(), 0, 0, 0, today.Nanosecond(), today.Location())
	assert.Equal(t, "16. Juni 2022", h.NaturalDay(someday))

	yesterday = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, today.Nanosecond(), today.Location()).Add(-1 * time.Minute)
	assert.Equal(t, "gestern", h.NaturalDay(yesterday))

	h = createSourceHumanizer(t)
	assert.Equal(t, "today", h.NaturalDay(time.Now().Unix()))

	t.Run("invalid data", func(t *testing.T) {
		assert.Equal(t, "%!(string=test)", h.NaturalDay("test"))
	})
}

func TestHumanizer_NaturalTime(t *testing.T) {
	h := createSourceHumanizer(t)

	t.Run("test common use cases", func(t *testing.T) {
		now := time.Now()
		tests := []struct {
			expected string
			time     interface{}
		}{
			{"%!(string=test)", "test"},
			{"a second ago", now.Add(-time.Second)},
			{"30 seconds ago", now.Add(-30 * time.Second)},
			{"a minute ago", now.Add(-90 * time.Second)},
			{"2 minutes ago", now.Add(-2 * time.Minute)},
			{"an hour ago", now.Add(-time.Hour)},
			{"23 hours ago", now.Add(-(23*time.Hour + 50*time.Minute + 50*time.Second))},
			{"1 day ago", now.AddDate(0, 0, -1)},
			{"1 year, 4 months ago", now.AddDate(0, 0, -500)},
			{"a second from now", now.Add(time.Second)},
			{"30 seconds from now", now.Add(30 * time.Second)},
			{"a minute from now", now.Add(90 * time.Second)},
			{"2 minutes from now", now.Add(2 * time.Minute)},
			{"an hour from now", now.Add(time.Hour)},
			{"23 hours from now", now.Add(23*time.Hour + 50*time.Minute + 50*time.Second)},
			{"1 day from now", now.AddDate(0, 0, 1)},
			{"2 days, 6 hours from now", now.Add(2*24*time.Hour + 6*time.Hour)},
			{"1 year, 4 months from now", now.AddDate(0, 0, 500)},
		}

		for _, tt := range tests {
			t.Run(fmt.Sprintf("%v => %s", tt.time, tt.expected), func(t *testing.T) {
				assert.Equal(t, tt.expected, h.NaturalTime(tt.time))
			})
		}
	})

	t.Run("test now", func(t *testing.T) {
		tests := []struct {
			time interface{}
		}{
			{time.Now()},
			{time.Now().Add(-time.Microsecond)},
			{time.Now().Add(-500 * time.Microsecond)},
			{time.Now().Add(time.Microsecond)},
			{time.Now().Add(500 * time.Microsecond)},
			{time.Now().In(time.UTC)},
		}

		ciResults := []string{"now", "a second from now", "a second ago"}
		for _, tt := range tests {
			if os.Getenv("CI") == "" {
				assert.Equal(t, "now", h.NaturalTime(tt.time))
			} else {
				assert.Contains(t, ciResults, h.NaturalTime(tt.time))
			}
		}

	})
}

func TestHumanizer_TimeSince(t *testing.T) {
	h := createSourceHumanizer(t)

	d := time.Date(2007, 8, 14, 13, 46, 0, 0, time.Local)
	day := time.Hour * 24
	week := day * 7
	month := day * 30
	year := day * 365

	t.Run("test invalid input", func(t *testing.T) {
		assert.Equal(t, "%!(string=test)", h.TimeSince("test"))
		assert.Equal(t, "<nil>", h.TimeSince(nil))
		assert.Equal(t, "%!(string=test)", h.TimeUntil("test"))
		assert.Equal(t, "<nil>", h.TimeUntil(nil))
	})

	t.Run("test equal datetimes", func(t *testing.T) {
		assert.Equal(t, "0 minutes", h.TimeSince(d, WithNow(d)))
	})

	t.Run("test ignore microseconds and seconds", func(t *testing.T) {
		assert.Equal(t, "0 minutes", h.TimeSince(d, WithNow(d.Add(time.Microsecond))))
		assert.Equal(t, "0 minutes", h.TimeSince(d, WithNow(d.Add(time.Second))))
	})

	t.Run("other", func(t *testing.T) {
		assert.Equal(t, "1 minute", h.TimeSince(d, WithNow(d.Add(time.Minute))))
		assert.Equal(t, "1 hour", h.TimeSince(d, WithNow(d.Add(time.Hour))))
		assert.Equal(t, "1 day", h.TimeSince(d, WithNow(d.Add(day))))
		assert.Equal(t, "1 week", h.TimeSince(d, WithNow(d.Add(week))))
		assert.Equal(t, "1 month", h.TimeSince(d, WithNow(d.Add(month))))
		assert.Equal(t, "1 year", h.TimeSince(d, WithNow(d.Add(year))))
	})

	t.Run("multiple", func(t *testing.T) {
		now := d.Add(2*day + 6*time.Hour)
		assert.Equal(t, "2 days, 6 hours", h.TimeSince(d, WithNow(now)))
		now = d.Add(2*week + 2*day)
		assert.Equal(t, "2 weeks, 2 days", h.TimeSince(d, WithNow(now)))
	})

	t.Run("display first unit", func(t *testing.T) {
		now := d.Add(2*week + 3*time.Hour + 4*time.Minute)
		assert.Equal(t, "2 weeks", h.TimeSince(d, WithNow(now)))
		now = d.Add(4*day + 5*time.Minute)
		assert.Equal(t, "4 days", h.TimeSince(d, WithNow(now)))
	})

	t.Run("test display second before first", func(t *testing.T) {
		assert.Equal(t, "0 minutes", h.TimeSince(d, WithNow(d.Add(-time.Microsecond))))
		assert.Equal(t, "0 minutes", h.TimeSince(d, WithNow(d.Add(-time.Minute))))
		assert.Equal(t, "0 minutes", h.TimeSince(d, WithNow(d.Add(-time.Hour))))
		assert.Equal(t, "0 minutes", h.TimeSince(d, WithNow(d.Add(-day))))
		assert.Equal(t, "0 minutes", h.TimeSince(d, WithNow(d.Add(-week))))
		assert.Equal(t, "0 minutes", h.TimeSince(d, WithNow(d.Add(-month))))
		assert.Equal(t, "0 minutes", h.TimeSince(d, WithNow(d.Add(-year))))

		now := d.Add(-2*day - 6*time.Hour)
		assert.Equal(t, "0 minutes", h.TimeSince(d, WithNow(now)))

		now = d.Add(-2*week - 2*day)
		assert.Equal(t, "0 minutes", h.TimeSince(d, WithNow(now)))

		now = d.Add(-2*week - 2*time.Hour - 4*time.Minute)
		assert.Equal(t, "0 minutes", h.TimeSince(d, WithNow(now)))

		now = d.Add(-4*day - 5*time.Minute)
		assert.Equal(t, "0 minutes", h.TimeSince(d, WithNow(now)))
	})

	t.Run("test different timezones", func(t *testing.T) {
		now := time.Now()
		nowDtz := time.Now().In(time.FixedZone("myzone", -210*60))

		assert.Equal(t, "0 minutes", h.TimeSince(now))
		assert.Equal(t, "0 minutes", h.TimeSince(now.In(time.UTC)))
		assert.Equal(t, "0 minutes", h.TimeSince(nowDtz))
		assert.Equal(t, "0 minutes", h.TimeSince(nowDtz, WithNow(now)))

		assert.Equal(t, "0 minutes", h.TimeUntil(now))
		assert.Equal(t, "0 minutes", h.TimeUntil(now.In(time.UTC)))
		assert.Equal(t, "0 minutes", h.TimeUntil(nowDtz))
		assert.Equal(t, "0 minutes", h.TimeUntil(nowDtz, WithNow(now)))
	})

	t.Run("test leap year", func(t *testing.T) {
		start := time.Date(2016, 12, 25, 0, 0, 0, 0, time.Local)
		assert.Equal(t, "1 week", h.TimeUntil(start.Add(week), WithNow(start)))
		assert.Equal(t, "1 week", h.TimeSince(start, WithNow(start.Add(week))))
	})

	t.Run("test leap year eve", func(t *testing.T) {
		dd := time.Date(2016, 12, 31, 0, 0, 0, 0, time.Local)
		now := time.Date(2016, 12, 31, 18, 0, 0, 0, time.Local)
		assert.Equal(t, "0 minutes", h.TimeSince(dd.Add(day), WithNow(now)))
		assert.Equal(t, "0 minutes", h.TimeUntil(dd.Add(-week), WithNow(now)))
	})

	t.Run("test thousand years ago", func(t *testing.T) {
		d := time.Date(2007, 8, 14, 13, 46, 0, 0, time.Local)
		dd := time.Date(1007, 8, 14, 13, 46, 0, 0, time.Local)
		assert.Equal(t, "1,000 years", h.TimeSince(dd, WithNow(d)))
		assert.Equal(t, "1,000 years", h.TimeUntil(d, WithNow(dd)))
	})

	t.Run("test depth", func(t *testing.T) {
		dd := d.Add(year + month + week + day + time.Hour)
		tests := []struct {
			value    time.Time
			depth    int
			expected string
		}{
			{dd, 1, "1 year"},
			{dd, 2, "1 year, 1 month"},
			{dd, 3, "1 year, 1 month, 1 week"},
			{dd, 4, "1 year, 1 month, 1 week, 1 day"},
			{dd, 5, "1 year, 1 month, 1 week, 1 day, 1 hour"},
			{dd, 6, "1 year, 1 month, 1 week, 1 day, 1 hour"},
			{d.Add(time.Hour), 5, "1 hour"},
			{d.Add(4 * time.Minute), 3, "4 minutes"},
			{d.Add(time.Hour + time.Minute), 1, "1 hour"},
			{d.Add(day + time.Hour), 1, "1 day"},
			{d.Add(week + day), 1, "1 week"},
			{d.Add(month + week), 1, "1 month"},
			{d.Add(year + month), 1, "1 year"},
			{d.Add(year + week + day), 3, "1 year"},
		}

		for _, tt := range tests {
			t.Run(tt.expected, func(t *testing.T) {
				assert.Equal(t, tt.expected, h.TimeSince(d, WithNow(tt.value), WithDepth(tt.depth)))
				assert.Equal(t, tt.expected, h.TimeUntil(tt.value, WithNow(d), WithDepth(tt.depth)))
			})
		}
	})

	t.Run("test with direct now time", func(t *testing.T) {
		now := d.Add(2*week + 3*time.Hour + 4*time.Minute)
		assert.Equal(t, "2 weeks", h.TimeSinceFrom(d, now))
		now = d.Add(4*day + 5*time.Minute)
		assert.Equal(t, "4 days", h.TimeSinceFrom(d, now))

		dd := time.Date(2016, 12, 31, 0, 0, 0, 0, time.Local)
		now = time.Date(2016, 12, 31, 18, 0, 0, 0, time.Local)
		assert.Equal(t, "0 minutes", h.TimeSinceFrom(dd.Add(day), now))
		assert.Equal(t, "0 minutes", h.TimeUntilFrom(dd.Add(-week), now))
	})
}

func TestFormatFunctions(t *testing.T) {

	t.Run("direct functions", func(t *testing.T) {
		h := createSourceHumanizer(t)
		assert.NotEmpty(t, h.Now())
		assert.NotEmpty(t, h.Date())
		assert.NotEmpty(t, h.Time())
	})
}
