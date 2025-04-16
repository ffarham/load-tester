package utils

import (
	"encoding/json"
	"math"
	"os"
)

// ReadJsonFile reads a JSON file and returns the contents as a map.
func ReadJsonFile(filename string) (interface{}, error) {

	contents, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var jsonData interface{}
	err = json.Unmarshal(contents, &jsonData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

// SafeDiv returns NaN when dividing by zero, instead of an infinity value.
func SafeDiv(a, b float64) float64 {
	if b == 0 {
		return math.NaN()
	}
	return a / b
}
