package po

type Message struct {
	Comment  *Comment
	Context  string
	ID       string
	IDPlural string
	Str      map[int]string
}

type Messages map[string]map[string]*Message

func NewMessage() *Message {
	return &Message{
		Comment: &Comment{
			References: make([]*Reference, 0),
			Flags:      make([]string, 0),
		},
		Context:  "",
		ID:       "",
		IDPlural: "",
		Str:      make(map[int]string),
	}
}

func (m *Message) AddReference(ref *Reference) {
	if m.Comment == nil {
		m.Comment = &Comment{}
	}

	m.Comment.AddReference(ref)
}

func (m *Message) Less(q *Message) bool {
	if m.Comment != nil && q.Comment != nil {
		return m.Comment.Less(q.Comment)
	}

	if a, b := m.ID, q.ID; a != b {
		return a < b
	}

	return false
}

func (m *Message) Merge(other *Message) {
	newReferences := make([]*Reference, 0)
	for _, otherRef := range other.Comment.References {
		for _, ref := range m.Comment.References {
			if ref.Equal(otherRef) {
				continue
			}

			newReferences = append(newReferences, otherRef)
		}
	}
	m.Comment.References = append(m.Comment.References, newReferences...)

	newFlags := make([]string, 0)
	for _, otherFlag := range other.Comment.Flags {
		for _, flag := range m.Comment.Flags {
			if flag == otherFlag {
				continue
			}

			newFlags = append(newFlags, otherFlag)
		}
	}
	m.Comment.Flags = append(m.Comment.Flags, newFlags...)
	m.Comment.sort()

	if m.IDPlural == "" && other.IDPlural != "" {
		m.IDPlural = other.IDPlural
	}
}

func (m Messages) Add(msg *Message) {
	if _, ok := m[msg.Context]; !ok {
		m[msg.Context] = make(map[string]*Message)
	}

	if _, ok := m[msg.Context][msg.ID]; ok {
		m[msg.Context][msg.ID].Merge(msg)
	} else {
		m[msg.Context][msg.ID] = msg
	}
}
