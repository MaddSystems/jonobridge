package data_complement_test

import (
	"fmt"
	"os"
	"recursoconfiable/features/data_complement"
	"testing"
)

func TestFetchComplemetData(t *testing.T) {
	COMPLEMENT_URL := os.Getenv("PLATES_URL")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("", COMPLEMENT_URL)
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	// COMPLEMENT_URL := "https://pluto.dudewhereismy.com.mx/imei/search?appId=2911"
	apiResponse, err := data_complement.FetchComplemetData(COMPLEMENT_URL)
	if err != nil {
		t.Errorf("Expected nil, but got error: %v", err)
	} else {

		t.Logf("Test passed successfully with no errors. Response: %s", apiResponse)
	}

}
