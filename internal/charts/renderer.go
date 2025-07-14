package charts

import "time"

type ChartRenderer interface {
	Render(stats map[string]time.Duration, outputFile string) error
}
