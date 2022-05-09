# xspreak

xspreak automatically extracts strings that use a string alias from the localize package.
The extracted strings are stored in a `.pot` file and can then be easily translated.
The translation produces `.po` or `.mo` files which are used by spreak for looking up translations.

The extracted strings can then be passed to a [Localizer](https://pkg.go.dev/github.com/vorlif/spreak#Localizer) 
or a [Locale](https://pkg.go.dev/github.com/vorlif/spreak#Locale) which returns the matching translation.

Example:

```go
package main

import (
	"fmt"

	"golang.org/x/text/language"

	"github.com/vorlif/spreak"
	"github.com/vorlif/spreak/localize"
)

// This string is extracted because the type is localize.Singular
var ApplicationName localize.Singular = "Beautiful app"

func main() {
	bundle, err := spreak.NewBundle(
		spreak.WithSourceLanguage(language.English),
		spreak.WithDomainPath(spreak.NoDomain, "../locale"),
		spreak.WithLanguage(language.German, language.Spanish, language.Chinese),
	)
	if err != nil {
		panic(err)
	}

	t := spreak.NewLocalizer(bundle, language.Spanish)

	// Message lookup of the extracted string
	fmt.Println(t.Get(ApplicationName))
	// Output:
	// Hermosa app
}
```

## How to install

```bash
go install github.com/vorlif/spreak/xspreak@main
xspreak -help
```

## What can be extracted?

### Global variables and constants

Global variables and constants are extracted if the type is localize.Singular or localize.MsgID.
Thereby localize.Singular and localize.MsgID are always equivalent and can be used synonymously.

```go
package main

import "github.com/vorlif/spreak/localize"

const Weekday localize.Singular = "weekday"

var ApplicationName localize.Singular = "app"
```

### Local variables
Local variables are extracted if the type is localize.Singular or localize.MsgID.

```go
package main

import "github.com/vorlif/spreak/localize"

func init() {
	holiday := localize.Singular("Christmas")
}
```

### Variable assignments

Assignments to variables are extracted if the type is localize.Singular or localize.MsgID.

```go
package main

import "github.com/vorlif/spreak/localize"

var ApplicationName = "app"

func init() {
	var holiday localize.Singular

	holiday = "Mother's Day"

	ApplicationName = "App for you"
}
```

### Argument of function calls

Function calls to **global functions** are extracted if the parameter type is from the localize package.
The parameters of a function are grouped together to form a message.
Thus a message can be created with singular, plural, a context and a domain.

```go
package main

import "github.com/vorlif/spreak/localize"

func noop(name localize.Singular, plural localize.Plural, ctx localize.Context) {}

func init() {
	// Extracted as a message with singular, plural and a context
	noop("I have %d car", "I have %d cars", "cars")
}
```

### Attributes at struct initialization

Struct initializations are extracted if the struct was **defined globally** and
the attribute type *comes from the localize package*.
The attributes of a struct are grouped together to create a message.
Thus a message can be created with singular, plural, a context and a domain.

```go
package main

import "github.com/vorlif/spreak/localize"

type MyMessage struct {
	// Defined as singular and plural
	Text   localize.Singular
	Plural localize.Plural
	Tmp    string
}

func main() {
	msg := &MyMessage{
		// Extracted as a message with singular and plural
		Text:   "Hello planet",
		Plural: "Hello planets",

		// not extracted - type string
		Tmp: "tmp",
	}
}
```

### Values from an array initialization

Arrays are extracted if the type is localize.Singular or a struct that contains parameter
types from the localize package.

```go
package main

import "github.com/vorlif/spreak/localize"

var weekdays = []localize.MsgID{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

type MyMessage struct {
	Text   localize.Singular
	Plural localize.Plural
	Tmp    string
}

func main() {
	animals := []MyMessage{
		{Text: "%d dog", Plural: "%d dogs"},
		{Text: "%d cat", Plural: "%d cat"},
		{Text: "%d horse", Plural: "%d horses"},
	}
}
```

### Error texts

Strings can be extracted from `errors.New` if xspreak is called with the `-e` option.

```go
package main

import "errors"

var ErrInvalidAnimal = errors.New("this is not a valid animal")

```

### Comments

Comments can be left for translators.
These are extracted, stored in the `.pot` file and displayed to the translator.

```go
package main

import "github.com/vorlif/spreak/localize"

// TRANSLATORS: This comment is automatically extracted by xspreak
// and can be used to leave useful hints for the translators.
//
// This comment is not extracted because a blank line was inserted above it.
const InvalidName localize.Singular = "The name has an invalid format"
```

### Exclude from extraction

Strings can be ignored.

```go
package main

import "github.com/vorlif/spreak/localize"

// xspreak: ignore
const MagicName localize.Singular = ".%$($§($(%"
```

## Roadmap

* [ ] Map initialization
* [ ] `append` to array
* [ ] String backtracing