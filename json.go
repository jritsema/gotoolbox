package gotoolbox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

//ReadJSONFile reads JSON from a file
func ReadJSONFile(fileName string) (map[string]string, error) {

	jsonFile, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %s: %w", fileName, err)
	}
	defer jsonFile.Close()
	b, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %s : %w", fileName, err)
	}
	var result map[string]string
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling file: %s : %w", fileName, err)
	}
	return result, nil
}

// HttpGetJSON fetches the contents of the given URL
// and decodes it as JSON into the given result,
// which should be a pointer to the expected data.
func HttpGetJSON(url string, result interface{}) error {

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("cannot fetch URL %q: %v", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected http GET status: %s", resp.Status)
	}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return fmt.Errorf("cannot decode JSON: %v", err)
	}
	return nil
}
