package main

import (
	"flag"
)

func main() {
	mode := flag.Bool("gui", false, "Running in GUI mode")

	flag.Parse()
	if *mode {
		RunGUI()
		return
	} else {
		RunCLI()
	}

}
