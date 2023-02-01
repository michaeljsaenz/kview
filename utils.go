package main

import (
	"bytes"
	"fmt"
	"strings"
)

// used by labels, annotations, ...
func convertMapToString(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

// check for string(error) in slice
func checkForError(slice []string) (string, bool) {
	checkValue := slice[0]
	listOfErrors := []string{"i/o timeout", "context deadline exceeded", "connection refused"}
	for _, error := range listOfErrors {
		if strings.Contains(checkValue, error) {
			return "Error: " + checkValue + " (validate cluster access and restart)", true
		}
	}
	return "", false
}
