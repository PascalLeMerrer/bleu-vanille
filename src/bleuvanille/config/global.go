package config

import (
	"log"
	"os"
	"strconv"
)

const (
	// AdminEmail is the login of the administrator account
	AdminEmail = "admin@bleuvanille.com"

	// TestUserEmail is the login of the basic test account used by the test scripts
	TestUserEmail = "testuser@bleuvanille.com"

	// NoReplyAddress Email address to be used as sender of email for which no response is expected
	NoReplyAddress = "noreply@bleuvanille.com"

	// ResetTokenDuration is the Time To Live for a passord reset token, in minutes
	ResetTokenDuration = 5

	// SessionDuration is the duration of the session, in hours
	SessionDuration = 1

	// SessionTokenDuration defines the time to live of the token, in hours
	SessionTokenDuration = 96
)

var (
	// DatabaseName is the name of the database used to store the service data
	DatabaseName = getEnv("DatabaseName")

	// DatabaseRootPassword is the password of the ArangoDB root account
	DatabaseRootPassword = getEnv("DatabaseRootPassword")

	// DatabasePassword is the password of the ArangoDB account to be used
	DatabasePassword = getEnv("DatabasePassword")

	// DatabaseProtocol must be http or https
	DatabaseProtocol = getEnvWithDefault("DatabaseProtocol", "http")

	// DatabaseHost is the name or IP of the ArangoDB server to be used
	DatabaseHost = getEnvWithDefault("DatabaseHost", "localhost")

	// DatabasePort is the port on which ArangoDB listens to
	DatabasePort = getNumericEnv("DatabasePort")

	// DatabaseUser is the name of the ArangoDB account to be used
	DatabaseUser = getEnv("DatabaseUser")

	// HostName is the name of the server.
	HostName = getEnv("BleuVanilleName")

	// HostPort is the port on which the server listens to.
	HostPort = getNumericEnv("BleuVanillePort")

	// SMTPPort is the port Number of the SMTP server
	SMTPPort = getEnv("SMTPPort")

	// SMTPPassword is the password of the account used for authentication on the SMTP SMTPServer
	SMTPPassword = getEnv("SMTPPassword")

	// SMTPServer The SMTP server to be used to send emails
	SMTPServer = getEnv("SMTPServer")

	// SMTPUser is the email address of the account used for authentication on the SMTP SMTPServer
	SMTPUser = getEnv("SMTPUser")

	// ProductionMode indicates if the instance is running in production.
	ProductionMode = getBooleanEnvWithDefault("ProductionMode", false)

	// Debug mode displays useful information during dev
	Debug = getEnv("Debug")
)

// ServerDebug when server logs should be in verbose mode
func ServerDebug() bool {
	return Debug == "server" || Debug == "all"
}

// DbDebug returns true when database request have to be displayed in console
func DbDebug() bool {
	return Debug == "db" || Debug == "database" || Debug == "all"
}

func getNumericEnv(name string) int {
	value := os.Getenv(name)

	if value == "" {
		log.Printf("Please define the environment variable %v, then relaunch the server.\n", name)
		return -1
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("The environment variable %v, must be an integer. Please fix it, then relaunch the server.\n", name)

	}
	return result
}

func getEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Printf("Please define the environment variable %v, then relaunch the server.\n", name)
	}
	return value
}

func getEnvWithDefault(name string, defaultValue string) string {
	value := os.Getenv(name)

	if value != "" {
		return value
	}

	return defaultValue
}

func getBooleanEnv(name string) bool {
	value := os.Getenv(name)

	if value == "" {
		log.Printf("Please define the environment variable %v, then relaunch the server.\n", name)
		return false
	}

	result, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("The environment variable %v, must be a boolean. Please fix it, then relaunch the server.\n", name)
		result = false
	}
	return result
}

func getBooleanEnvWithDefault(name string, defaultValue bool) bool {
	value := os.Getenv(name)

	if value == "" {
		return defaultValue
	}

	result, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("The environment variable %v, must be a boolean. Please fix it, then relaunch the server.\n", name)
		result = false
	}
	return result
}
