package main

import (
	"flag"
	"log"

	"github.com/DmitriyKolesnikM8O/calendar-parser/internal/cli"
	"github.com/DmitriyKolesnikM8O/calendar-parser/internal/gui"
)

var (
	mode = flag.Bool("gui", false, "Running in GUI mode")
)

func main() {

	flag.Parse()
	if *mode {
		err := gui.RunGUI()
		if err != nil {
			log.Fatalf(err.Error())
		}
	} else {
		err := cli.RunCLI()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

}
