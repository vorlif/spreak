package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	IsVerbose       bool
	SourceDir       string
	OutputDir       string
	OutputFile      string
	CommentPrefixes []string
	ExtractErrors   bool
	ErrorContext    string

	DefaultDomain   string
	WriteNoLocation bool
	WrapWidth       int
	DontWrap        bool

	OmitHeader      bool
	CopyrightHolder string
	PackageName     string
	BugsAddress     string

	MaxDepth int
	Args     []string

	Timeout time.Duration
}

func NewDefault() *Config {
	return &Config{
		IsVerbose:       false,
		SourceDir:       "",
		OutputDir:       "./",
		OutputFile:      "",
		CommentPrefixes: []string{"TRANSLATORS"},
		ExtractErrors:   false,
		ErrorContext:    "errors",

		DefaultDomain:   "messages",
		WriteNoLocation: false,
		WrapWidth:       80,
		DontWrap:        false,

		OmitHeader:      false,
		CopyrightHolder: "THE PACKAGE'S COPYRIGHT HOLDER",
		PackageName:     "PACKAGE VERSION",
		BugsAddress:     "",

		MaxDepth: 20,
		Timeout:  15 * time.Minute,
	}
}

func (c *Config) Prepare() error {
	c.ErrorContext = strings.TrimSpace(c.ErrorContext)
	c.DefaultDomain = strings.TrimSpace(c.DefaultDomain)
	if c.DefaultDomain == "" {
		return errors.New("a default domain is required")
	}

	if c.Timeout < 1*time.Minute {
		return errors.New("the value for Timeout must be at least one minute")
	}

	currentDir, errC := os.Getwd()
	if errC != nil {
		return errC
	}

	if c.SourceDir == "" {
		c.SourceDir = currentDir
	}

	//nolint:revive
	if absP, err := filepath.Abs(c.SourceDir); err != nil {
		return err
	} else {
		c.SourceDir = absP
	}

	if c.OutputFile != "" {
		c.OutputDir = filepath.Dir(c.OutputFile)
		c.OutputFile = filepath.Base(c.OutputFile)
	} else {
		c.OutputFile = c.DefaultDomain + ".pot"
	}

	//nolint:revive
	if absP, err := filepath.Abs(c.OutputDir); err != nil {
		return err
	} else {
		c.OutputDir = absP
	}

	if c.DontWrap {
		c.WrapWidth = -1
	}

	return nil
}
