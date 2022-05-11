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

func (p *Message) AddReference(ref *Reference) {
	if p.Comment == nil {
		p.Comment = &Comment{}
	}

	p.Comment.AddReference(ref)
}

func (p *Message) Less(q *Message) bool {
	if p.Comment.Less(q.Comment) {
		return true
	}

	if a, b := p.Context, q.Context; a != b {
		return a < b
	}
	if a, b := p.ID, q.ID; a != b {
		return a < b
	}
	if a, b := p.IDPlural, q.IDPlural; a != b {
		return a < b
	}
	return false
}

func (p *Message) Merge(other *Message) {
	newReferences := make([]*Reference, 0)
	for _, ref := range p.Comment.References {
		for _, otherRef := range other.Comment.References {
			if ref.Line == otherRef.Line && ref.Path == otherRef.Path {
				continue
			}

			newReferences = append(newReferences, otherRef)
		}
	}
	p.Comment.References = append(p.Comment.References, newReferences...)

	newFlags := make([]string, 0)
	for _, flag := range p.Comment.Flags {
		for _, otherFlag := range other.Comment.Flags {
			if flag == otherFlag {
				continue
			}

			newFlags = append(newFlags, otherFlag)
		}
	}
	p.Comment.Flags = append(p.Comment.Flags, newFlags...)

	if p.IDPlural == "" && other.IDPlural != "" {
		p.IDPlural = other.IDPlural
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
