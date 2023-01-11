package environment

import "os"

const (
	// Dev config
	DEFAULT_DB_HOST            = "localhost:27017"
	DEFAULT_DB_USER            = "root"
	DEFAULT_DB_PASS            = "example"
	DEFAULT_DB_NAME            = "sno2"
	DEFAULT_QUEUE_HOST         = "localhost:5672"
	DEFAULT_UNITS_SERVICE_HOST = "localhost:5001"
)

func DbHost() string {
	return getEnvironment("GOSERVER_DB_HOST", DEFAULT_DB_HOST)
}

func DbUser() string {
	return getEnvironment("GOSERVER_DB_USER", DEFAULT_DB_USER)
}

func DbPass() string {
	return getEnvironment("GOSERVER_DB_PASS", DEFAULT_DB_PASS)
}

func DbName() string {
	return getEnvironment("GOSERVER_DB_NAME", DEFAULT_DB_NAME)
}

func QueueHost() string {
	return getEnvironment("GOSERVER_QUEUE_HOST", DEFAULT_QUEUE_HOST)
}

func UnitsServiceHost() string {
	return getEnvironment("GOSERVER_UNITS_SERVICE_HOST", DEFAULT_UNITS_SERVICE_HOST)
}

func QueueConnectionString() string {
	return "amqp://guest:guest@" + QueueHost() + "/"
}

func DbConnectionString() string {
	return "mongodb://" + DbUser() +
		":" + DbPass() +
		"@" + DbHost()
}

func getEnvironment(name string, defaultOnEmpty string) string {

	setting := os.Getenv(name)
	if setting == "" {
		setting = defaultOnEmpty
	}
	return setting
}
