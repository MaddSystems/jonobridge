package data_complement_test

import (
	"testing"
	"vecfleet/features/data_complement"
)

func TestFetchComplemetData(t *testing.T) {
	COMPLEMENT_URL := "https://pluto.dudewhereismy.com.mx/imei/search?appId=1141"
	apiResponse, err := data_complement.FetchComplemetData(COMPLEMENT_URL)
	if err != nil {
		t.Errorf("Expected nil, but got error: %v", err)
	} else {

		t.Logf("Test passed successfully with no errors. Response: %s", apiResponse)
	}

}
