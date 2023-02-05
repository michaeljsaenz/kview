package main

import (
	"testing"
)

func TestConvertMapToString(t *testing.T) {
	testMap := make(map[string]string)
	testMap["key"] = "value"
	returnString := convertMapToString(testMap)
	expectedString := "key=\"value\"\n"
	if expectedString != returnString {
		t.Errorf("Did not get expected result. Got '%s', wanted '%s'", returnString, expectedString)
	}
}

func TestCheckForError(t *testing.T) {
	testErrorSlice := []string{"i/o timeout", "context deadline exceeded", "connection refused"}
	for _, error := range testErrorSlice {
		expectedTestValueString := "Error: " + error + " (validate cluster access and restart)"
		expectedTestErrorExistsTrue := true
		errorSlice := []string{error}
		returnedValueString, returnedBool := checkForError(errorSlice)
		if (expectedTestErrorExistsTrue != returnedBool) || (expectedTestValueString != returnedValueString) {
			t.Errorf("Did not get expected result. Got '%s' and '%v', wanted '%s' and '%v'", returnedValueString, returnedBool,
				expectedTestValueString, expectedTestErrorExistsTrue)
		}
	}
}
