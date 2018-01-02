package main

import (
	"os"

	"sync"

	"context"

	"flag"

	"strconv"
	"time"

	"clickyab.com/crane/commands"
	"clickyab.com/crane/models/cron"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	_ "github.com/clickyab/services/mysql/connection/mysql"
	"github.com/clickyab/services/safe"
	"github.com/sirupsen/logrus"
)

var (
	wg = &sync.WaitGroup{}
)

func main() {
	config.Initialize(commands.Organization, commands.AppName, commands.Prefix, commands.DefaultConfig())
	config.DumpConfig(os.Stdout)
	defer initializer.Initialize()()
	s := time.Now()
	logrus.Debugf("start cron %s ", commands.AppName)

	def, _ := strconv.Atoi(time.Now().Format("20060102"))
	fDate := flag.Int("d", def, "date for cron job")
	flag.Parse()

	// if d=-1 the cron runs for yesterday
	if *fDate == -1 {
		yesterday, _ := strconv.Atoi(time.Now().AddDate(0, 0, -1).Format("20060102"))
		*fDate = yesterday
	}

	ctx := context.Background()
	wg.Add(4)

	safe.GoRoutine(ctx, func() {
		defer func() {
			logrus.Debugf("web impression worker done in %f seconds", time.Since(s).Seconds())
			wg.Done()
		}()
		assert.Nil(cron.WebImp(*fDate))
	})
	safe.GoRoutine(ctx, func() {
		defer func() {
			logrus.Debugf("web click worker done in %f seconds", time.Since(s).Seconds())
			wg.Done()
		}()
		assert.Nil(cron.WebClick(*fDate))
	})
	safe.GoRoutine(ctx, func() {
		defer func() {
			logrus.Debugf("app impression worker done in %f seconds", time.Since(s).Seconds())
			wg.Done()
		}()
		assert.Nil(cron.AppImp(*fDate))
	})
	safe.GoRoutine(ctx, func() {
		defer func() {
			logrus.Debugf("app click worker done in %f seconds", time.Since(s).Seconds())
			wg.Done()
		}()
		assert.Nil(cron.AppClick(*fDate))
	})

	wg.Wait()

	logrus.Debugf("cron %s finished and toke %f seconds for all tasks", commands.AppName, time.Since(s).Seconds())

}
