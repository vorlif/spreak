package localize

// MsgID is an alias type for string which is used by xspreak for extracting strings
// MsgID and Singular are synonymous and both represent the ID to identify the message.
type MsgID = string

// Singular is an alias type for string which is used by xspreak for extracting strings
// MsgID and Singular are synonymous and both represent the ID to identify the message.
type Singular = string

// Plural is an alias type for string which is used by xspreak for extracting strings.
type Plural = string

// Context is an alias type for string which is used by xspreak for extracting strings.
type Context = string

// Domain is an alias type for string which is used by xspreak for extracting strings.
type Domain = string

type Localizable interface {
	GetMsgID() string
	GetPluralID() string
	GetContext() string
	GetVars() []interface{}
	GetCount() int
	HasDomain() bool
	GetDomain() string
}

// Message is a simple struct representing a message without a domain.
// Can be converted to a translated string by Localizers or a Locale.
type Message struct {
	Singular Singular
	Plural   Plural
	Context  Context
	Vars     []interface{}
	Count    int
}

var _ Localizable = (*Message)(nil)

func (m *Message) GetMsgID() string { return m.Singular }

func (m *Message) GetPluralID() string { return m.Plural }

func (m *Message) GetContext() string { return m.Context }

func (m *Message) GetVars() []interface{} { return m.Vars }

func (m *Message) GetCount() int { return m.Count }

func (m *Message) HasDomain() bool { return false }

func (m *Message) GetDomain() string { return "" }
