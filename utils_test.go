package main

import (
	"testing"
)

func TestConvertMapToString(t *testing.T) {
	testMap := make(map[string]string)
	testMap["key"] = "value"
	expectedString := "key=\"value\"\n"
	returnString := convertMapToString(testMap)
	if expectedString != returnString {
		t.Errorf("Did not get expected result. Got '%s', wanted '%s'", returnString, expectedString)
	}
}
