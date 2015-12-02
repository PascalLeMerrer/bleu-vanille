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

	// DatabasePassword is the password of the Postgresql account to be used
	DatabasePassword = getEnv("DatabasePassword")

	// DatabasePort is the port on which Postgresql listens to
	DatabasePort = getNumericEnv("DatabasePort")

	// DatabaseUser is the name of the Postgresql account to be used
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
)

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

func getBooleanEnvWithDefault(name string, defaultvalue bool) bool {
	value := os.Getenv(name)

	if value == "" {
		return defaultvalue;
	}
	
	result, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("The environment variable %v, must be a boolean. Please fix it, then relaunch the server.\n", name)
		result = false
	}
	return result
}