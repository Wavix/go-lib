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

func TestGetAuthServicePath(t *testing.T) {
	tests := []struct {
		name     string
		hostname string
		expected string
	}{
		{
			name:     "default URL when hostname does not contain 'qa'",
			hostname: "app01.zone.wavix.com",
			expected: "https://private-api.wavix.com",
		},
		{
			name:     "qa environment with 3 parts",
			hostname: "qa123.wavix.zone",
			expected: "https://private-api.qa123.wavix.zone",
		},
		{
			name:     "complex hostname with qa - extract last 3 parts",
			hostname: "some1.some2.some3.qa123.wavix.zone",
			expected: "https://private-api.qa123.wavix.zone",
		},
		{
			name:     "default when hostname does not contain qa",
			hostname: "app01.nl.wavix.com",
			expected: "https://private-api.wavix.com",
		},
		{
			name:     "default when hostname has less than 3 parts",
			hostname: "tts.wavix.com",
			expected: "https://private-api.wavix.com",
		},
		{
			name:     "default with wavix.com domain",
			hostname: "wavix.com",
			expected: "https://private-api.wavix.com",
		},
		{
			name:     "default with subdomain wavix.com",
			hostname: "some.wavix.com",
			expected: "https://private-api.wavix.com",
		},
		{
			name:     "qa environment with wavix.dev",
			hostname: "qa1.wavix.dev",
			expected: "https://private-api.qa1.wavix.dev",
		},
		{
			name:     "qa subdomain with app prefix",
			hostname: "app01.qa1.wavix.dev",
			expected: "https://private-api.qa1.wavix.dev",
		},
		{
			name:     "extract last 3 parts from hostname with qa",
			hostname: "server1.qa1.wavix.dev",
			expected: "https://private-api.qa1.wavix.dev",
		},
		{
			name:     "complex hostname with qa",
			hostname: "some.host.qa1.wavix.dev",
			expected: "https://private-api.qa1.wavix.dev",
		},
		{
			name:     "default when hostname is empty",
			hostname: "",
			expected: "https://private-api.wavix.com",
		},
		{
			name:     "qa with fewer than 3 parts",
			hostname: "qa1.wavix",
			expected: "https://private-api.qa1.wavix",
		},
		{
			name:     "qa with exactly 3 parts",
			hostname: "qa1.wavix.zone",
			expected: "https://private-api.qa1.wavix.zone",
		},
		{
			name:     "numbers in hostname without qa",
			hostname: "app123.prod456.wavix.com",
			expected: "https://private-api.wavix.com",
		},
		{
			name:     "qa with numbers in hostname",
			hostname: "app123.qa456.wavix.zone",
			expected: "https://private-api.qa456.wavix.zone",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProvider := MockHostnameProvider{MockHostname: tt.hostname}
			result := utils.GetAuthServicePath(mockProvider)
			assert.Equal(t, tt.expected, result)
		})
	}
}
