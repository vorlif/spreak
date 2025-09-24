module example.com/spreak

go 1.25.0

require (
	github.com/Xuanwo/go-locale v1.1.3
	github.com/vorlif/spreak/v2 v2.0.0
	golang.org/x/text v0.29.0
)

require golang.org/x/sys v0.28.0 // indirect

replace github.com/vorlif/spreak/v2 v2.0.0 => ../
