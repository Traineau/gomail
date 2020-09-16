package helpers

import (
	"fmt"
	"strconv"
)

// ParseInt64 helper to avoid code repetition
func ParseInt64(stringToParse string) (int64, error) {
	intID, err := strconv.ParseInt(stringToParse, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse string to int")
	}
	return intID, nil
}
