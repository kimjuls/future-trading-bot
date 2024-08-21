package utils

import (
	"os"

	"github.com/sirupsen/logrus"
	logger "github.com/sirupsen/logrus"
)

type Log struct {
	Msg    interface{}
	Fields logger.Fields
	Level  logger.Level
}

var Logger = make(chan Log)

func init() {
	go func() {
		logger.SetFormatter(&logger.TextFormatter{
			ForceColors: true, FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05",
		})
		logger.SetOutput(os.Stdout)
		logger.SetLevel(logger.TraceLevel)

		for log := range Logger {
			logger.WithFields(log.Fields).Log(log.Level, log.Msg)
		}
	}()
}

func Info(msg string) {
	log(msg, logrus.InfoLevel)
}

func Warn(msg string) {
	log(msg, logrus.WarnLevel)
}

func Fatal(msg string) {
	log(msg, logrus.FatalLevel)
}

func Panic(msg string) {
	log(msg, logrus.PanicLevel)
}

func log(msg string, level logrus.Level) {
	Logger <- Log{Msg: msg, Level: level}
}
