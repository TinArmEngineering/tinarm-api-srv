// dbo/setup.go

package dbo

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// Dev config
	DB_HOST    = "127.0.0.1:3306"
	DB_USER    = "root"
	DB_PASS    = "tinarm"
	DB_NAME    = "tinarm_server_dev"
	QUEUE_HOST = "localhost:5672"
)

var DB *gorm.DB

func ConnectDatabase() {

	dbHost := getEnvironment("GOSERVER_DB_HOST", DB_HOST)
	dbUser := getEnvironment("GOSERVER_DB_USER", DB_USER)
	dbPass := getEnvironment("GOSERVER_DB_PASS", DB_PASS)
	dbName := getEnvironment("GOSERVER_DB_NAME", DB_NAME)

	dbConnection := dbUser + ":" + dbPass + "@tcp(" + dbHost + ")/"
	dbConnectionString := dbConnection + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	log.Printf("Using connection string: " + dbConnectionString)

	ensureDbExists(dbConnection, dbName)

	db, err := gorm.Open(mysql.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Create/Migrate schema
	err = db.AutoMigrate(&Job{})

	if err != nil {
		return
	}

	DB = db
}

// Post to RabbitMQ
func Enqueue(body string) {

	log.Printf("Queueing: " + body)

	queueHost := getEnvironment("GOSERVER_QUEUE_HOST", QUEUE_HOST)
	conn, err := amqp.Dial("amqp://guest:guest@" + queueHost + "/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rectangle_mesh_queue", // name
		true,                   // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	failOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}

func getEnvironment(name string, defaultOnEmpty string) string {

	setting := os.Getenv(name)
	if setting == "" {
		setting = defaultOnEmpty
	}
	return setting
}

func ensureDbExists(dbConnection string, dbName string) {

	db, err := sql.Open("mysql", dbConnection)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		panic(err)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
