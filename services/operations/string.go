package operations

import (
	"fmt"
	"strconv"

	"math/rand"
)

func Replace(input string, value string) string {
	return value
}

func RandomInt(input string, upper_limit string, lower_limit string) string {
	// Convert value to int
	upper, err := strconv.Atoi(upper_limit)
	if err != nil {
		return input
	}
	lower, err := strconv.Atoi(lower_limit)
	if err != nil {
		return input
	}
	return fmt.Sprintf("%d", rand.Intn(upper-lower)+lower)
}
