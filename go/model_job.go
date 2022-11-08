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

// Job - The basic job resource
type Job struct {

	// The unique id of this job
	Id *interface{} `json:"id,omitempty"`

	Stator Statorgeometry `json:"stator"`
}

// AssertJobRequired checks if the required fields are not zero-ed
func AssertJobRequired(obj Job) error {
	elements := map[string]interface{}{
		"stator": obj.Stator,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	if err := AssertStatorgeometryRequired(obj.Stator); err != nil {
		return err
	}
	return nil
}

// AssertRecurseJobRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of Job (e.g. [][]Job), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseJobRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aJob, ok := obj.(Job)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertJobRequired(aJob)
	})
}
