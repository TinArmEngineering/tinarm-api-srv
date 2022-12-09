// dbo/setup.go

package dbo

import (
	"database/sql"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	// Dev config
	DB_HOST = "127.0.0.1:3306"
	DB_USER = "root"
	DB_PASS = "tinarm"
	DB_NAME = "tinarm_server"
)

var DB *gorm.DB

func ConnectDatabase() {

	dbHost := os.Getenv("GOSERVER_DB_HOST")
	if dbHost == "" {
		dbHost = DB_HOST
	}

	dbUser := os.Getenv("GOSERVER_DB_USER")
	if dbUser == "" {
		dbUser = DB_USER
	}

	dbPass := os.Getenv("GOSERVER_DB_PASS")
	if dbPass == "" {
		dbPass = DB_USER
	}

	dbConnection := DB_USER + ":" + DB_PASS + "@tcp(" + DB_HOST + ")/"
	dbConnectionString := dbConnection + DB_NAME + "?charset=utf8mb4&parseTime=True&loc=Local"

	ensureDbExists(dbConnection, DB_NAME)

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
