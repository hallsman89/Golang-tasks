/*
 * Candy Server
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type BuyCandyResponse struct {
	Change float64 `json:"change"`
	Thanks string  `json:"thanks"`
}

type BuyCandyRequest struct {

	// amount of money put into vending machine
	Money int32 `json:"money"`

	// kind of candy
	CandyType string `json:"candyType"`

	// number of candy
	CandyCount int32 `json:"candyCount"`
}

// AssertBuyCandyRequestRequired checks if the required fields are not zero-ed
func AssertBuyCandyRequestRequired(obj BuyCandyRequest) error {
	elements := map[string]interface{}{
		"money":      obj.Money,
		"candyType":  obj.CandyType,
		"candyCount": obj.CandyCount,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertBuyCandyRequestConstraints checks if the values respects the defined constraints
func AssertBuyCandyRequestConstraints(obj BuyCandyRequest) error {
	return nil
}
