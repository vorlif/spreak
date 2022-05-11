package spreak

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestErrMissingDomain(t *testing.T) {
	err := &ErrMissingDomain{}
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "spreak")

	err.Language = language.Afrikaans
	assert.Contains(t, err.Error(), err.Language.String())

	err.Domain = "mydomain"
	assert.Contains(t, err.Error(), err.Domain)
}

func TestErrMissingContext(t *testing.T) {
	err := &ErrMissingContext{}
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "spreak")

	err.Language = language.Afrikaans
	assert.Contains(t, err.Error(), err.Language.String())

	err.Domain = "mydomain"
	assert.Contains(t, err.Error(), err.Domain)

	err.Context = "mycontext"
	assert.Contains(t, err.Error(), err.Context)
}

func TestErrMissingMessageId(t *testing.T) {
	err := &ErrMissingMessageID{
		Language: language.Danish,
		Domain:   "mydomain",
		Context:  "mycontext",
		MsgID:    "mymsgid",
	}
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "spreak")

	assert.Contains(t, err.Error(), err.Language.String())
	assert.Contains(t, err.Error(), err.Domain)
	assert.Contains(t, err.Error(), err.Context)
	assert.Contains(t, err.Error(), err.MsgID)
}

func TestErrMissingTranslation(t *testing.T) {
	err := &ErrMissingTranslation{
		Language: language.Danish,
		Domain:   "mydomain",
		Context:  "mycontext",
		MsgID:    "mymsgid",
		Idx:      1022,
	}
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "spreak")

	assert.Contains(t, err.Error(), err.Language.String())
	assert.Contains(t, err.Error(), err.Domain)
	assert.Contains(t, err.Error(), err.Context)
	assert.Contains(t, err.Error(), err.MsgID)
	assert.Contains(t, err.Error(), strconv.Itoa(err.Idx))
}
