package kinvest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/suapapa/go_kinvest/internal/oapi"
)

// ClientConfig holds the configuration for the Kinvest client
type ClientConfig struct {
	AppKey    string
	AppSecret string
	Account   string // 계좌번호 XXXXXXXX-XX
}

// NewClientConfigFromEnv creates a new ClientConfig from environment variables
func NewClientConfigFromEnv() (*ClientConfig, error) {
	appKey := apiEnvs["APPKEY"]
	appSecret := apiEnvs["APPSECRET"]
	account := apiEnvs["ACCOUNT"]
	if appKey == "" || appSecret == "" || account == "" {
		return nil, fmt.Errorf("set KINVEST_APPKEY, KINVEST_APPSECRET, KINVEST_ACCOUNT env vars")
	}
	return &ClientConfig{
		AppKey:    appKey,
		AppSecret: appSecret,
		Account:   account,
	}, nil
}

// Client is the main client for the Kinvest API
type Client struct {
	oc *oapi.Client

	appKey    string
	appSecret string
	account   string

	ip  string
	mac string

	token       string
	tokenExpiry time.Time

	hash string
}

// NewClient creates a new Kinvest client
// It uses the provided config to set up the client
// If the config is nil, it will use the environment variables
// KINVEST_APPKEY, KINVEST_APPSECRET, KINVEST_ACCOUNT, KINVEST_MAC
func NewClient(config *ClientConfig) (*Client, error) {
	ip, mac, err := getLocalIPAndMAC()
	if err != nil {
		return nil, fmt.Errorf("failed to get local IP and MAC: %w", err)
	}

	if config == nil {
		config, err = NewClientConfigFromEnv()
		if err != nil {
			return nil, fmt.Errorf("failed to create client config: %w", err)
		}
	}

	c := &Client{
		appKey:    config.AppKey,
		appSecret: config.AppSecret,
		account:   config.Account,
		ip:        ip,
		mac:       mac,
	}
	if c.appKey == "" {
		c.appKey = apiEnvs["APPKEY"]
	}
	if c.appSecret == "" {
		c.appSecret = apiEnvs["APPSECRET"]
	}
	if c.account == "" {
		c.account = apiEnvs["ACCOUNT"]
	}
	if c.appKey == "" || c.appSecret == "" || c.account == "" {
		return nil, fmt.Errorf("invalid config: appKey, appSecret, account must be set")
	}

	// addReqAuthHeader := func(ctx context.Context, req *http.Request) error {
	// 	return nil
	// }

	c.oc, err = oapi.NewClient(
		prodAddr,
		// oapi.WithRequestEditorFn(addReqAuthHeader),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create oapi client: %w", err)
	}

	if err := c.refreshToken(); err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	hash, err := c.genHashKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate hash key: %w", err)
	}
	c.hash = hash

	return c, nil
}

func (c *Client) refreshToken() error {
	// Check if the token is either expired or will expire within 1 minute
	if c.token != "" && c.tokenExpiry.After(time.Now().Add(1*time.Minute)) {
		return nil
	}
	if err := c.getToken(); err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}
	return nil
}

func (c *Client) getToken() error {
	reqBody := mustCreateJsonReader(map[string]any{
		"grant_type": "client_credentials",
		"appkey":     c.appKey,
		"appsecret":  c.appSecret,
	})
	req, err := oapi.NewPostOauth2TokenPRequestWithBody(
		c.oc.Server,
		&oapi.PostOauth2TokenPParams{
			ContentType: &jsonContentType,
		},
		jsonContentType,
		reqBody,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.oc.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	ret := mustUnmarshalJsonBody(resp.Body)
	c.token = ret["token_type"].(string) + " " + ret["access_token"].(string)
	if c.token == "" {
		return fmt.Errorf("empty token")
	}
	c.tokenExpiry, err = time.Parse(time.RFC3339, ret["access_token_expired"].(string))
	if err != nil {
		return fmt.Errorf("failed to parse token expiry: %w", err)
	}
	if c.tokenExpiry.Before(time.Now()) {
		return fmt.Errorf("token expired")
	}

	return nil
}

func (c *Client) genHashKey() (string, error) {
	reqBody := mustCreateJsonReader(map[string]any{
		"ACCOUNT":      c.account,
		"ACNT_PRDT_CD": "01",
		"OVRS_EXCG_CD": "SHAA",
	})
	req, err := oapi.NewPostUapiHashkeyRequestWithBody(
		c.oc.Server,
		&oapi.PostUapiHashkeyParams{
			ContentType: &jsonContentType,
			Appkey:      &c.appKey,
			Appsecret:   &c.appSecret,
		},
		jsonContentType,
		reqBody,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.oc.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	ret := mustUnmarshalJsonBody(resp.Body)
	hash, ok := ret["HASH"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected response format: %v", ret)
	}
	return hash, nil
}
