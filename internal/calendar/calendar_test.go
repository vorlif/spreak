package calendar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isLeapYear(t *testing.T) {
	for _, year := range []int{2004, 2008, 2012, 2016, 2020, 2024, 2028, 2032, 2036, 2040} {
		assert.True(t, IsLeap(year))
	}

	for _, year := range []int{2005, 2006, 2007, 1985, 2034} {
		assert.False(t, IsLeap(year), year)
	}
}

func TestLeapDays(t *testing.T) {
	assert.Equal(t, 2, LeapDays(2016, 2022))
	assert.Equal(t, -2, LeapDays(2022, 2016))
	assert.Equal(t, 0, LeapDays(2001, 2003))
}

func TestDaysInMonth(t *testing.T) {
	assert.Equal(t, 31, DaysInMonth(2022, 05))
	assert.Equal(t, 29, DaysInMonth(2024, 02))
}
