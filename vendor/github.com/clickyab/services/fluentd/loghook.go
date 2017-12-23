package fluentd

import (
	"context"

	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	"github.com/evalphobia/logrus_fluent"
	"github.com/sirupsen/logrus"
)

var (
	host       = config.RegisterString("services.fluentd.host", "fluentd.monitoring", "fluentd host")
	port       = config.RegisterInt64("services.fluentd.port", 24224, "fluentd port")
	active     = config.RegisterBoolean("services.fluentd.enable", false, "fluentd enable")
	defaultTag = config.RegisterString("services.fluentd.tag", "change.me", "fluentd default tag")
	allLevels  = config.RegisterBoolean("services.fluentd.all_levels", false, "send all logs, also info and debugs")
)

type hook struct {
}

func (hook) Initialize(ctx context.Context) {
	if !active.Bool() {
		return
	}
	hook, err := logrus_fluent.New(host.String(), port.Int())
	if err != nil {
		logrus.Errorf("fluentd logger failed, if this is in production check for the problem: %s", err)
		return
	}

	// set custom fire level
	l := []logrus.Level{
		logrus.PanicLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.WarnLevel,
	}
	if allLevels.Bool() {
		l = append(l, logrus.InfoLevel, logrus.DebugLevel)
	}
	hook.SetLevels(l)

	// set static tag
	hook.SetTag(defaultTag.String())
	// filter func
	// TODO : write more filter to handle more type (for clarification on logger side)
	hook.AddFilter("error", logrus_fluent.FilterError)
	logrus.AddHook(hook)

	go func() {
		<-ctx.Done()
		// somehow it can make a race, but there is no objection, its dying anyway
		tmp := hook.Fluent
		hook.Fluent = nil
		err := tmp.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()
}

func init() {
	initializer.Register(&hook{}, 0)
}
