package main

import (
	"flag"
	"log"

	"github.com/DmitriyKolesnikM8O/calendar-parser/internal/cli"
	"github.com/DmitriyKolesnikM8O/calendar-parser/internal/gui"
)

var (
	mode     = flag.Bool("gui", false, "Running in GUI mode")
	charType = flag.String("chart", "bar", "Char type: bar or pie")
)

func main() {

	flag.Parse()
	if *mode {
		err := gui.RunGUI(*charType)
		if err != nil {
			log.Fatalf(err.Error())
		}
	} else {
		err := cli.RunCLI(*charType)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

}
