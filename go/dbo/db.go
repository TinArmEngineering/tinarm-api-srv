package dbo

import (
	"context"
	"log"
	"os"

	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// Dev config
	MONGODB_HOST = "localhost:27017"
	MONGODB_USER = "root"
	MONGODB_PASS = "example"
	MONGODB_NAME = "sno2"
)

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

func dbConnectionString() string {
	return "mongodb://" + dbUser() +
		":" + dbPass() +
		"@" + dbHost()
}

// This is a user defined method to close resources.
// This method closes mongoDB connection and cancel context.
func close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {

	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

// This is a user defined method that returns mongo.Client,
// context.Context, context.CancelFunc and error.
// mongo.Client will be used for further database operation.
// context.Context will be used set deadlines for process.
// context.CancelFunc will be used to cancel context and
// resource associated with it.
func connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

// insertOne is a user defined method, used to insert
// documents into collection returns result of InsertOne
// and error if any.
func insertOne(
	client *mongo.Client,
	ctx context.Context,
	dataBase,
	col string,
	doc interface{}) (*mongo.InsertOneResult, error) {

	// select database and collection ith Client.Database method
	// and Database.Collection method
	collection := client.Database(dataBase).Collection(col)

	// InsertOne accept two argument of type Context
	// and of empty interface
	result, err := collection.InsertOne(ctx, doc)
	return result, err
}

// insertMany is a user defined method, used to insert
// documents into collection returns result of
// InsertMany and error if any.
func insertMany(
	client *mongo.Client,
	ctx context.Context,
	dataBase,
	col string,
	docs []interface{}) (*mongo.InsertManyResult, error) {

	// select database and collection ith Client.Database
	// method and Database.Collection method
	collection := client.Database(dataBase).Collection(col)

	// InsertMany accept two argument of type Context
	// and of empty interface
	result, err := collection.InsertMany(ctx, docs)
	return result, err
}

func getEnvironment(name string, defaultOnEmpty string) string {

	setting := os.Getenv(name)
	if setting == "" {
		setting = defaultOnEmpty
	}
	return setting
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func InsertJob(job interface{}) (string, error) {
	return insertOneThing("jobs", job)
}

func InsertMaterial(material interface{}) (string, error) {
	return insertOneThing("material", material)
}

func insertOneThing(thingName string, thing interface{}) (string, error) {

	client, ctx, cancel, err := connect(dbConnectionString())
	if err != nil {
		return "", err
	}

	insertOneResult, err := insertOne(client, ctx, dbName(), thingName, thing)
	if err != nil {
		return "", err
	}

	defer close(client, ctx, cancel)

	return insertOneResult.InsertedID.(primitive.ObjectID).String(), err
}
