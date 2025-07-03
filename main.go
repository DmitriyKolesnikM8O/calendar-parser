package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

/*
6 - красный цвет
2 - зеленый цвет
_ - синий цвет
3 - фиолетовый цвет
4 - фламинго
5 - желтый
8 - серый
11 - ярко красный
7 - павлин
*/

var (
	tokenFile  = "token.json"
	colorNames = map[string]string{
		"6":  "red",
		"2":  "green",
		"":   "blue",
		"3":  "violet",
		"4":  "flamingo",
		"5":  "yellow",
		"8":  "grey",
		"11": "bright red",
		"7":  "bright blue",
	}

	timeStart, timeEnd string

	mode = flag.Bool("gui", false, "Running in GUI mode")
)

func getToken(config *oauth2.Config) (*oauth2.Token, error) {

	if token, err := tokenFromFile(tokenFile); err == nil {
		return token, nil
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Printf("Go to link and auth:\n%v\n", authURL)
	fmt.Println("write here auth code:")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("error when reading auth code: %v", err)
	}

	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("error when try to receive token: %v", err)
	}

	if err := saveToken(tokenFile, token); err != nil {
		log.Printf("ATTEMPT: can`t save token: %v", err)
	}

	return token, nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(token)
}

func formatHours(hours float64) string {
	h := int(hours)
	m := int((hours - float64(h)) * 60)
	return fmt.Sprintf("%d h. %02d min.", h, m)
}

// TODO: it is bad to send struct as parameter???????
// TODO: too much function, I think need separated functions
// TODO:duration in minutes -> duration in hours?????
func statistics(eventsColorTime map[string][]struct {
	Start    *calendar.EventDateTime
	End      *calendar.EventDateTime
	Duration time.Duration
}, timeStart string, timeEnd string) {

	eventsColorSummaryTimeInHours := make(map[string]float64)
	for colorID, timeRanges := range eventsColorTime {
		color := colorNames[colorID]
		if color == "" {
			log.Printf("error when parsing color in statistics func\n")
		}

		var totalDuration time.Duration

		for _, timeRange := range timeRanges {
			totalDuration += timeRange.Duration
		}
		totalHours := totalDuration.Hours()
		fullHours := int(totalHours)
		minutes := int((totalHours - float64(fullHours)) * 60)

		eventsColorSummaryTimeInHours[color] = totalHours
		fmt.Printf("Total time for color %s on period from %s to %s - %d h. %d m.\n",
			color, timeStart, timeEnd, fullHours, minutes)
	}

	keys := make([]string, 0, len(eventsColorSummaryTimeInHours))
	for key := range eventsColorSummaryTimeInHours {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return eventsColorSummaryTimeInHours[keys[i]] < eventsColorSummaryTimeInHours[keys[j]]
	})

	colorMap := map[string]string{
		"flamingo":    "#DE8157",
		"violet":      "#9B2AC9",
		"yellow":      "#EFD10F",
		"blue":        "#1634DB",
		"green":       "#17D427",
		"bright red":  "#FF0000",
		"red":         "#E3450B",
		"grey":        "#7D7877",
		"bright blue": "#031D9C",
	}

	description := map[string]string{
		"yellow":      "time instead",
		"red":         "important events",
		"grey":        "trains",
		"blue":        "useful activities",
		"bright red":  "another option for important events",
		"green":       "sleep",
		"violet":      "useless activities",
		"flamingo":    "cooking and eating",
		"bright blue": "anouther useful activities",
	}

	var values []opts.BarData
	for _, colorId := range keys {
		values = append(values, opts.BarData{
			Value: eventsColorSummaryTimeInHours[colorId],
			Name:  fmt.Sprintf("%s - (%s)", description[colorId], formatHours(eventsColorSummaryTimeInHours[colorId])),
			ItemStyle: &opts.ItemStyle{
				Color: colorMap[colorId],
			},
		})
	}

	bar := charts.NewBar()

	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: "Diagramm of spended time from " + timeStart + " to " + timeEnd,
	}), charts.WithAnimation(true), charts.WithXAxisOpts(opts.XAxis{
		Name: "color",
	}), charts.WithYAxisOpts(opts.YAxis{
		Name: "Duration in hours",
	}), charts.WithLegendOpts(opts.Legend{
		Data: []string{
			"Yellow - time instead",
			"Red - important events",
			"Grey - trains",
			"Blue - useful activities",
			"Bright Red - another option for important events",
			"Green - sleep",
			"Violet - useless activities",
			"Flamingo - cooking and eating",
		},
		Top:   "10%",
		Right: "80%",
	}))

	bar.SetXAxis(keys).AddSeries("Time", values)

	f, _ := os.Create("bar.html")
	bar.Render(f)

	// for _, key := range keys {
	// 	fmt.Printf("%s %d\n", key, eventsColorSummaryTimeInMinutes[key])
	// }

}

// TODO: too much main function
// TODO: long celectors and other garbage in code
func main() {

	flag.Parse()
	if *mode || len(os.Args) == 1 {
		RunGUI()
	}

	ctx := context.Background()

	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("error when reading info from credentials: %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("error when load configuration: %v", err)
	}

	token, err := getToken(config)
	if err != nil {
		log.Fatalf("error when try to receive token in main: %v", err)
	}

	client := config.Client(ctx, token)
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("error when creating a service: %v", err)
	}

	calendars, err := srv.CalendarList.List().Do()
	if err != nil {
		log.Fatalf("error when receive a lists of calendars: %v", err)
	}

	listsOfCalendars := make(map[int]string)
	fmt.Println("ivailable calendars:")
	for i, item := range calendars.Items {
		fmt.Printf("%d - %s (%s)\n", i, item.Summary, item.Id)
		listsOfCalendars[i] = item.Id
	}

	var numberOfCalendar int
	fmt.Printf("Which calendar could you want to use: ")
	_, err = fmt.Scanf("%d", &numberOfCalendar)
	if err != nil {
		log.Fatalf("error when parsing number of calendar: %v", err)
	}

	fmt.Printf("Write time start for parsing (format: YYYY-MM-DD): ")
	_, err = fmt.Scanf("%s", &timeStart)
	if err != nil {
		log.Fatalf("error when parsing time of events: %v", err)
	}
	fmt.Printf("Write time end for parsing (format: YYYY-MM-DD): ")
	_, err = fmt.Scanf("%s", &timeEnd)
	if err != nil {
		log.Fatalf("error when parsing time end of events: %v", err)
	}

	var AllEvents []*calendar.Events
	pageToken := ""
	for {
		events, err := srv.Events.List(listsOfCalendars[numberOfCalendar]).TimeMin(timeStart + "T10:00:00+03:00").TimeMax(timeEnd + "T23:59:59+03:00").PageToken(pageToken).Do()
		if err != nil {
			log.Fatalf("error when receive events: %v", err)
		}
		AllEvents = append(AllEvents, events)
		pageToken = events.NextPageToken
		// fmt.Println(pageToken)
		if pageToken == "" {
			break
		}

	}

	eventsColorTime := make(map[string][]struct {
		Start, End *calendar.EventDateTime
		Duration   time.Duration
	})
	for _, events := range AllEvents {
		for _, event := range events.Items {
			// fmt.Printf("ColorId: %s Creator: %s, Start: %s, End: %s, Summary: %s\n", event.ColorId, event.Creator, event.Start, event.End,
			// 	event.Summary)

			var start, end time.Time
			if event.Start.DateTime == "" {
				start, err = time.Parse("2006-01-02", event.Start.Date)
				if err != nil {
					log.Fatalf("error parsing start time: %v", err)
				}
			} else {
				start, err = time.Parse(time.RFC3339, event.Start.DateTime)
				if err != nil {
					log.Fatalf("error parsing start time: %v", err)
				}
			}

			if event.Start.DateTime == "" {
				end, err = time.Parse("2006-01-02", event.End.Date)
				if err != nil {
					log.Fatalf("error parsing end time: %v", err)
				}
			} else {
				end, err = time.Parse(time.RFC3339, event.End.DateTime)
				if err != nil {
					log.Fatalf("error parsing end time: %v", err)
				}
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
	}

	// for i, k := range eventsColorTime {
	// 	for _, v := range k {
	// 		fmt.Printf("ColorId: %s, event: %s\n", i, v)
	// 	}
	// }

	statistics(eventsColorTime, timeStart, timeEnd)
}
