package gotoolbox

import (
	"bytes"
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
// returns an error if http response code is not 200
func HttpGetJSON(url string, result interface{}) error {

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("cannot fetch URL %q: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected http GET status: %s", resp.Status)
	}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return fmt.Errorf("cannot decode JSON: %w", err)
	}
	return nil
}

// HttpPutJSON encodes a struct as JSON and
// HTTP PUTs it to the specified endpoint
// returns an error if http response code is not 200
func HttpPutJSON(url string, o interface{}) error {

	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(o)
	req, err := http.NewRequest(http.MethodPut, url, payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unexpected http PUT error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected http PUT status: %s", resp.Status)
	}
	return nil
}

// HttpPostJSON encodes a struct as JSON and
// HTTP POSTs it to the specified endpoint
// returns an error if http response code does not match specified
func HttpPostJSON(url string, o interface{}, httpStatusCodeToCheck int) error {

	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(o)
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unexpected http PUT error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != httpStatusCodeToCheck {
		return fmt.Errorf("unexpected http POST status: %s", resp.Status)
	}
	return nil
}
