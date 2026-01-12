package helpers

import (
	"encoding/hex"
	"strconv"
	"time"
)

// TODO: Manage de different possible errors

func DateAndTime(hexString string) any {
	intValue, ok := HexToLittleEndianDecimal(hexString).(int)
	if !ok {
		return nil
	}
	// Initial date: 1 january 2000, 00:00:00 am
	startTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	finalTime := startTime.Add(time.Duration(intValue) * time.Second)
	return finalTime
}

func LatLngValue(hexString string) any {
	intValue := HexToLittleEndianDecimal(hexString).(int)
	bitSize := len(hexString) * 4 // Cada caracter hexadecimal representa 4 bits
	signedIntValue := TwosComplement(intValue, bitSize)
	decimalValue := float64(signedIntValue) / 1000000
	return decimalValue
}

func TwosComplement(value int, bitSize int) int {
	if (value & (1 << (bitSize - 1))) != 0 {
		value = value - (1 << bitSize)
	}
	return value
}

func Percentage(hexString string) any {
	intValue, ok := HexToLittleEndianDecimal(hexString).(int)
	if !ok {
		return nil
	}

	// Convertir a string y asegurar que tenga al menos 2 dígitos
	stringValue := strconv.Itoa(intValue)
	if len(stringValue) < 2 {
		stringValue = "0" + stringValue
	}

	// Insertar el punto decimal
	return stringValue[:len(stringValue)-2] + "." + stringValue[len(stringValue)-2:]
}

func DivideByThousand(hexString string) any {
	intValue, ok := HexToLittleEndianDecimal(hexString).(int)
	if !ok {
		return nil
	}
	return intValue / 1000
}

func DivideByHundred(hexString string) any {
	intValue, ok := HexToLittleEndianDecimal(hexString).(int)
	if !ok {
		return nil
	}
	return intValue / 100
}

func DivideByTen(hexString string) any {
	intValue, ok := HexToLittleEndianDecimal(hexString).(int)
	if !ok {
		return nil
	}
	return intValue / 10
}

func BooleanValue(hexString string) any {
	intValue := HexToLittleEndianDecimal(hexString)
	if intValue == 1 {
		return true
	} else {
		return false
	}
}

func HexToBinary(hexString string) any {
	decimalValue, err := strconv.ParseUint(hexString, 16, 64)
	if err != nil {
		return err
	}
	binaryString := strconv.FormatUint(decimalValue, 2)
	return binaryString
}

func HexLittleEndianToBinary(hexString string) any {
	value := HexToLittleEndianDecimal(hexString).(int)
	value64 := uint64(value)
	binaryString := strconv.FormatUint(value64, 2)
	return binaryString
}

func HexToLittleEndianDecimal(hexString string) any {
	littleEndian := ""
	for i := len(hexString) - 2; i >= 0; i -= 2 {
		littleEndian += hexString[i : i+2]
	}
	return HexToInt(littleEndian)
}

func HexToInt(hexString string) any {
	decimalValue, _ := strconv.ParseInt(hexString, 16, 64)
	return int(decimalValue)
}

func HexToUTF8(hexString string) (string, error) {

	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return "", err
	}

	// Conversion to UTF-8
	return string(bytes), nil
}
