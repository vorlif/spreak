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
