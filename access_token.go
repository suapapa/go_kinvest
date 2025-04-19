package kinvest

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/goccy/go-yaml"
)

var defaultAccessTokenPath string

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("failed to get current working directory: %w", err))
	}
	defaultAccessTokenPath = path.Join(wd, "kinvest_access_token.yaml")
}

type accessToken struct {
	TokenType   string    `json:"token_type" yaml:"token_type"`
	AccessToken string    `json:"access_token" yaml:"access_token"`
	ExpiresIn   time.Time `json:"expires_in" yaml:"expires_in"`
}

func LoadAccessToken(tokenPath string) (*accessToken, error) {
	ret := &accessToken{}
	f, err := os.Open(tokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open token file: %w", err)
	}
	defer f.Close()

	switch ext := getExt(tokenPath); ext {
	case "json":
		if err := unmarshalJsonBody(f, ret); err != nil {
			return nil, fmt.Errorf("failed to unmarshal token: %w", err)
		}
	case "yaml", "yml":
		if err := unmarshalYamlBody(f, ret); err != nil {
			return nil, fmt.Errorf("failed to unmarshal token: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}

	if ret.TokenType == "" || ret.AccessToken == "" || ret.ExpiresIn.IsZero() {
		return nil, fmt.Errorf("invalid token data")
	}
	if ret.IsExpired() {
		return nil, fmt.Errorf("token expired")
	}

	return ret, nil
}

func (t *accessToken) Save(tokenPath string) error {
	f, err := os.Create(tokenPath)
	if err != nil {
		return fmt.Errorf("failed to create token file: %w", err)
	}
	defer f.Close()

	switch ext := getExt(tokenPath); ext {
	case "json":
		if err := json.NewEncoder(f).Encode(t); err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}
	case "yaml", "yml":
		if err := yaml.NewEncoder(f).Encode(t); err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}
	default:
		return fmt.Errorf("unsupported file extension: %s", ext)
	}

	return nil
}

func (t *accessToken) IsExpired() bool {
	expiresIn := t.ExpiresIn
	if expiresIn.IsZero() {
		return true
	}

	// Check if the token is expired or will expire within 1 minute
	expiresIn = expiresIn.Add(-1 * time.Minute)

	return expiresIn.Before(time.Now())
}

func (t *accessToken) Authorization() string {
	ret := t.TokenType + " " + t.AccessToken
	return ret
}
