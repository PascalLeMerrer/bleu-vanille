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
func Error(context *echo.Context, message string) {
	_log("ERROR", context, message)
}

// Debug writes an debug message
func Debug(context *echo.Context, message string) {
	_log("DEBUG", context, message)
}

// Info writes any other message
func Info(context *echo.Context, message string) {
	_log("INFO", context, message)
}

// Info writes any other message
func Fatal(v ...interface{}) {
	golog.Fatal(v)
}

func Printf(format string, v ...interface{}) {
	golog.Printf(format, v)
}

// Middleware retrieves the session for an authenticated user
// It also deletes session for expired token
func Middleware() echo.HandlerFunc {
	return func(context *echo.Context) error {
		correlationId := uuid.NewV4().String()

		context.Set("correlationId", correlationId)
		return nil
	}
}

func _log(level string, context *echo.Context, message string) {
	var correlationId string
	nowstr := time.Now().Format(time.UnixDate)

	if context == nil {
		correlationId = "no correlationId"
	} else {
		if correlationIdAux, ok := context.Get("correlationId").(string); !ok {
			golog.Println(nowstr + " - no correlationId - " + message)
			return
		} else {
			correlationId = correlationIdAux
		}
	}
	golog.Println(level + " - " + nowstr + " - " + correlationId + " - " + message)
}
