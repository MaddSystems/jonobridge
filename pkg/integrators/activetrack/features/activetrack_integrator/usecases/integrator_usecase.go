package usecases

import (
	"activetrack/features/activetrack_integrator/data"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MaddSystems/jonobridge/common/utils"
)

func RunIntegrator(imei, dataStr string) error {
	// First get the vehicle plates for this IMEI from our JSON data

	platesStr, err := utils.GetPlates(imei)
	if err != nil {
		platesStr = "Desconocido" // Use a default value instead of empty string
	}
	utils.VPrint("Plates: %s", platesStr)
	//}

	// Now get the spoofed IMEI
	// Try to use the file in the container's root directory first
	jsonPath := "/data_spoof.json"

	// Also check the path within jonobridge as a fallback
	fallbackPath := filepath.Join(os.Getenv("HOME"), "jonobridge", "pkg", "integrators", "activetrack", "data_spoof.json")

	// Try container path first, then fallback path
	imei_spoof, spooferror := utils.GetSpoofimeiFromJson(imei, jsonPath)

	if spooferror != nil {
		// If container path failed, try the fallback path
		utils.VPrint("First lookup failed, trying fallback path...")
		imei_spoof, spooferror = utils.GetSpoofimeiFromJson(imei, fallbackPath)

		if spooferror != nil {
			utils.VPrint("JSON lookup failed: %v, trying API fallback...", spooferror)
			// Fallback to the original method if our JSON lookup fails
			imei_spoof, spooferror = utils.GetImei_spoof(imei)
			if spooferror != nil {
				//utils.VPrint("Error getting IMEI spoof: %v", spooferror)
				//utils.VPrint("No spoofed IMEI found, not sending data")
				return fmt.Errorf("no spoofed IMEI found for IMEI %s", imei)
			}
		}
	}
	utils.VPrint("IMEI spoofed: %s", imei_spoof)

	// Replace the original IMEI with the spoofed IMEI in the dataStr
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(dataStr), &jsonData); err != nil {
		utils.VPrint("Error unmarshalling JSON data: %v", err)
		return fmt.Errorf("error unmarshalling JSON data: %w", err)
	}

	// Replace the IMEI in the JSON data
	jsonData["IMEI"] = imei_spoof

	// Marshal back to JSON string
	modifiedDataStr, err := json.Marshal(jsonData)
	if err != nil {
		utils.VPrint("Error marshalling modified JSON data: %v", err)
		return fmt.Errorf("error marshalling modified JSON data: %w", err)
	}

	utils.VPrint("Replaced original IMEI %s with spoofed IMEI %s in data", imei, imei_spoof)

	// Send the data via SOAP with error handling
	if err := data.ProcessAndSendActiveTrack(platesStr, string(modifiedDataStr)); err != nil {
		utils.VPrint("Error processing and sending data: %v", err)
		return fmt.Errorf("error processing and sending data: %w", err)
	}

	return nil
}
