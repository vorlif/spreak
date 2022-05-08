package main

import alias "github.com/vorlif/spreak/localize"

func varAssignFunc() (string, string) {
	var name alias.Singular

	name = "Bob"
	name = "Bobby"

	var alternativeName string
	alternativeName = "no localize assign local"

	return name, alternativeName
}

func changeApplicationName() {
	applicationName = "application"
}

type assignStruct struct{}

func (assignStruct) testAssign() string {
	var name alias.Singular
	name = "john"
	name = "doe"
	return name
}

func assignFunc(singular alias.Singular) {
	singular = "assign function param"
}
