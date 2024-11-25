package auth

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/wavix/go-lib/logger"
	"github.com/wavix/go-lib/utils"
)

type PublicAuthUser struct {
	ID int
}

type PublicAuthResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message,omitempty"`
	User    PublicAuthUser `json:"user,omitempty"`
}

const (
	PublicAuthGenericError = "Access denied"
)

var log = logger.New("PublicAuth")

var transport = &http.Transport{
	TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	DisableKeepAlives:     false,
	IdleConnTimeout:       90 * time.Second,
	MaxIdleConns:          500,
	MaxIdleConnsPerHost:   100,
	ForceAttemptHTTP2:     true,
	ExpectContinueTimeout: 1 * time.Second,
}

var sharedClient = &http.Client{Transport: transport}

func Public(appid string, ipAddr string) PublicAuthResponse {
	gw := utils.GetAuthServicePath(utils.DefaultHostnameProvider{})
	authServicePath := fmt.Sprintf("%s/private/auth/public?appid=%s&ip=%s", gw, appid, ipAddr)

	req, err := http.NewRequest("GET", authServicePath, nil)
	if err != nil {
		log.Error().Msgf("Failed to create request: %v", err)
		return errorResponse(PublicAuthGenericError)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := sharedClient.Do(req)
	if err != nil {
		log.Error().Msgf("Failed to send request: %v", err)
		return errorResponse(PublicAuthGenericError)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Msgf("Failed to read response body: %v", err)
		return errorResponse(PublicAuthGenericError)
	}

	var authResponseBody PublicAuthResponse
	err = json.Unmarshal(body, &authResponseBody)
	if err != nil {
		log.Error().Msgf("Failed to unmarshal response body: %v", err)
		return errorResponse(PublicAuthGenericError)
	}

	if resp.StatusCode != http.StatusOK {
		log.Error().Msgf("Failed to send request: %s", authResponseBody.Message)
		return errorResponse(PublicAuthGenericError)
	}

	return authResponseBody
}

func errorResponse(msg string) PublicAuthResponse {
	return PublicAuthResponse{
		Success: false,
		Message: msg,
	}
}
