package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type HostnameProvider interface {
	Hostname() (string, error)
}

type DefaultHostnameProvider struct{}

func (d DefaultHostnameProvider) Hostname() (string, error) {
	return os.Hostname()
}

func GetAuthServicePath(provider HostnameProvider) string {
	apiGateway := "https://private.api.wavix.com"
	authService := os.Getenv("AUTH_SERVICE")
	if authService != "" {
		return authService
	}

	hostname, err := provider.Hostname()
	if err != nil {
		return apiGateway
	}

	re := regexp.MustCompile(`\d+`)
	hostnameFormatted := re.ReplaceAllString(hostname, "")

	if strings.Contains(hostnameFormatted, "qa.") {
		parts := strings.Split(hostname, ".")
		if len(parts) < 3 {
			return fmt.Sprintf("https://private.api.%s", hostname)
		}

		secondLevelDomain := strings.Join(parts[len(parts)-3:], ".")
		return fmt.Sprintf("https://private.api.%s", secondLevelDomain)
	}

	return apiGateway
}
