package config

import "strings"

const (
	SpreakPackagePath         = "github.com/vorlif/spreak"
	SpreakLocalizePackagePath = SpreakPackagePath + "/localize"
	XSpreakPackagePath        = SpreakPackagePath + "/xspreak"
)

func IsValidSpreakPackage(pkg string) bool {
	return pkg == SpreakPackagePath ||
		pkg == SpreakLocalizePackagePath ||
		strings.HasPrefix(pkg, XSpreakPackagePath)
}

func ShouldScanPackage(pkg string) bool {
	if ShouldScanStruct(pkg) {
		return true
	}

	// We need the definitions of the localizer and locale methods.
	return pkg == SpreakPackagePath
}

func ShouldScanStruct(pkg string) bool {
	if !strings.HasPrefix(pkg, SpreakPackagePath) {
		return true
	}

	// We need the definitions of the message. All other packages can be ignored.
	return pkg == SpreakLocalizePackagePath || strings.HasPrefix(pkg, XSpreakPackagePath)
}

func ShouldExtractPackage(pkg string) bool {
	if !strings.HasPrefix(pkg, SpreakPackagePath) {
		return true
	}

	return strings.HasPrefix(pkg, XSpreakPackagePath)
}
