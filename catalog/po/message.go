package po

// Message is a representation of a single message in a catalog.
type Message struct {
	Comment  *Comment
	Context  string
	ID       string
	IDPlural string
	Str      map[int]string
}

func NewMessage() *Message {
	return &Message{
		Comment: &Comment{
			References: make([]*Reference, 0),
			Flags:      make([]string, 0),
		},
		Context:  "",
		ID:       "",
		IDPlural: "",
		Str:      make(map[int]string, 1),
	}
}

func (m *Message) AddReference(ref *Reference) {
	if m.Comment == nil {
		m.Comment = NewComment()
	}

	m.Comment.AddReference(ref)
}

func (m *Message) Merge(other *Message) {
	if other == nil {
		return
	}

	if m.Comment == nil {
		m.Comment = NewComment()
	}
	m.Comment.Merge(other.Comment)

	if m.IDPlural == "" && other.IDPlural != "" {
		m.IDPlural = other.IDPlural
	}
}
