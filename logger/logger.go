package logger

import (
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func SetupLog(level string) {
	logrus.SetLevel(logLevel(level))
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(TimestampFormatter())
}

func logLevel(level string) logrus.Level {
	switch strings.ToUpper(level) {
	case "INFO":
		return logrus.InfoLevel
	case "DEBUG":
		return logrus.DebugLevel
	case "WARN":
		return logrus.WarnLevel
	case "ERROR":
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}

func randString(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

const LogRefKey = "reference"

//LogEntryWithRef returns a logrus Entry with a random unique value for requestId field
func LogEntryWithRef() *logrus.Entry {
	return logrus.WithField(LogRefKey, randString(10))
}

//TimestampFormatter returns a custom logrus Formatter with timestamp enabled
func TimestampFormatter() logrus.Formatter {
	formatter := new(logrus.JSONFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.PrettyPrint = false

	return formatter
}
