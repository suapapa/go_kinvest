package kinvest

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const defaultAccessTokenPath = "./kinvest_access_token.json"

type AccessToken struct {
	TokenType   string    `json:"token_type"`
	AccessToken string    `json:"access_token"`
	ExpiresIn   time.Time `json:"expires_in"`
}

func LoadAccessToken(tokenPath string) (*AccessToken, error) {
	ret := &AccessToken{}
	f, err := os.Open(tokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open token file: %w", err)
	}
	defer f.Close()
	if err := unmarshalJsonBody(f, ret); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}
	if ret.TokenType == "" || ret.AccessToken == "" || ret.ExpiresIn.IsZero() {
		return nil, fmt.Errorf("invalid token data")
	}
	if ret.IsExpired() {
		return nil, fmt.Errorf("token expired")
	}

	return ret, nil
}

func (t *AccessToken) Save(tokenPath string) error {
	f, err := os.Create(tokenPath)
	if err != nil {
		return fmt.Errorf("failed to create token file: %w", err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(t); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	return nil
}

func (t *AccessToken) IsExpired() bool {
	expiresIn := t.ExpiresIn
	if expiresIn.IsZero() {
		return true
	}

	// Check if the token is expired or will expire within 1 minute
	expiresIn = expiresIn.Add(-1 * time.Minute)

	return expiresIn.Before(time.Now())
}

func (t *AccessToken) Authorization() string {
	ret := t.TokenType + " " + t.AccessToken
	return ret
}
