package main

import (
	"flag"

	"github.com/DmitriyKolesnikM8O/calendar-parser/internal/cli"
	"github.com/DmitriyKolesnikM8O/calendar-parser/internal/gui"
)

var (
	mode = flag.Bool("gui", false, "Running in GUI mode")
)

func main() {

	flag.Parse()
	if *mode {
		gui.RunGUI()
		return
	} else {
		cli.RunCLI()
	}

}
