package main

import (
	"errors"
	alias "errors"
)

/* TRANSLATORS: comment a
comment b
comment c */
// comment d
var ErrGlobal = errors.New(
	// test comment
	"global error")

// TRANSLATORS:  comment e
// comment f
//
// comment g
var ErrAliasGloba = alias.New("global alias error")

func localError() error {
	// comment local error
	_ = errors.New("local error")
	// comment local alias error
	_ = errors.New("local alias error")

	return errors.New("return error")
}
