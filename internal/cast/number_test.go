package cast

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToNumber(t *testing.T) {
	s := "5.5"
	i := 10
	i8 := int8(8)
	i16 := int16(16)
	i32 := int32(32)
	i64 := int64(64)
	ui8 := uint8(8)
	ui16 := uint16(16)
	ui32 := uint32(32)
	ui64 := uint64(64)
	wd := time.Now().Weekday()

	tests := []struct {
		name    string
		arg     interface{}
		wantRes float64
		wantErr assert.ErrorAssertionFunc
	}{
		{"int", 1, 1, assert.NoError},
		{"int8", int8(1), 1, assert.NoError},
		{"int16", int16(1), 1, assert.NoError},
		{"int32", int32(1), 1, assert.NoError},
		{"int64", int64(1), 1, assert.NoError},
		{"int8 pointer", &i8, float64(i8), assert.NoError},
		{"int16 pointer", &i16, float64(i16), assert.NoError},
		{"int32 pointer", &i32, float64(i32), assert.NoError},
		{"int64 pointer", &i64, float64(i64), assert.NoError},
		{"uint", uint(5), 5, assert.NoError},
		{"uint8", uint8(5), 5, assert.NoError},
		{"uint16", uint16(5), 5, assert.NoError},
		{"uint32", uint32(5), 5, assert.NoError},
		{"uint64", uint64(5), 5, assert.NoError},
		{"uint8 pointer", &ui8, float64(ui8), assert.NoError},
		{"uint16 pointer", &ui16, float64(ui16), assert.NoError},
		{"uint32 pointer", &ui32, float64(ui32), assert.NoError},
		{"uint64 pointer", &ui64, float64(ui64), assert.NoError},
		{"int pointer", &i, float64(i), assert.NoError},
		{"float64", 1.11111, 1.11111, assert.NoError},
		{"int string", "5", 5, assert.NoError},
		{"string", s, 5.5, assert.NoError},
		{"string pointer", &s, 5.5, assert.NoError},
		{"weekday", wd, float64(wd), assert.NoError},
		{"weekday pointer", wd, float64(wd), assert.NoError},
		{"nil pointer", nil, 0, assert.Error},
		{"invalid string", "hello world", 0, assert.Error},
		{"json number", json.Number("1234"), 1234, assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := ToNumber(tt.arg)
			if !tt.wantErr(t, err, fmt.Sprintf("ToNumber(%v)", tt.arg)) {
				return
			}
			assert.Equalf(t, tt.wantRes, gotRes, "ToNumber(%v)", tt.arg)
		})
	}
}
