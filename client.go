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

	ip      string
	mac     string
	account string

	appKey    string
	appSecret string
	token     *AccessToken
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

	return c, nil
}

func (c *Client) refreshToken() error {
	if c.token != nil && !c.token.IsExpired() {
		return nil
	}

	tokenPath := apiEnvs["TOKEN_PATH"]
	if tokenPath == "" {
		tokenPath = defaultAccessTokenPath
	}

	var err error
	if c.token == nil || c.token.IsExpired() {
		if fileExists(tokenPath) {
			c.token, err = LoadAccessToken(tokenPath)
			if err == nil {
				return nil
			}
		}
		c.token, err = c.getToken()
		if err != nil {
			return fmt.Errorf("failed to get token: %w", err)
		}
	}

	err = c.token.Save(tokenPath)
	if err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	return nil
}

func (c *Client) getToken() (*AccessToken, error) {
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
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.oc.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	data := mustUnmarshalJsonBody(resp.Body)

	return &AccessToken{
		TokenType:   data["token_type"].(string),
		AccessToken: data["access_token"].(string),
		ExpiresIn:   time.Now().Add(time.Duration((data["expires_in"].(float64))) * time.Second),
	}, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	err := c.refreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return c.oc.Client.Do(req)
}
