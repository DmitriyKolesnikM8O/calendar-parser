package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func RunCLI() {
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

	statistics(eventsColorTime, timeStart, timeEnd)
}
