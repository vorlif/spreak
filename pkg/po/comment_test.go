package po

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReference(t *testing.T) {
	ref := Reference{Path: "abc"}
	assert.Equal(t, "abc", ref.String())

	ref.Line = 10
	assert.Contains(t, ref.String(), "10")

	ref.Line = 20
	assert.Contains(t, ref.String(), "20")
}

func TestComment_AddReference(t *testing.T) {
	c := NewComment()
	require.NotNil(t, c)
	c.References = nil
	assert.Nil(t, c.References)

	ref := &Reference{Path: "bc"}
	c.AddReference(ref)
	assert.NotNil(t, c.References)
	assert.Contains(t, c.References, ref)

	ref = &Reference{Path: "aa"}
	c.AddReference(ref)
	assert.Equal(t, 2, len(c.References))
	assert.Equal(t, "aa", c.References[0].Path)
}

func TestComment_HasFlag(t *testing.T) {
	flag := "test"
	c := NewComment()
	assert.False(t, c.HasFlag(flag))
	c.AddFlag(flag)
	assert.True(t, c.HasFlag(flag))
}

func TestComment_AddFlag(t *testing.T) {
	flag := "test"
	c := NewComment()
	assert.False(t, c.HasFlag(flag))
	c.AddFlag(flag)
	c.AddFlag(flag)
	c.AddFlag(flag)
	c.AddFlag(flag)
	assert.True(t, c.HasFlag(flag))
	assert.Equal(t, 1, len(c.Flags))
}

func TestComment_Less(t *testing.T) {
	c := NewComment()
	cRef := &Reference{Path: "b"}
	c.AddReference(cRef)
	o := NewComment()
	oRef := &Reference{Path: "a"}
	o.AddReference(oRef)

	t.Run("First sorted by the path", func(t *testing.T) {
		assert.True(t, c.Less(o))

		oRef.Path = "c"
		assert.False(t, c.Less(o))
	})

	t.Run("If the path is the equal, it is sorted by the line number", func(t *testing.T) {
		cRef.Path = "a"
		cRef.Line = 2
		oRef.Path = "a"
		oRef.Line = 1
		assert.True(t, c.Less(o))
		oRef.Line = 3
		assert.False(t, c.Less(o))
	})

	t.Run("If the path and the row number are the same, the sorting is done by the column.", func(t *testing.T) {
		cRef.Line = 1
		oRef.Line = 1
		cRef.Column = 2
		oRef.Column = 1
		assert.True(t, c.Less(o))
		oRef.Line = 3
		assert.False(t, c.Less(o))
	})

	t.Run("Different number of references are ignored", func(t *testing.T) {
		cRef.Column = 1
		oRef.Column = 1
		cRef.Path = "l"
		c.AddReference(&Reference{Path: "z"})
		c.AddReference(&Reference{Path: "m"})

		assert.True(t, c.Less(o))
	})
}

func TestCommentSort(t *testing.T) {
	t.Run("References are sorted", func(t *testing.T) {
		refs := []*Reference{
			{Path: "b", Line: 2, Column: 10},
			{Path: "b", Line: 1, Column: 20},
			{Path: "a", Line: 1, Column: 10},
			{Path: "c", Line: 2, Column: 20},
			{Path: "a", Line: 1, Column: 20},
			{Path: "d", Line: 1, Column: 10},
		}
		c := NewComment()

		for _, ref := range refs {
			c.AddReference(ref)
		}

		c.sortReferences()
		require.Len(t, c.References, 6)
		assert.Equal(t, refs[2], c.References[0])
		assert.Equal(t, refs[4], c.References[1])
		assert.Equal(t, refs[1], c.References[2])
		assert.Equal(t, refs[0], c.References[3])
		assert.Equal(t, refs[3], c.References[4])
		assert.Equal(t, refs[5], c.References[5])
	})
}
