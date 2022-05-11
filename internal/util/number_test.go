package util

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToNumber(t *testing.T) {
	s := "5.5"
	i := 10
	wd := time.Now().Weekday()

	tests := []struct {
		name    string
		args    interface{}
		wantRes float64
		wantErr assert.ErrorAssertionFunc
	}{
		{"int", 1, 1, assert.NoError},
		{"int8", int8(1), 1, assert.NoError},
		{"int16", int16(1), 1, assert.NoError},
		{"int32", int32(1), 1, assert.NoError},
		{"int64", int64(1), 1, assert.NoError},
		{"uint", uint(5), 5, assert.NoError},
		{"uint8", uint8(5), 5, assert.NoError},
		{"uint16", uint16(5), 5, assert.NoError},
		{"uint32", uint32(5), 5, assert.NoError},
		{"uint64", uint64(5), 5, assert.NoError},
		{"int pointer", &i, float64(i), assert.NoError},
		{"int string", "5", 5, assert.NoError},
		{"string", s, 5.5, assert.NoError},
		{"string pointer", &s, 5.5, assert.NoError},
		{"weekday", wd, float64(wd), assert.NoError},
		{"weekday pointer", wd, float64(wd), assert.NoError},
		{"nil pointer", nil, 0, assert.Error},
		{"invalid string", "hello world", 0, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := ToNumber(tt.args)
			if !tt.wantErr(t, err, fmt.Sprintf("ToNumber(%v)", tt.args)) {
				return
			}
			assert.Equalf(t, tt.wantRes, gotRes, "ToNumber(%v)", tt.args)
		})
	}
}
