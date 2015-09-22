package log

import "github.com/labstack/echo"

// log package wraps the default GO log package
// It uses the echo context to add a correlation Id in logs,
// so we can trace the processing of a given request across logs
// It also format all logs in the same way

// Error writes an error message - Don't forget to provide
// contextual information to make easy the debug process
func Error(context *echo.Context, message string) {

}

// Debug writes an debug message
func Debug(context *echo.Context, message string) {

}

// Info writes any other message
func Info(context *echo.Context, message string) {

}
