package main

import (
	"github.com/vorlif/testdata/sub"

	"github.com/vorlif/spreak/localize"
)

// TRANSLATORS: This is not extracted
var globalSlice = []localize.Singular{
	// TRANSLATORS: numbers 1 to 4
	// only for "one" extracted
	"one", "two", "three", "four",
}

var globalStructSlice = []localize.Message{
	{
		// TRANSLATORS: For singular extracted
		Singular: "global struct slice singular",
		Plural:   "global struct slice plural",
		Context:  "global ctx",
		Vars:     nil,
		Count:    0,
	},
	{
		"global struct slice singular 2",
		// TRANSLATORS: For plural extracted
		"global struct slice plural 2",
		"",
		nil,
		0,
	},
}

func localSliceFunc() []string {
	globalSlice = append(globalSlice, "five")

	localSlice := []localize.Singular{"six", "seven", "eight", "nine"}
	localSlice = append(localSlice, "ten")

	_ = []localize.Message{
		{
			Singular: "local struct slice singular",
			Plural:   "local struct slice plural",
			Context:  "local ctx",
			Vars:     nil,
			Count:    0,
		},
		{
			"local struct slice singular 2",
			"local struct slice plural 2",
			"local struct slice ctx 2",
			nil,
			0,
		},
	}

	subs := []*sub.Sub{
		{
			Text:   "struct slice msgid",
			Plural: "struct slice plural",
		},
	}

	_ = []OneLineStruct{
		{
			A: "A1",
			B: "B1",
			C: "C1",
		},
		{"A2", "B2", "C2"},
	}

	_ = append(subs,
		&sub.Sub{
			Text:   "struct msgid arr1",
			Plural: "struct plural arr1",
		},
		&sub.Sub{
			Text:   "struct msgid arr2",
			Plural: "struct plural arr2",
		},
	)

	return localSlice
}
