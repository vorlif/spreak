package main

import (
	"github.com/vorlif/testdata/sub"

	"github.com/vorlif/spreak/localize"
)

var _ = sub.Sub{
	Text:   "global struct msgid",
	Plural: "global struct plural",
}

type OneLineStruct struct {
	A, B, C localize.Singular
}

func structLocalTest() []*sub.Sub {

	// TRANSLATORS: Struct init
	_ = OneLineStruct{
		A: "A3",
		B: "B3",
		C: "C3",
	}

	_ = OneLineStruct{"A4", "B4", "C4"}

	item := &sub.Sub{
		Text:   "local struct msgid",
		Plural: "local struct plural",
	}

	item.Text = "struct attr assign"

	return []*sub.Sub{item}
}
