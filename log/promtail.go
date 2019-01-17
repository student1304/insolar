package log

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/afiskon/promtail-client/promtail"
	"github.com/insolar/insolar/core"
)

type promtailAdapter struct {
	impl  promtail.Client
	level string
}

func newPromtailAdapter(agentUrl string, levelStr string) *promtailAdapter {
	level := promtail.DEBUG
	switch levelStr {
	case "Info":
		level = promtail.INFO
	case "Warn":
		level = promtail.WARN
	case "Error":
		level = promtail.ERROR
	default:
		levelStr = "Debug"
		log.Printf("Unknown log level '%s' - ignored\n", levelStr)
	}

	conf := promtail.ClientConfig{
		PushURL:            agentUrl,
		Labels:             "{job=\"insolard\"}",
		BatchWait:          5 * time.Second,
		BatchEntriesNumber: 10000,
		SendLevel:          level,
		PrintLevel:         level,
	}

	impl, err := promtail.NewClientProto(conf)
	if err != nil {
		log.Panicf("promtail.NewClientProto: %s", err)
	}

	return &promtailAdapter{impl: impl, level: levelStr}
}

func (pa *promtailAdapter) SetLevel(level string) error {
	// level is set in newPromtailAdapter
	return nil
}

func (pa *promtailAdapter) GetLevel() string {
	return pa.level
}

func (pa *promtailAdapter) Debug(args ...interface{}) {
	pa.impl.Debugf("%s", fmt.Sprint(args...))
}

func (pa *promtailAdapter) Debugln(args ...interface{}) {
	pa.impl.Debugf("%s", fmt.Sprintln(args...))
}

func (pa *promtailAdapter) Debugf(fmt string, args ...interface{}) {
	pa.impl.Debugf(fmt, args...)
}

func (pa *promtailAdapter) Info(args ...interface{}) {
	pa.impl.Infof("%s", fmt.Sprint(args...))
}

func (pa *promtailAdapter) Infoln(args ...interface{}) {
	pa.impl.Infof("%s", fmt.Sprintln(args...))
}

func (pa *promtailAdapter) Infof(fmt string, args ...interface{}) {
	pa.impl.Infof(fmt, args...)
}

func (pa *promtailAdapter) Warn(args ...interface{}) {
	pa.impl.Warnf("%s", fmt.Sprint(args...))
}

func (pa *promtailAdapter) Warnln(args ...interface{}) {
	pa.impl.Warnf("%s", fmt.Sprintln(args...))
}

func (pa *promtailAdapter) Warnf(fmt string, args ...interface{}) {
	pa.impl.Warnf(fmt, args...)
}

func (pa *promtailAdapter) Error(args ...interface{}) {
	pa.impl.Errorf("%s", fmt.Sprint(args...))
}

func (pa *promtailAdapter) Errorln(args ...interface{}) {
	pa.impl.Errorf("%s", fmt.Sprintln(args...))
}

func (pa *promtailAdapter) Errorf(fmt string, args ...interface{}) {
	pa.impl.Errorf(fmt, args...)
}

func (pa *promtailAdapter) Fatal(args ...interface{}) {
	pa.impl.Errorf("FATAL %s", fmt.Sprint(args...))
}

func (pa *promtailAdapter) Fatalln(args ...interface{}) {
	pa.impl.Errorf("FATAL %s", fmt.Sprintln(args...))
}

func (pa *promtailAdapter) Fatalf(fmt string, args ...interface{}) {
	pa.impl.Errorf("FATAL "+fmt, args...)
}

func (pa *promtailAdapter) Panic(args ...interface{}) {
	msg := fmt.Sprint(args...)
	pa.impl.Errorf("PANIC %s", msg)
	panic(msg)
}

func (pa *promtailAdapter) Panicln(args ...interface{}) {
	msg := fmt.Sprintln(args...)
	pa.impl.Errorf("PANIC %s", msg)
	panic(msg)
}

func (pa *promtailAdapter) Panicf(format string, args ...interface{}) {
	pa.impl.Errorf("PANIC "+format, args...)
	panic(fmt.Sprintf(format, args...))
}

func (pa *promtailAdapter) SetOutput(w io.Writer) {
	// do nothing, at least yet
	// SetOutput make little sense for the Promtail adapter
}

func (pa *promtailAdapter) WithField(string, string) core.Logger {
	return pa
	//panic("implement me")
}
