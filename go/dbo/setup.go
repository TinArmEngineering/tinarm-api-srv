// dbo/setup.go

package dbo

import (
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	// Dev config
	DB_CONNECTION = "root:tinarm@tcp(127.0.0.1:40000)/"
	DB_NAME       = "hellodb"

	// Test config
	// DB_CONNECTION = "root:tinarm@tcp(host.docker.internal)/"
	// DB_NAME       = "tinarm_test"

	DB_CONNECTION_STRING = DB_CONNECTION + DB_NAME + "?charset=utf8mb4&parseTime=True&loc=Local"
)

var DB *gorm.DB

func ConnectDatabase() {

	ensureDbExists()

	db, err := gorm.Open(mysql.Open(DB_CONNECTION_STRING), &gorm.Config{})
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

func ensureDbExists() {

	db, err := sql.Open("mysql", DB_CONNECTION)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + DB_NAME)
	if err != nil {
		panic(err)
	}
}
