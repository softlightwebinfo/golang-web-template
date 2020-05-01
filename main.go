package main

import (
	"codeunic.com/DocumentationApp/apps"
)

func main() {
	main := apps.Main{}
	main.Initialize()
	main.MidGzip()
	main.MidLogger(true)
	main.MidTrailingSlash(true)
	main.Cache(true)
	main.Prometheus()
	main.Recover()
	main.Start()
}
