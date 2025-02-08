package ctx

import (
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type (
	ContextLogger    struct{}
	ContextRequestId struct{}
)

var (
	CustomLogger = &log.Logger{
		Out: os.Stdout,
		Formatter: &log.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
			PrettyPrint:     true,
		},
		//Formatter:    &log.TextFormatter{},
		Hooks:        make(log.LevelHooks),
		Level:        log.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
)
