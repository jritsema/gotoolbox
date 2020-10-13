package gotoolbox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//ReadJSONFile reads JSON from a file
func ReadJSONFile(fileName string) (map[string]string, error) {

	jsonFile, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v: %w", fileName, err)
	}
	defer jsonFile.Close()
	b, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v : %w", fileName, err)
	}
	var result map[string]string
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling file: %v : %w", fileName, err)
	}
	return result, nil
}
