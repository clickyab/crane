package version

import (
	"context"
	"strconv"
	"time"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/initializer"

	"github.com/sirupsen/logrus"
)

// Data is the application version in detail
type Data struct {
	Hash      string    `json:"hash"`
	Short     string    `json:"short_hash"`
	Date      time.Time `json:"commit_date"`
	Count     int64     `json:"build_number"`
	BuildDate time.Time `json:"build_date"`
}

// GetVersion return the application version in detail
func GetVersion() Data {
	c := Data{}
	c.Count, _ = strconv.ParseInt(count, 10, 64)
	c.Date, _ = time.Parse(time.RFC1123Z, date)
	c.Hash = hash
	c.Short = short
	c.BuildDate, _ = time.Parse("01-02-06-15-04-05", build)

	return c
}

// LogVersion return an logrus entry with version information attached
func LogVersion() *logrus.Entry {
	ver := GetVersion()
	return logrus.WithFields(
		logrus.Fields{
			"Commit hash":       ver.Hash,
			"Commit short hash": ver.Short,
			"Commit date":       ver.Date.Format(time.RFC3339),
			"Build date":        ver.BuildDate.Format(time.RFC3339),
		},
	)
}

type show struct {
}

func (show) Initialize(ctx context.Context) {
	done := ctx.Done()
	assert.NotNil("[BUG] context is not cancelable")

	LogVersion().Debug("Start")
	go func() {
		<-done
		LogVersion().Debug("Done")
	}()
}

func init() {
	initializer.Register(&show{}, -10)
}
