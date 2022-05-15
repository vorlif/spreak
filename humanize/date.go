package humanize

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/vorlif/spreak/internal/calendar"
	"github.com/vorlif/spreak/internal/util"
)

const (
	DateFormat          = "DATE_FORMAT"
	TimeFormat          = "TIME_FORMAT"
	DateTimeFormat      = "DATETIME_FORMAT"
	YearMonthFormat     = "YEAR_MONTH_FORMAT"
	MonthDayFormat      = "MONTH_DAY_FORMAT"
	ShortDateFormat     = "SHORT_DATE_FORMAT"
	ShortDatetimeFormat = "SHORT_DATETIME_FORMAT"
)

type gettextEntry struct {
	context  string
	singular string
	plural   string
}

// NaturalDay returns for time values that are tomorrow, today or yesterday compared to present day
// the representing string.
func (h *Humanizer) NaturalDay(i interface{}) string {
	t, err := util.ToTime(i)
	if err != nil {
		return fmt.Sprintf("%v", i)
	}

	value := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	now := time.Now().In(value.Location())
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	delta := value.Sub(today).Hours() / 24
	switch delta {
	case 0:
		return h.loc.Get("today")
	case 1:
		return h.loc.Get("tomorrow")
	case -1:
		return h.loc.Get("yesterday")
	default:
		return t.String()
	}
}

var naturalTimeStrings = map[string]gettextEntry{
	// Translators: delta will contain a string like '2 months' or '1 month, 2 weeks'
	"past-day": {"", "%[1]v ago", ""},
	// Translators: please keep a non-breaking space (U+00A0) between count
	// and time unit.
	"past-hour": {"", "an hour ago", "%[1]v hours ago"},
	// Translators: please keep a non-breaking space (U+00A0) between count
	// and time unit.
	"past-minute": {"", "a minute ago", "%[1]v minutes ago"},
	// Translators: please keep a non-breaking space (U+00A0) between count
	// and time unit.
	"past-second": {"", "a second ago", "%[1]v seconds ago"},
	"now":         {"", "now", ""},
	// Translators: please keep a non-breaking space (U+00A0) between count
	// and time unit.
	"future-second": {"", "a second from now", "%[1]v seconds from now"},
	// Translators: please keep a non-breaking space (U+00A0) between count
	// and time unit.
	"future-minute": {"", "a minute from now", "%[1]v minutes from now"},
	// Translators: please keep a non-breaking space (U+00A0) between count
	// and time unit.
	"future-hour": {"", "an hour from now", "%[1]v hours from now"},
	// Translators: delta will contain a string like '2 months' or '1 month, 2 weeks'
	"future-day": {"", "%[1]v from now", ""},
}

var naturalPastSubstrings = map[string]gettextEntry{
	"year":   {"naturaltime-past", "%[1]v year", "%[1]v years"},
	"month":  {"naturaltime-past", "%[1]v month", "%[1]v months"},
	"week":   {"naturaltime-past", "%[1]v week", "%[1]v weeks"},
	"day":    {"naturaltime-past", "%[1]v day", "%[1]v days"},
	"hour":   {"naturaltime-past", "%[1]v hour", "%[1]v hours"},
	"minute": {"naturaltime-past", "%[1]v minute", "%[1]v minutes"},
}

var naturalFutureSubstrings = map[string]gettextEntry{
	// Translators: 'naturaltime-future' strings will be included in
	//  '%[1]v from now'.
	"year":   {"naturaltime-future", "%[1]v year", "%[1]v years"},
	"month":  {"naturaltime-future", "%[1]v month", "%[1]v months"},
	"week":   {"naturaltime-future", "%[1]v week", "%[1]v weeks"},
	"day":    {"naturaltime-future", "%[1]v day", "%[1]v days"},
	"hour":   {"naturaltime-future", "%[1]v hour", "%[1]v hours"},
	"minute": {"naturaltime-future", "%[1]v minute", "%[1]v minutes"},
}

// NaturalTime shows for a time value how many seconds, minutes, or hours ago
// compared to current timestamp return representing string.
func (h *Humanizer) NaturalTime(i interface{}) string {
	t, err := util.ToTime(i)
	if err != nil {
		return fmt.Sprintf("%v", i)
	}

	now := time.Now().In(t.Location())
	if t.Before(now) {
		delta := now.Sub(t)
		deltaSec := int64(delta.Truncate(time.Second).Seconds())
		if int64(delta.Round(time.Second).Hours()) >= 24 {
			entry := naturalTimeStrings["past-day"]
			timeSince := h.TimeSinceFrom(t, now, withTimeStrings(naturalPastSubstrings))
			return h.loc.Getf(entry.singular, timeSince)
		} else if deltaSec == 0 {
			entry := naturalTimeStrings["now"]
			return h.loc.Get(entry.singular)
		} else if deltaSec < 60 {
			entry := naturalTimeStrings["past-second"]
			return h.loc.NGetf(entry.singular, entry.plural, deltaSec, deltaSec)
		} else if floorDivision(delta.Round(time.Second).Seconds(), 60) < 60 {
			count := int64(math.Floor(float64(deltaSec) / 60))
			entry := naturalTimeStrings["past-minute"]
			return h.loc.NGetf(entry.singular, entry.plural, count, count)
		} else {
			count := int64(math.Floor(math.Floor(float64(deltaSec)/60) / 60))
			entry := naturalTimeStrings["past-hour"]
			return h.loc.NGetf(entry.singular, entry.plural, count, count)
		}
	}

	delta := t.Sub(now).Round(time.Second)
	deltaSec := int64(delta.Truncate(time.Second).Seconds())
	if int64(delta.Round(time.Second).Hours()) >= 24 {
		entry := naturalTimeStrings["future-day"]
		timeSince := h.TimeUntilFrom(t, now, withTimeStrings(naturalFutureSubstrings))
		return h.loc.Getf(entry.singular, timeSince)
	} else if deltaSec == 0 {
		entry := naturalTimeStrings["now"]
		return h.loc.Get(entry.singular)
	} else if deltaSec < 60 {
		entry := naturalTimeStrings["future-second"]
		return h.loc.NGetf(entry.singular, entry.plural, deltaSec, deltaSec)
	} else if floorDivision(delta.Round(time.Second).Seconds(), 60) < 60 {
		count := int64(math.Floor(float64(deltaSec) / 60))
		entry := naturalTimeStrings["future-minute"]
		return h.loc.NGetf(entry.singular, entry.plural, count, count)
	} else {
		count := int64(math.Floor(math.Floor(float64(deltaSec)/60) / 60))
		entry := naturalTimeStrings["future-hour"]
		return h.loc.NGetf(entry.singular, entry.plural, count, count)
	}
}

var timeSinceStrings = map[string]gettextEntry{
	"year":   {"", "%[1]v year", "%[1]v years"},
	"month":  {"", "%[1]v month", "%[1]v months"},
	"week":   {"", "%[1]v week", "%[1]v weeks"},
	"day":    {"", "%[1]v day", "%[1]v days"},
	"hour":   {"", "%[1]v hour", "%[1]v hours"},
	"minute": {"", "%[1]v minute", "%[1]v minutes"},
}

var timeSinceChunks = []struct {
	name    string
	seconds int64
}{
	{"year", 60 * 60 * 24 * 365},
	{"month", 60 * 60 * 24 * 30},
	{"week", 60 * 60 * 24 * 7},
	{"day", 60 * 60 * 24},
	{"hour", 60 * 60},
	{"minute", 60},
}

type timeSinceOptions struct {
	now         time.Time
	reverse     bool
	timeStrings map[string]gettextEntry
	depth       int
}

type TimeSinceOption func(opt *timeSinceOptions)

func WithDepth(depth int) TimeSinceOption {
	return func(opt *timeSinceOptions) {
		opt.depth = depth
	}
}

func withTimeStrings(timeStrings map[string]gettextEntry) TimeSinceOption {
	return func(opt *timeSinceOptions) {
		opt.timeStrings = timeStrings
	}
}

func WithReverse(reverse bool) TimeSinceOption {
	return func(opt *timeSinceOptions) {
		opt.reverse = reverse
	}
}

func WithNow(now time.Time) TimeSinceOption {
	return func(opt *timeSinceOptions) {
		opt.now = now
	}
}

// TimeSince take two datetime objects and return the time between d and now as a nicely
// formatted string, e.g. "10 minutes". If d occurs after now, return
// "0 minutes".
//
// Units used are years, months, weeks, days, hours, and minutes.
// Seconds and microseconds are ignored. Up to `depth` adjacent units will be
// displayed.  For example, "2 weeks, 3 days" and "1 year, 3 months" are
// possible outputs, but "2 weeks, 3 hours" and "1 year, 5 days" are not.
//
// `depth` is an optional integer to control the number of adjacent time
// units returned.
//
// Adapted from
// https://web.archive.org/web/20060617175230/http://blog.natbat.co.uk/archive/2003/Jun/14/time_since
func (h *Humanizer) TimeSince(inputTime interface{}, opts ...TimeSinceOption) string {
	d, err := util.ToTime(inputTime)
	if err != nil {
		return fmt.Sprintf("%v", inputTime)
	}

	o := &timeSinceOptions{
		reverse:     false,
		timeStrings: nil,
		depth:       -1,
	}
	for _, opt := range opts {
		opt(o)
	}

	if o.timeStrings == nil || len(o.timeStrings) == 0 {
		o.timeStrings = timeSinceStrings
	}
	if o.depth <= 0 {
		o.depth = 2
	}
	now := o.now
	if now.IsZero() {
		now = time.Now()
	}

	if o.reverse {
		d, now = now, d
	}

	delta := now.In(time.UTC).Unix() - d.In(time.UTC).Unix()

	// Deal with leapyears by subtracing the number of leapdays
	leapdays := calendar.LeapDays(d.Year(), now.Year())
	if leapdays != 0 {
		if calendar.IsLeap(d.Year()) {
			leapdays--
		} else if calendar.IsLeap(now.Year()) {
			leapdays++
		}
	}
	delta -= 60 * 60 * 24 * int64(leapdays)

	// ignore microseconds
	if delta <= 0 {
		//  d is in the future compared to now, stop processing.
		entry := o.timeStrings["minute"]
		return h.loc.NPGetf(entry.context, entry.singular, entry.plural, 0, 0)
	}

	since := delta

	i := sort.Search(len(timeSinceChunks), func(i int) bool {
		chunk := timeSinceChunks[i]
		count := floorDivision(float64(since), float64(chunk.seconds))
		return count != 0
	})

	if i == len(timeSinceChunks) {
		entry := o.timeStrings["minute"]
		return h.loc.NPGetf(entry.context, entry.singular, entry.plural, 0, 0)
	}

	var result []string
	currentDepth := 0
	for i < len(timeSinceChunks) && currentDepth < o.depth {
		chunk := timeSinceChunks[i]
		count := floorDivision(float64(since), float64(chunk.seconds))
		if count <= 0 {
			break
		}
		entry := o.timeStrings[chunk.name]
		result = append(result, h.loc.NPGetf(entry.context, entry.singular, entry.plural, count, count))
		since -= chunk.seconds * count
		currentDepth++
		i++
	}

	return strings.Join(result, h.loc.Get(", "))
}

func (h *Humanizer) TimeSinceFrom(d time.Time, now time.Time, opts ...TimeSinceOption) string {
	opts = append(opts, WithNow(now))
	return h.TimeSince(d, opts...)
}

// TimeUntil works like TimeSince, but return a string measuring the time until the given time.
func (h *Humanizer) TimeUntil(d interface{}, opts ...TimeSinceOption) string {
	parsedTime, err := util.ToTime(d)
	if err != nil {
		return fmt.Sprintf("%v", d)
	}

	opts = append(opts, WithReverse(true))
	return h.TimeSince(parsedTime, opts...)
}

func (h *Humanizer) TimeUntilFrom(d time.Time, now time.Time, opts ...TimeSinceOption) string {
	opts = append(opts, WithNow(now), WithReverse(true))
	return h.TimeSince(d, opts...)
}

func (h *Humanizer) FormatTime(t time.Time, format string) string {
	switch format {
	case DateFormat:
		format = h.format.DateFormat
	case TimeFormat:
		format = h.format.TimeFormat
	case DateTimeFormat:
		format = h.format.DateTimeFormat
	case YearMonthFormat:
		format = h.format.YearMonthFormat
	case MonthDayFormat:
		format = h.format.MonthDayFormat
	case ShortDateFormat:
		format = h.format.ShortDateFormat
	case ShortDatetimeFormat:
		format = h.format.ShortDatetimeFormat
	}

	tf := newTimeFormatter(h, t)
	return tf.format(format)
}

func (h *Humanizer) Now() string {
	return h.FormatTime(time.Now(), DateTimeFormat)
}

func (h *Humanizer) Time() string {
	return h.FormatTime(time.Now(), TimeFormat)
}

func (h *Humanizer) Date() string {
	return h.FormatTime(time.Now(), DateFormat)
}
