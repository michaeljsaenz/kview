package utils

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

func ConvertMapToString(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

// check for string(error) in slice
func CheckForError(slice []string) (string, bool) {
	if slice != nil {
		checkValue := slice[0]
		listOfErrors := []string{"i/o timeout", "context deadline exceeded", "connection refused", "Bad Request"}
		for _, error := range listOfErrors {
			if strings.Contains(checkValue, error) {
				return "Error: " + checkValue + " (validate cluster access and restart)", true
			}
		}
	}
	return "", false
}

// return pointer to int64
func CreateInt64(num int64) *int64 {
	return &num
}

func RemoveANSIEscapeCodes(input string) string {
	regex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	clean := regex.ReplaceAllString(input, "")
	return clean
}
