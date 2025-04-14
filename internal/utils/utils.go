package utils

import (
	"encoding/json"
	"os"
)

// readJsonFile reads a JSON file and returns the contents as a map.
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
