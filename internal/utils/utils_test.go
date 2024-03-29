package utils

import (
	"testing"
)

func TestConvertMapToString(t *testing.T) {
	testMap := make(map[string]string)
	testMap["key"] = "value"
	returnString := ConvertMapToString(testMap)
	expectedString := "key=\"value\"\n"
	if expectedString != returnString {
		t.Errorf("Did not get expected result. Got '%s', wanted '%s'", returnString, expectedString)
	}
}

func TestCheckForError(t *testing.T) {
	testForErrorSlice := []string{"i/o timeout", "context deadline exceeded", "connection refused", "...no error found...", "Bad Request"}
	for _, error := range testForErrorSlice {
		expectedTestValueString := "Error: " + error + " (validate cluster access and restart)"
		errorSlice := []string{error}
		if error == "...no error found..." {
			returnedValueString, returnedBool := CheckForError(errorSlice)
			if (returnedBool != false) || (returnedValueString != "") {
				t.Errorf("Did not get expected result. Got '%s' and '%v', wanted '%s' and '%v'",
					returnedValueString, returnedBool, expectedTestValueString, false)
			}
		} else {
			returnedValueString, returnedBool := CheckForError(errorSlice)
			if (returnedBool != true) || (returnedValueString != expectedTestValueString) {
				t.Errorf("Did not get expected result. Got '%s' and '%v', wanted '%s' and '%v'",
					returnedValueString, returnedBool, expectedTestValueString, true)
			}
		}
	}
}
