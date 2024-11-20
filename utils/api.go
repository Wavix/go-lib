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
	apiGateway := "https://api.wavix.com"
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
		return fmt.Sprintf("https://api.%s", hostname)
	}

	return apiGateway
}
