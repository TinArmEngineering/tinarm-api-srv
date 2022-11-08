/*
 * ta-solve
 *
 * The unnamed Tin Arm solver API
 *
 * API version: 1.0
 * Contact: api@tinarmengineering.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"log"
	"net/http"

	openapi "github.com/tinarmengineering/tinarm-api-srv/go"
)

func main() {

	log.Printf("Server started")

	DefaultApiService := openapi.NewDefaultApiService()
	DefaultApiController := openapi.NewDefaultApiController(DefaultApiService)

	router := openapi.NewRouter(DefaultApiController)

	log.Fatal(http.ListenAndServe(":8080", router))
}
