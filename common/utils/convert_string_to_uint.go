package utils

import (
	"log"
	"strconv"
)

func ConvStrToUint(s, varName string) uint64 {
	var result uint64
	if s == "" {
		result = uint64(0)
		return result
	}
	result, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		log.Printf("ERROR: [utils - ConvStrToUint] Error while convert %s string to uint: %v\n", varName, err)
		result = uint64(0)
	}
	return result
}
