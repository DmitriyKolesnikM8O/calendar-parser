package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var (
	tokenFile = "token.json"
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

func main() {
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

	lists_of_calendars := make(map[int]string)
	fmt.Println("ivailable calendars:")
	for i, item := range calendars.Items {
		fmt.Printf("%d - %s (%s)\n", i, item.Summary, item.Id)
		lists_of_calendars[i] = item.Id
	}

	var number_of_calendar int
	fmt.Printf("Which calendar could you want to use: ")
	_, err = fmt.Scanf("%d", &number_of_calendar)
	if err != nil {
		log.Fatalf("error when parsing number of calendar: %v", err)
	}

	//TODO: normal parsing for time zones
	events, err := srv.Events.List(lists_of_calendars[number_of_calendar]).TimeMin("2025-05-24T10:00:00+03:00").TimeMax("2025-05-31T10:00:00+03:00").Do()
	if err != nil {
		log.Fatalf("error when receive events: %v", err)
	}
	for _, event := range events.Items {
		fmt.Printf("ColorId: %s Creator: %s, Start: %s, End: %s, Summary: %s\n", event.ColorId, event.Creator, event.Start, event.End,
			event.Summary)
	}
}
