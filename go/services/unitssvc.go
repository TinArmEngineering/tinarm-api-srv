package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	env "github.com/tinarmengineering/sno2-srv-go/go/environment"
)

type ForceQuantity struct {
	From Quantity `json:"from"`
	To   Quantity `json:"to"`
}

type Quantity struct {
	Magnitude float64 `json:"magnitude"`
	Units     []Unit  `json:"units"`
}

type Unit struct {
	Name     string `json:"name"`
	Exponent int    `json:"exponent"`
}

type UnitError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
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

	json_data, err := json.Marshal(qty)
	if err != nil {
		return Quantity{}, err
	}

	return callUnitsService("convert_to_base_units", json_data)
}

// Checks that the specified 'from' Quantity can be converted to the
// specified 'to' Quantity
//
// Example input:
//
//	from/to := svc.Quantity{
//		Magnitude: 1,
//		Units: []svc.Unit{
//			{Name: "lb", Exponent: 1},
//			{Name: "foot", Exponent: -2},
//		},
//	}
func (unitsSvc UnitsSvc) ConvertTo(from Quantity, to Quantity) (Quantity, error) {

	json_data, err := json.Marshal(ForceQuantity{From: from, To: to})
	if err != nil {
		return Quantity{}, err
	}

	return callUnitsService("convert_to", json_data)
}

func callUnitsService(method string, json_data []byte) (Quantity, error) {

	res := Quantity{}

	resp, err := http.Post(
		"http://"+env.UnitsServiceHost()+"/"+method,
		"application/json",
		bytes.NewBuffer(json_data))

	if err != nil {
		return Quantity{}, err
	}

	if resp.StatusCode != http.StatusOK {
		unitError := UnitError{}
		json.NewDecoder(resp.Body).Decode(&unitError)
		return Quantity{}, fmt.Errorf("%v %v %v", resp.Status, unitError.Message, unitError.Type)
	}

	return res, json.NewDecoder(resp.Body).Decode(&res)
}
