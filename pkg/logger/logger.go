package logger

import (
	"flag"
	"os"

	"github.com/google/logger"
)

const logPath = "/var/log/notificaation-service.log"

var verbose = flag.Bool("verbose", false, "print info level logs to stdout")

func init() {
	flag.Parse()

	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	defer lf.Close()

	defer logger.Init("LoggerExample", *verbose, true, lf).Close()

}

func Info(v interface{}) {
	logger.Info(v)
}

func Error(err error) {
	logger.Error(err)
}
