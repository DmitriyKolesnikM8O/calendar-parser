package statistics

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"time"

	"github.com/DmitriyKolesnikM8O/calendar-parser/internal/connection"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"google.golang.org/api/calendar/v3"
)

var (
	OutputPath string
)

func GetDiagramPath() string {

	_, filename, _, _ := runtime.Caller(0)

	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	diagramDir := filepath.Join(rootDir, "diagramm")
	os.MkdirAll(diagramDir, 0755)

	return diagramDir
}

func Statistics(eventsColorTime map[string][]struct {
	Start    *calendar.EventDateTime
	End      *calendar.EventDateTime
	Duration time.Duration
}, timeStart string, timeEnd string) {

	eventsColorSummaryTimeInHours := make(map[string]float64)
	for colorID, timeRanges := range eventsColorTime {
		color := connection.ColorNames[colorID]
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
			Name:  fmt.Sprintf("%s - (%s)", description[colorId], connection.FormatHours(eventsColorSummaryTimeInHours[colorId])),
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

	OutputPath = filepath.Join(GetDiagramPath(), timeStart+" - "+timeEnd+".html")
	f, err := os.Create(OutputPath)
	if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
	}
	bar.Render(f)

}
