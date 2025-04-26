package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"loan-mgt/g-cram/internal/config"
	"loan-mgt/g-cram/internal/db"
	"log"
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

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func GetRefreshToken(cfg *config.Config, code string) (*TokenResponse, error) {
	// Define response structure

	fmt.Println(code)

	// Prepare form data
	form := url.Values{
		"code":          {code},
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {"http://localhost:8080/"},
	}

	fmt.Println(form.Encode())
	fmt.Println(form.Encode())
	fmt.Println(code)

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
		log.Printf("error response from token endpoint: %d %s",
			resp.StatusCode, string(bodyBytes))
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

func GetAccessToken(ctx context.Context, cfg *config.Config, db *db.Store, userId string) (string, error) {
	user, err := db.GetUser(ctx, userId)
	if err != nil {
		return "", fmt.Errorf("error getting user: %w", err)
	}

	if user.Token.String == "" {
		return "", fmt.Errorf("user has no refresh token")
	}

	form := url.Values{
		"refresh_token": {user.Token.String},
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
		"grant_type":    {"refresh_token"},
	}

	req, err := http.NewRequest("POST", cfg.TokenURI, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("error response from token endpoint: %d %s",
			resp.StatusCode, string(bodyBytes))
		return "", fmt.Errorf("error response from token endpoint: %d %s",
			resp.StatusCode, string(bodyBytes))
	}

	var tokenResp AccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("error parsing response: %w", err)
	}

	return tokenResp.AccessToken, nil
}
