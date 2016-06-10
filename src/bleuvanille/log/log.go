package log

import (
	"github.com/labstack/echo"
	"github.com/twinj/uuid"
	golog "log"
	"time"
)

// log package wraps the default GO log package
// It uses the echo context to add a correlation Id in logs,
// so we can trace the processing of a given request across logs
// It also format all logs in the same way

// Error writes an error message - Don't forget to provide
// contextual information to make easy the debug process
func Error(context echo.Context, message string) {
	_log("ERROR", context, message)
}

// Debug writes a debug message
func Debug(context echo.Context, message string) {
	_log("DEBUG", context, message)
}

// Info writes any other message
func Info(context echo.Context, message string) {
	_log("INFO", context, message)
}

// Fatal writes a crash message
func Fatal(v ...interface{}) {
	golog.Fatal(v)
}

// Printf format and writes a message to standard output
func Printf(format string, v ...interface{}) {
	golog.Printf(format, v)
}

// Middleware retrieves the session for an authenticated user
// It also deletes session for expired token
func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		correlationID := uuid.NewV4().String()

		context.Set("correlationID", correlationID)
		return next(context)

	}
}

func _log(level string, context echo.Context, message string) {
	var correlationID string
	nowstr := time.Now().Format(time.UnixDate)

	if context == nil {
		correlationID = "no correlationID"
	} else {
		correlationIDAux, ok := context.Get("correlationID").(string)
		if !ok {
			golog.Println(nowstr + " - no correlationID - " + message)
			return
		}
		correlationID = correlationIDAux
	}
	golog.Println(level + " - " + nowstr + " - " + correlationID + " - " + message)
}
