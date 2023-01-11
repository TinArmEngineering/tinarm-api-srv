package services

import (
	"bytes"
	"encoding/json"
	"net/http"

	env "github.com/tinarmengineering/sno2-srv-go/go/environment"
)

type Quantity struct {
	Magnitude float64 `json:"magnitude"`
	Units     []Unit  `json:"units"`
}

type Unit struct {
	Name     string `json:"name"`
	Exponent int    `json:"exponent"`
}

type UnitsSvc struct {
}

// Converts a quantity of whatever into an equivalent quantity of SI units
//
// Example input:
//
//	qty := svc.Quantity{
//		Magnitude: 1,
//		Units: []svc.Unit{
//			{Name: "lb", Exponent: 1},
//			{Name: "foot", Exponent: -2},
//		},
//	}
func (unitsSvc UnitsSvc) ConvertToBaseUnits(qty Quantity) (Quantity, error) {

	res := Quantity{}
	json_data, err := json.Marshal(qty)

	if err != nil {
		return res, err
	}

	url := "http://" + env.UnitsServiceHost() + "/convert_quantity"
	contentType := "application/json"
	resp, err := http.Post(url, contentType, bytes.NewBuffer(json_data))

	if err != nil {
		return res, err
	}

	json.NewDecoder(resp.Body).Decode(&res)

	return res, err
}
