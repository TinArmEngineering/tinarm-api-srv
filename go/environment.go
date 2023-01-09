package openapi

const (
	// Dev config
	MONGODB_HOST = "localhost:27017"
	MONGODB_USER = "root"
	MONGODB_PASS = "example"
	MONGODB_NAME = "sno2"
	QUEUE_HOST   = "localhost:5672"
)

func qHost() string {
	return getEnvironment("GOSERVER_QUEUE_HOST", QUEUE_HOST)
}

func dbHost() string {
	return getEnvironment("GOSERVER_DB_HOST", MONGODB_HOST)
}

func dbUser() string {
	return getEnvironment("GOSERVER_DB_USER", MONGODB_USER)
}

func dbPass() string {
	return getEnvironment("GOSERVER_DB_PASS", MONGODB_PASS)
}

func dbName() string {
	return getEnvironment("GOSERVER_DB_NAME", MONGODB_NAME)
}

func queueConnectionString() string {
	return "amqp://guest:guest@" + qHost() + "/"
}

func dbConnectionString() string {
	return "mongodb://" + dbUser() +
		":" + dbPass() +
		"@" + dbHost()
}
