package gui

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"

	"fyne.io/fyne/v2/widget"
	"github.com/DmitriyKolesnikM8O/calendar-parser/internal/connection"
	"github.com/DmitriyKolesnikM8O/calendar-parser/internal/statistics"
	"google.golang.org/api/calendar/v3"
)

func RunGUI() {
	a := app.New()
	w := a.NewWindow("Calendar Parser")
	w.Resize(fyne.NewSize(900, 600))

	srv, calendars, err := connection.InitializeCalendarService(w)
	if err != nil {
		return
	}

	calendarOptions := make([]string, 0, len(calendars.Items))
	calendarMap := make(map[string]string)
	for _, item := range calendars.Items {
		displayName := fmt.Sprintf("%s (%s)", item.Summary, item.Id)
		calendarOptions = append(calendarOptions, displayName)
		calendarMap[displayName] = item.Id
	}

	calendarsSelect := widget.NewSelect(calendarOptions, nil)
	if len(calendarOptions) > 0 {
		calendarsSelect.SetSelected(calendarOptions[0])
	}
	dateFrom := widget.NewEntry()
	dateFrom.SetPlaceHolder("YYYY-MM-DD")
	dateTo := widget.NewEntry()
	dateTo.SetPlaceHolder("YYYY-MM-DD")
	resultLabel := widget.NewLabel("Results will appear here after submission.")
	submitButton := widget.NewButton("Parse Calendar", nil)
	viewChatButton := widget.NewButton("View chart", nil)
	viewChatButton.Disable()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Select Calendar", Widget: calendarsSelect},
			{Text: "Start Date (YYYY-MM-DD)", Widget: dateFrom},
			{Text: "End Date (YYYY-MM-DD)", Widget: dateTo},
		},
	}

	submitButton.OnTapped = func() {

		_, err := time.Parse("2006-01-02", dateFrom.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid start date format: %v", err), w)
			return
		}
		_, err = time.Parse("2006-01-02", dateTo.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid end date format: %v", err), w)
			return
		}

		calendarID := calendarMap[calendarsSelect.Selected]
		eventsColorTime := make(map[string][]struct {
			Start    *calendar.EventDateTime
			End      *calendar.EventDateTime
			Duration time.Duration
		})

		pageToken := ""
		for {
			events, err := srv.Events.List(calendarID).
				TimeMin(dateFrom.Text + "T00:00:00+03:00").
				TimeMax(dateTo.Text + "T23:59:59+03:00").
				PageToken(pageToken).
				Do()
			if err != nil {
				dialog.ShowError(fmt.Errorf("error retrieving events: %v", err), w)
				return
			}

			for _, event := range events.Items {
				var start, end time.Time
				if event.Start.DateTime == "" {
					start, err = time.Parse("2006-01-02", event.Start.Date)
				} else {
					start, err = time.Parse(time.RFC3339, event.Start.DateTime)
				}
				if err != nil {
					log.Printf("error parsing start time: %v", err)
					continue
				}

				if event.End.DateTime == "" {
					end, err = time.Parse("2006-01-02", event.End.Date)
				} else {
					end, err = time.Parse(time.RFC3339, event.End.DateTime)
				}
				if err != nil {
					log.Printf("error parsing end time: %v", err)
					continue
				}

				duration := end.Sub(start)
				eventsColorTime[event.ColorId] = append(eventsColorTime[event.ColorId], struct {
					Start    *calendar.EventDateTime
					End      *calendar.EventDateTime
					Duration time.Duration
				}{
					Start:    event.Start,
					End:      event.End,
					Duration: duration,
				})
			}

			pageToken = events.NextPageToken
			if pageToken == "" {
				break
			}
		}

		eventsColorSummaryTimeInHours := make(map[string]float64)
		var resultText string
		for colorID, timeRanges := range eventsColorTime {
			color := connection.ColorNames[colorID]
			if color == "" {
				color = "unknown"
			}
			var totalDuration time.Duration
			for _, timeRange := range timeRanges {
				totalDuration += timeRange.Duration
			}
			totalHours := totalDuration.Hours()
			eventsColorSummaryTimeInHours[color] = totalHours
			resultText += fmt.Sprintf("Total time for %s: %s\n", color, connection.FormatHours(totalHours))
		}

		resultLabel.SetText(resultText)

		statistics.Statistics(eventsColorTime, dateFrom.Text, dateTo.Text)

		viewChatButton.Enable()

		dialog.ShowInformation("Success", "Chart generated as bar.html", w)
	}

	viewChatButton.OnTapped = func() {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("start", statistics.GetDiagramPath())
		case "darwin":
			cmd = exec.Command("open", statistics.GetDiagramPath())
		default:
			cmd = exec.Command("xdg-open", statistics.GetDiagramPath())
		}
		if err := cmd.Run(); err != nil {
			dialog.ShowError(fmt.Errorf("error opening chart: %v", err), w)
			return
		}
	}

	leftPanel := container.NewVBox(
		widget.NewLabel("Google calendar parser"),
		form,
		submitButton,
		container.NewHBox(resultLabel, layout.NewSpacer(), viewChatButton),
	)

	w.SetContent(leftPanel)
	w.ShowAndRun()
}
