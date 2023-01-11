package services

import (
	"context"
	"log"
	"reflect"

	"time"

	env "github.com/tinarmengineering/sno2-srv-go/go/environment"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	dbClient     *mongo.Client
	dbCtx        context.Context
	dbCancelFunc context.CancelFunc
}

func (db *Database) getDbClient() *mongo.Client {
	if db.dbClient == nil {
		db.connect()
	}
	return db.dbClient
}

func (db *Database) connect() (err error) {

	rb := bson.NewRegistryBuilder()
	rb.RegisterTypeMapEntry(bsontype.EmbeddedDocument, reflect.TypeOf(bson.M{}))
	reg := rb.Build()

	db.dbCtx, db.dbCancelFunc = context.WithTimeout(context.Background(), 30*time.Second)
	db.dbClient, err = mongo.Connect(db.dbCtx, options.Client().ApplyURI(env.DbConnectionString()).SetRegistry(reg))
	return err
}

func (db *Database) close() {

	defer db.dbCancelFunc()

	defer func() {
		if err := db.dbClient.Disconnect(db.dbCtx); err != nil {
			panic(err)
		}
	}()
}

func (db Database) InsertOne(col string, doc interface{}) (string, error) {

	collection := db.getDbClient().Database(env.DbName()).Collection(col)
	defer db.close()

	result, err := collection.InsertOne(db.dbCtx, doc)
	if err != nil {
		return "", err
	}

	mongoId, err := result.InsertedID.(primitive.ObjectID).MarshalText()
	return string(mongoId[:]), err
}

func (db Database) Query(col string, query, field interface{}) (result *mongo.Cursor, err error) {

	collection := db.getDbClient().Database(env.DbName()).Collection(col)
	defer db.close()

	result, err = collection.Find(
		db.dbCtx,
		query,
		options.Find().SetProjection(field))

	return
}

func (db Database) QueryById(col string, id string) *mongo.SingleResult {

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	var result = db.getDbClient().Database(env.DbName()).Collection(col).FindOne(db.dbCtx, bson.M{"_id": objectId})
	defer db.close()

	return result
}
