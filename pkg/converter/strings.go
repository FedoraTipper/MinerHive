package converter

import (
	"strconv"
	"strings"
)

func StringToIntSliceWithDashDelimiter(stringList string) []int {
	return StringToIntSlice(stringList, "-")
}

func StringToIntSlice(stringList, delimiter string) []int {
	split := strings.Split(stringList, delimiter)
	var intSlice []int
	for _, t := range split {
		intVal, _ := strconv.Atoi(t)
		intSlice = append(intSlice, intVal)
	}

	return intSlice
}
