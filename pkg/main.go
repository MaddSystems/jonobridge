package main

import (
    "fmt"
    "strings"
    "strconv"
)

func main() {
    a := "$$m158,864352045580768,AAA,35,19.563274155339307,-99.15777504444124,250214233431,A,9,12,98,76,1,2239,0,1348,0|0|0000|0000,0000,0000|0000|0000|80|0000,00000000,*80"

    parts := strings.Split(a, ",")
    
    // Get string values
    latitudeStr := parts[4]
    longitudeStr := parts[5]

    // Convert to float64
    latitudeFloat, _ := strconv.ParseFloat(latitudeStr, 64)
    longitudeFloat, _ := strconv.ParseFloat(longitudeStr, 64)

    // Print both string and float values
    fmt.Printf("Latitude  (string): %s\n", latitudeStr)
    fmt.Printf("Latitude  ( float): %.9f\n\n", latitudeFloat)
    
    fmt.Printf("Longitude (string): %s\n", longitudeStr)
    fmt.Printf("Longitude ( float): %.9f\n", longitudeFloat)
}