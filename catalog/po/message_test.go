package po

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessage_AddReference(t *testing.T) {
	m := NewMessage()
	require.NotNil(t, m)

	t.Run("references can be added", func(t *testing.T) {
		assert.Empty(t, m.Comment.References)
		ref := &Reference{Path: "path/to/file"}
		m.AddReference(ref)
		assert.Len(t, m.Comment.References, 1)

		ref2 := &Reference{Path: "path/to/other/file"}
		m.AddReference(ref2)
		assert.Len(t, m.Comment.References, 2)
	})

	t.Run("References can be added without initialization", func(t *testing.T) {
		m.Comment = nil

		ref := &Reference{Path: "path/to/file"}
		m.AddReference(ref)
		assert.Len(t, m.Comment.References, 1)
	})
}

func TestMessage_Less(t *testing.T) {
	m := NewMessage()
	o := NewMessage()

	t.Run("Primary is sorted by references", func(t *testing.T) {
		m.AddReference(&Reference{Path: "b"})
		o.AddReference(&Reference{Path: "a"})
		assert.True(t, m.Less(o))
	})

	t.Run("After that the sorting is done according to the msgid", func(t *testing.T) {
		m.Comment = nil
		o.Comment = nil

		m.ID = "a"
		o.ID = "b"

		assert.False(t, m.Less(o))
		m.ID = "c"
		assert.True(t, m.Less(o))
	})

	t.Run("Otherwise false is returned", func(t *testing.T) {
		m.ID = ""
		o.ID = ""
		assert.False(t, m.Less(o))
		assert.False(t, o.Less(m))
	})
}

func TestMessageSort(t *testing.T) {
	msg := NewMessage()
	msg.AddReference(&Reference{Path: "b"})
	o := NewMessage()
	o.AddReference(&Reference{Path: "a"})
	messages := []*Message{msg, o}
	sort.Slice(messages, func(i, j int) bool {
		return messages[j].Less(messages[i])
	})

	assert.Equal(t, messages[0], o)
	assert.Equal(t, messages[1], msg)
}

func TestMessage_Merge(t *testing.T) {
	t.Run("add message", func(t *testing.T) {
		msg := NewMessage()
		msg.AddReference(&Reference{Path: "b"})
		o := NewMessage()
		o.AddReference(&Reference{Path: "a"})
		o.Comment.AddFlag("flag-a")

		msg.Merge(o)
		assert.Len(t, msg.Comment.References, 2)
		assert.Len(t, msg.Comment.Flags, 1)
	})

	t.Run("nil message", func(t *testing.T) {
		f := func() {
			msg := NewMessage()
			msg.Merge(nil)
		}
		assert.NotPanics(t, f)
	})

	t.Run("updates plural id", func(t *testing.T) {
		msg := NewMessage()
		msg.ID = "test"

		other := NewMessage()
		other.ID = "test"
		other.IDPlural = "test_plural"

		msg.Merge(other)
		assert.Equal(t, other.IDPlural, msg.IDPlural)
	})

	t.Run("create comments struct", func(t *testing.T) {
		msg := NewMessage()
		msg.Comment = nil
		assert.Nil(t, msg.Comment)
		msg.Merge(NewMessage())
		assert.NotNil(t, msg.Comment)
	})
}

func TestMessages_Add(t *testing.T) {
	t.Run("creates context map", func(t *testing.T) {
		ctx := "context"
		messages := make(Messages)
		assert.NotContains(t, messages, ctx)

		msg := NewMessage()
		msg.Context = ctx
		msg.ID = "test"
		messages.Add(msg)

		if assert.Contains(t, messages, ctx) {
			assert.Contains(t, messages[ctx], msg.ID)
			assert.Equal(t, messages[ctx][msg.ID], msg)
		}
	})

	t.Run("update existing", func(t *testing.T) {
		messages := make(Messages)
		msg := NewMessage()
		msg.ID = "test"

		messages.Add(msg)

		other := NewMessage()
		other.ID = msg.ID
		messages.Add(other)

		assert.Equal(t, messages[msg.Context][msg.ID], msg)
		assert.Len(t, messages, 1)
		assert.Len(t, messages[msg.Context], 1)
	})
}
