package main

import (
	"bytes"
	"fmt"
)

// used by labels, annotations, ...
func convertMapToString(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}
