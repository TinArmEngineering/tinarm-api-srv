package openapi

import (
	"context"
	"log"
	"os"
	"reflect"

	"time"

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
	db.dbClient, err = mongo.Connect(db.dbCtx, options.Client().ApplyURI(dbConnectionString()).SetRegistry(reg))
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

	collection := db.getDbClient().Database(dbName()).Collection(col)

	defer db.close()

	result, err := collection.InsertOne(db.dbCtx, doc)
	if err != nil {
		return "", err
	}

	mongoId, err := result.InsertedID.(primitive.ObjectID).MarshalText()
	return string(mongoId[:]), err
}

func (db Database) Query(col string, query, field interface{}) (result *mongo.Cursor, err error) {

	collection := db.getDbClient().Database(dbName()).Collection(col)

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

	var result = db.getDbClient().Database(dbName()).Collection(col).FindOne(
		db.dbCtx,
		bson.M{"_id": objectId})

	defer db.close()

	return result
}

func getEnvironment(name string, defaultOnEmpty string) string {

	setting := os.Getenv(name)
	if setting == "" {
		setting = defaultOnEmpty
	}
	return setting
}

func GetMaterialById(id string) (Material, error) {

	result := Database{}.QueryById("materials", id)

	material := Material{}
	err := result.Decode(&material)

	var materialId interface{} = id
	material.Id = &materialId

	return material, err
}
