package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func RunGUI() {
	a := app.New()
	w := a.NewWindow("Calendar parser")

	calendars := widget.NewSelect([]string{"Calendar 1", "Calendar 2"}, nil)
	dateFrom := widget.NewEntry()
	dateTo := widget.NewEntry()
	result := widget.NewLabel()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Calendars", Widget: calendars},
			{Text: "Write time start for parsing (format: YYYY-MM-DD): ", Widget: dateFrom},
			{Text: "Write time end for parsing (format: YYYY-MM-DD): ", Widget: dateTo},
		},
		OnSubmit: func() {

		},
	}

	w.SetContent(widget.NewLabel("Calendar parser"))
	w.ShowAndRun()
}
