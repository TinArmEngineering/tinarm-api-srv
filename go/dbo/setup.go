// dbo/setup.go

package dbo

import (
	"database/sql"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	// Dev config
	DB_HOST = "127.0.0.1:3306"
	DB_USER = "root"
	DB_PASS = "tinarm"
	DB_NAME = "tinarm_server_dev"
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
