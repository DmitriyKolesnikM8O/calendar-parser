package charts

import "time"

type BarChartRenderer struct{}

func (b *BarChartRenderer) Render(stats map[string]time.Duration, outputFile string) error {
	return nil
}
