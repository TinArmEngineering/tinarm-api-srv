/*
 * ta-solve
 *
 * The unnamed Tin Arm solver API
 *
 * API version: 1.0
 * Contact: api@tinarmengineering.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	dbo "github.com/tinarmengineering/tinarm-api-srv/go/dbo"
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

// DefaultApiService is a service that implements the logic for the DefaultApiServicer
// This service should implement the business logic for every endpoint for the DefaultApi API.
// Include any external packages or services that will be required by this service.
type DefaultApiService struct {
}

// NewDefaultApiService creates a default api service
func NewDefaultApiService() DefaultApiServicer {

	ensureDbExists(DB_CONNECTION, DB_NAME)

	db, err := gorm.Open(mysql.Open(DB_CONNECTION_STRING), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&dbo.Job{})

	return &DefaultApiService{}
}

// DeleteJobsId - Delete Job
func (s *DefaultApiService) DeleteJobsId(ctx context.Context, id interface{}) (ImplResponse, error) {

	var job dbo.Job
	getDb().Delete(&job, id)
	return Response(200, nil), nil
}

// GetJobsId - Get Job
func (s *DefaultApiService) GetJobsId(ctx context.Context, id interface{}) (ImplResponse, error) {

	var job dbo.Job
	getDb().First(&job, id)

	if job.ID == 0 {
		return Response(404, nil), nil
	}

	return Response(200, job), nil
}

// PostRectanglejobs - Create RectangleJob
func (s *DefaultApiService) PostRectanglejobs(ctx context.Context, rectanglejob Rectanglejob) (ImplResponse, error) {

	b, err := json.Marshal(rectanglejob.Geometry)
	if err != nil {
		panic("failed to serialise geometry")
	}

	getDb().Create(&dbo.Job{Data: string(b)})

	return Response(200, nil), nil
}

// PostStatorjobs - Create StatorJob
func (s *DefaultApiService) PostStatorjobs(ctx context.Context, statorjob Statorjob) (ImplResponse, error) {
	// TODO - update PostStatorjobs with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("PostStatorjobs method not implemented")
}

func ensureDbExists(connection string, name string) {

	db, err := sql.Open("mysql", connection)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + name)
	if err != nil {
		panic(err)
	}
}

func getDb() (db *gorm.DB) {
	db, err := gorm.Open(mysql.Open(DB_CONNECTION_STRING), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
