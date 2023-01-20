package main

import (
	"testing"

	"k8s.io/client-go/tools/clientcmd"
)

func TestGetCurrentContext(t *testing.T) {
	testClientConfig, _ := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{
			CurrentContext: "",
		}).RawConfig()
	expectedCurrentContext := testClientConfig.CurrentContext
	currentContext := getCurrentContext()

	if expectedCurrentContext != currentContext {
		t.Errorf("Did not get expected result. Got '%s', wanted '%s'", currentContext, expectedCurrentContext)
	}
}
