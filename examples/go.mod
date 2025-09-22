module example.com/spreak

go 1.21

toolchain go1.21.4

require (
	github.com/Xuanwo/go-locale v1.1.0
	github.com/vorlif/spreak v0.0.0
	golang.org/x/text v0.14.0
)

require golang.org/x/sys v0.5.0 // indirect

replace github.com/vorlif/spreak v0.0.0 => ../
