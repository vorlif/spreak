package plural

import (
	"encoding/json"
	"os"
	"testing"
)

type fixture struct {
	PluralForm string
	Fixture    []int
}

func TestParser_Parse(t *testing.T) {

	f, err := os.Open("../../testdata/pluralfixtures.json")
	if err != nil {
		t.Fatal(err)
	}
	dec := json.NewDecoder(f)
	var fixtures []fixture
	err = dec.Decode(&fixtures)
	if err != nil {
		t.Fatal(err)
	}

	for _, data := range fixtures {
		forms, err := Parse(data.PluralForm)
		if err != nil {
			t.Errorf("'%s' triggered error: %s", data.PluralForm, err)
		} else if forms == nil || forms.node == nil {
			t.Logf("'%s' compiled to nil", data.PluralForm)
			t.Fail()
		} else {
			for n, e := range data.Fixture {
				i := forms.IndexForN(n)
				if i != e {
					t.Logf("'%s' with n = %d, expected %d, got %d, compiled to", data.PluralForm, n, e, i)
					t.Fail()
				}
				if i == -1 {
					break
				}
			}
		}
	}
}
