package data_complement

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"vecfleet/features/data_complement/models"

	"net/http"
)

func FetchComplemetData(url string) (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error haciendo la solicitud: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("error read response: %v", err)
		}

		var apiResponse models.ComplementData
		if err := json.Unmarshal(body, &apiResponse); err != nil {
			return "", fmt.Errorf("error deserialized JSON: %v", err)
		}
		responseStr, _ := apiResponse.Marshal()
		return string(responseStr), nil
	} else if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("status not foiund")
	} else {
		return "", fmt.Errorf("error undefined, status code: %d", resp.StatusCode)
	}
}
