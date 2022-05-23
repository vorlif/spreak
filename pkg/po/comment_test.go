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

		c.sort()
		require.Len(t, c.References, 6)
		assert.Equal(t, refs[2], c.References[0])
		assert.Equal(t, refs[4], c.References[1])
		assert.Equal(t, refs[1], c.References[2])
		assert.Equal(t, refs[0], c.References[3])
		assert.Equal(t, refs[3], c.References[4])
		assert.Equal(t, refs[5], c.References[5])
	})
}

func TestCommentMerge(t *testing.T) {
	t.Run("ignore nil comment", func(t *testing.T) {
		f := func() {
			c := NewComment()
			c.Merge(nil)
		}

		assert.NotPanics(t, f)
	})

	t.Run("There are no duplicate references", func(t *testing.T) {
		c := NewComment()
		c.AddReference(&Reference{Path: "a"})
		c.AddReference(&Reference{Path: "b"})
		c.AddReference(&Reference{Path: "c"})
		other := NewComment()
		dupRef := &Reference{Path: "a"}
		other.AddReference(dupRef)
		other.AddReference(&Reference{Path: "d"})
		other.AddReference(&Reference{Path: "e"})

		assert.Len(t, c.References, 3)
		c.mergeReferences(other)
		assert.Len(t, c.References, 5)
		assert.Len(t, other.References, 3)
	})

	t.Run("only new references", func(t *testing.T) {
		c := NewComment()
		other := NewComment()
		other.AddReference(&Reference{Path: "d"})
		other.AddReference(&Reference{Path: "e"})

		assert.Len(t, c.References, 0)
		c.mergeReferences(other)
		assert.Len(t, c.References, 2)
	})

	t.Run("There are no duplicate flags", func(t *testing.T) {
		c := NewComment()
		c.AddFlag("a")
		c.AddFlag("b")
		c.AddFlag("c")
		other := NewComment()
		other.AddFlag("a")
		other.AddFlag("d")
		other.AddFlag("e")

		assert.Len(t, c.Flags, 3)
		c.Merge(other)
		assert.Len(t, c.Flags, 5)
		assert.Len(t, other.Flags, 3)
	})

	t.Run("only new flags", func(t *testing.T) {
		c := NewComment()
		other := NewComment()
		other.AddFlag("d")
		other.AddFlag("e")

		assert.Len(t, c.Flags, 0)
		c.Merge(other)
		assert.Len(t, c.Flags, 2)
	})

	t.Run("update extracted and translator", func(t *testing.T) {
		comm := NewComment()
		comm.Extracted = "\na\nb\nc"
		comm.Translator = "x\ny"

		other := NewComment()
		other.Extracted = "d\nb\ne\nf"
		other.Translator = "z\nx\n"

		comm.Merge(other)
		assert.Equal(t, "a\nb\nc\nd\ne\nf", comm.Extracted)
		assert.Equal(t, "x\ny\nz", comm.Translator)
	})
}

func Test_mergeStringArrays(t *testing.T) {
	left := []string{"a", "b", "c"}
	right := []string{"d", "b", "e", "f"}
	want := []string{"a", "b", "c", "d", "e", "f"}
	assert.Equal(t, want, mergeStringArrays(left, right))

	left = []string{"z", "x", ""}
	right = []string{"x", "y"}
	want = []string{"z", "x", "", "y"}
	assert.Equal(t, want, mergeStringArrays(left, right))
}