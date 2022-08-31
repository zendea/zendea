package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
)

const timeFormat = "2006-01-02 15:04:05"

var log zerolog.Logger

func init() {
	zerolog.CallerSkipFrameCount = 3
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: timeFormat}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf(" | %s", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf(" | %s", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf(" %s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s ", i))
	}
	output.FormatCaller = func(i interface{}) string {
		var c string
		if cc, ok := i.(string); ok {
			c = cc
		}
		if len(c) > 0 {
			cwd, err := os.Getwd()
			if err == nil {
				c = strings.TrimPrefix(c, cwd)
				c = strings.TrimPrefix(c, "/")
			}
		}
		return "| " + c
	}
	log = zerolog.New(output).With().Timestamp().Logger()

}

//Debug : Level 0
func Debug(format string, v ...interface{}) {
	log.Debug().Caller().Msgf(format, v...)
}

//Info : Level 1
func Info(format string, v ...interface{}) {
	log.Info().Caller().Msgf(format, v...)
}

//Warn : Level 2
func Warn(format string, v ...interface{}) {
	log.Warn().Caller().Msgf(format, v...)
}

//Error : Level 3
func Error(format string, v ...interface{}) {
	log.Error().Caller().Msgf(format, v...)
}

//Fatal : Level 4
func Fatal(format string, v ...interface{}) {
	log.Fatal().Caller().Msgf(format, v...)
}
