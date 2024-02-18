package logs

import(
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

var logLevelMap = map[string]zerolog.Level{
	"trace": zerolog.TraceLevel,
	"debug": zerolog.DebugLevel,	
	"info": zerolog.InfoLevel,
	"warn": zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
	"fatal": zerolog.FatalLevel,
	"panic": zerolog.PanicLevel,
} 

func Init(logLevel string, dev bool) {
	// Define the logger
	if dev {
		// Initiate the logger.
		cw := zerolog.ConsoleWriter{
			Out: os.Stderr,
			NoColor: false,
			TimeFormat: time.TimeOnly, // "15:04:05"
		}
		Logger = zerolog.New(cw).Level(logLevelMap[logLevel]).With().Timestamp().Logger()
		
	} else {
		Logger = zerolog.New(os.Stderr).Level(logLevelMap[logLevel]).With().Timestamp().Caller().Logger()
	}
}
