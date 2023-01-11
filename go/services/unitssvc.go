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
	unitMap  map[string][]Unit
	fieldMap map[string]string
}

func (unitsSvc UnitsSvc) getMaps() (map[string]string, map[string][]Unit) {

	if len(unitsSvc.unitMap) == 0 {

		unitsSvc.unitMap = make(map[string][]Unit)

		unitsSvc.unitMap["displacement"] = []Unit{{Name: "m", Exponent: 1}}
		unitsSvc.unitMap["mass"] = []Unit{{Name: "kg", Exponent: 1}}
		unitsSvc.unitMap["area"] = []Unit{{Name: "m", Exponent: 2}}
		unitsSvc.unitMap["volume"] = []Unit{{Name: "m", Exponent: 3}}
		unitsSvc.unitMap["force"] = []Unit{{Name: "N", Exponent: 1}}

		unitsSvc.unitMap["pressure"] = []Unit{
			{Name: "N", Exponent: 1},
			{Name: "m", Exponent: -2}}

		unitsSvc.unitMap["volumetric_mass_density"] = []Unit{
			{Name: "kg", Exponent: 1},
			{Name: "m", Exponent: -3}}

		unitsSvc.unitMap["thermal_conductivity"] = []Unit{
			{Name: "W", Exponent: 1},
			{Name: "m", Exponent: -1},
			{Name: "K", Exponent: -1}}

		unitsSvc.unitMap["electrical_conductivity"] = []Unit{
			{Name: "S", Exponent: 1},
			{Name: "m", Exponent: -1}}

		unitsSvc.unitMap["magnetic_permeability"] = []Unit{
			{Name: "H", Exponent: 1},
			{Name: "m", Exponent: -1}}

		unitsSvc.unitMap["permittivity"] = []Unit{
			{Name: "F", Exponent: 1},
			{Name: "m", Exponent: -1}}

		unitsSvc.unitMap["heat_capacity"] = []Unit{
			{Name: "J", Exponent: 1},
			{Name: "K", Exponent: -1}}

		unitsSvc.fieldMap = make(map[string]string)

		// Fields from Material
		unitsSvc.fieldMap["width"] = "displacement"
		unitsSvc.fieldMap["depth"] = "displacement"
		unitsSvc.fieldMap["depth"] = "displacement"
		unitsSvc.fieldMap["heat_conductivity"] = "thermal_conductivity"
		unitsSvc.fieldMap["electric_conductivity"] = "electrical_conductivity"
		unitsSvc.fieldMap["relative_permeability"] = "magnetic_permeability"
		unitsSvc.fieldMap["relative_permittivity"] = "permittivity"
		unitsSvc.fieldMap["heat_capacity"] = "heat_capacity"
		unitsSvc.fieldMap["density"] = "volumetric_mass_density"
	}

	return unitsSvc.fieldMap, unitsSvc.unitMap
}

func (unitsSvc UnitsSvc) StandardFieldUnits(fieldName string) []Unit {
	fieldMap, unitMap := unitsSvc.getMaps()
	if val, ok := fieldMap[fieldName]; ok {
		return unitMap[val]
	}
	return []Unit{}
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

	return callUnitsService("convert_to_base_units", qty)
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

	return callUnitsService("convert_to", ForceQuantity{From: from, To: to})
}

func callUnitsService(method string, v any) (Quantity, error) {

	result := Quantity{}

	json_data, err := json.Marshal(v)
	if err != nil {
		return Quantity{}, err
	}

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

	return result, json.NewDecoder(resp.Body).Decode(&result)
}
