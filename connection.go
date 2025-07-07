package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
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

func initializeCalendarService(w fyne.Window) (*calendar.Service, *calendar.CalendarList, error) {
	ctx := context.Background()

	b, err := os.ReadFile("credentials.json")
	if err != nil {
		if w != nil {
			dialog.ShowError(fmt.Errorf("error reading credentials: %v", err), w)
		} else {
			log.Fatalf("error when reading info from credentials: %v", err)
		}
		return nil, nil, err
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		if w != nil {
			dialog.ShowError(fmt.Errorf("error loading configuration: %v", err), w)
		} else {
			log.Fatalf("error when load configuration: %v", err)
		}
		return nil, nil, err
	}

	token, err := getToken(config)
	if err != nil {
		if w != nil {
			dialog.ShowError(fmt.Errorf("error retrieving token: %v", err), w)
		} else {
			log.Fatalf("error when try to receive token in main: %v", err)
		}
		return nil, nil, err
	}

	client := config.Client(ctx, token)
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		if w != nil {
			dialog.ShowError(fmt.Errorf("error creating calendar service: %v", err), w)
		} else {
			log.Fatalf("error when creating a service: %v", err)
		}
		return nil, nil, err
	}

	calendars, err := srv.CalendarList.List().Do()
	if err != nil {
		if w != nil {
			dialog.ShowError(fmt.Errorf("error retrieving calendars: %v", err), w)
		} else {
			log.Fatalf("error when receive a lists of calendars: %v", err)
		}
		return nil, nil, err
	}

	return srv, calendars, nil
}
