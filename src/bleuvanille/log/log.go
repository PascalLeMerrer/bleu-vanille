package log

import (
	"time"
	golog "log"
	"github.com/labstack/echo"
	"github.com/twinj/uuid"
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
	correlationId := context.Get("correlationId")
	nowstr := time.Now().Format(time.UnixDate)
	
	if id, ok := correlationId.(string) ; ok {
		golog.Println(level + " - " + nowstr + " - " + id + " - " + message);	
	} else {
		golog.Println(nowstr + " - no correlationId - " + message);	
	}
}
