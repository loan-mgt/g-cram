package service

import (
	"encoding/json"
	"fmt"
	"io"
	"loan-mgt/g-cram/internal/config"
	"net/http"
	"net/url"
	"strings"
)

type TokenResponse struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	Scope                 string `json:"scope"`
	TokenType             string `json:"token_type"`
	IDToken               string `json:"id_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
}

func GetTokens(cfg *config.Config, code string) (*TokenResponse, error) {
	// Define response structure

	// Prepare form data
	form := url.Values{
		"code":          {code},
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {"http://localhost:8080/"},
	}

	// Create request
	req, err := http.NewRequest("POST", cfg.TokenURI, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error response from token endpoint: %d %s",
			resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &tokenResp, nil
}
