package main

import (
	"testing"

	"github.com/go-playground/assert"
	"github.com/wavix/go-lib/utils"
)

type MockHostnameProvider struct {
	MockHostname string
}

func (m MockHostnameProvider) Hostname() (string, error) {
	return m.MockHostname, nil
}

func TestGetAuthServicePathWithoutENV(t *testing.T) {
	mockProvider := MockHostnameProvider{MockHostname: "app01.zone.wavix.net"}
	gateway := utils.GetAuthServicePath(mockProvider)
	assert.Equal(t, "https://api.wavix.com", gateway)
}

func TestGetAuthServicePathQAEnv(t *testing.T) {
	mockProvider := MockHostnameProvider{MockHostname: "qa123.wavix.zone"}
	gateway := utils.GetAuthServicePath(mockProvider)
	assert.Equal(t, "https://api.qa123.wavix.zone", gateway)
}
