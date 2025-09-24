package catalog

import (
	"encoding/json"
	"fmt"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak/v2/catalog/cldrplural"
)

func ExampleNewJSONCatalog() {
	cat := NewJSONCatalog(language.German, "domain")
	data := []byte(`{
	"help": "Hilfe",
	"car": {
		"one": "Auto",
		"other": "Autos"
	}
}`)

	if err := json.Unmarshal(data, cat); err != nil {
		panic(err)
	}

	tr, err := cat.Lookup("", "help")
	if err != nil {
		panic(err)
	}

	fmt.Println(tr)
	// Output:
	// Hilfe
}

func ExampleNewJSONCatalogWithMessages() {
	messages := make(JSONMessages)
	messages["car"] = &JSONMessage{
		Translations: map[cldrplural.Category]string{
			cldrplural.One:   "Car",
			cldrplural.Other: "Cars",
		},
	}
	messages["help_ctx"] = &JSONMessage{
		Translations: map[cldrplural.Category]string{
			cldrplural.Other: "Help",
		},
	}

	cat, err := NewJSONCatalogWithMessages(language.English, "", messages)
	if err != nil {
		panic(err)
	}
	fmt.Println(cat.Lookup("", "car"))

	res, err := json.Marshal(cat)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))

	// Output:
	// Cars <nil>
	// {"car":{"one":"Car","other":"Cars"},"help_ctx":"Help"}
}

func ExampleApplyPluralCategoriesToJSONMessage() {
	msg := &JSONMessage{
		Translations: map[cldrplural.Category]string{
			cldrplural.One:   "Car",
			cldrplural.Other: "Cars",
		},
	}

	fmt.Println("Before", msg.Translations)

	ApplyPluralCategoriesToJSONMessage(language.Polish, msg)

	fmt.Println("After", msg.Translations)

	// Output:
	// Before map[One:Car Other:Cars]
	// After map[One:Car Few: Many: Other:Cars]
}
