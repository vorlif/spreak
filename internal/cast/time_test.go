package cast

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToTime(t *testing.T) {
	now := time.Now()

	t.Run("test errors", func(t *testing.T) {
		now := time.Now()
		tests := []struct {
			name    string
			arg     interface{}
			wantErr assert.ErrorAssertionFunc
		}{
			{"time", now, assert.NoError},
			{"zero time", time.Time{}, assert.NoError},
			{"positive duration", 5 * time.Minute, assert.NoError},
			{"negative duration", -5 * time.Hour, assert.NoError},
			{"number to time", 1_000, assert.NoError},
			{"pointer", &now, assert.NoError},
			{"nil", nil, assert.Error},
			{"invalid time", "string", assert.Error},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				d, err := ToTime(tt.arg)
				tt.wantErr(t, err, fmt.Sprintf("ToTime(%v)", tt.arg))
				assert.NotNil(t, d)
			})
		}
	})

	t.Run("test simple time input", func(t *testing.T) {
		d, err := ToTime(now)
		assert.NoError(t, err)
		assert.Equal(t, now, d)
	})

	t.Run("test pointer", func(t *testing.T) {
		d, err := ToTime(&now)
		assert.NoError(t, err)
		assert.Equal(t, now, d)
	})

	t.Run("test numbers", func(t *testing.T) {
		d, err := ToTime(1_000)
		assert.NoError(t, err)
		assert.Equal(t, int64(1000), d.Unix())
	})

	t.Run("test empty", func(t *testing.T) {
		d, err := ToTime(time.Time{})
		assert.NoError(t, err)
		assert.True(t, d.IsZero())
		assert.Equal(t, time.Time{}, d)
	})

	t.Run("test duration", func(t *testing.T) {
		d, err := ToTime(5 * time.Minute)
		assert.NoError(t, err)
		mins := time.Until(d).Minutes()
		assert.Truef(t, mins >= 4.75 && mins <= 5.25, fmt.Sprintf("%v", mins))

		d, err = ToTime(-10 * time.Hour)
		assert.NoError(t, err)
		hours := time.Since(d).Hours()
		assert.Truef(t, hours >= 9.9 && hours <= 10.1, fmt.Sprintf("%v", hours))
	})
}
