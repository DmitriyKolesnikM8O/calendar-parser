package main

import (
	"flag"
)

var (
	mode = flag.Bool("gui", false, "Running in GUI mode")
)

func main() {

	flag.Parse()
	if *mode {
		RunGUI()
		return
	} else {
		RunCLI()
	}

}
