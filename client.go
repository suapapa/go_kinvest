package kinvest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/suapapa/go_kinvest/internal/oapi"
)

type ClientConfig struct {
	AppKey    string
	AppSecret string
	Account   string // 계좌번호
	Mac       string // MAC 주소
}

func NewClientConfigFromEnv() (*ClientConfig, error) {
	appKey := apiEnvs["APPKEY"]
	appSecret := apiEnvs["APPSECRET"]
	account := apiEnvs["CANO"]
	mac := apiEnvs["MAC"]
	if appKey == "" || appSecret == "" || account == "" {
		return nil, fmt.Errorf("set KINVEST_APPKEY, KINVEST_APPSECRET, KINVEST_CANO env vars")
	}
	return &ClientConfig{
		AppKey:    appKey,
		AppSecret: appSecret,
		Account:   account,
		Mac:       mac,
	}, nil
}

type Client struct {
	oc *oapi.Client

	appKey    string
	appSecret string
	account   string
	mac       string

	token       string
	tokenExpiry time.Time

	hash string
}

func NewClient(config *ClientConfig) (*Client, error) {
	var err error
	c := &Client{
		appKey:    config.AppKey,
		appSecret: config.AppSecret,
		account:   config.Account,
		mac:       config.Mac,
	}
	c.oc, err = oapi.NewClient(
		prodAddr,
		// oapi.WithRequestEditorFn(addReqAuthHeader),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create oapi client: %w", err)
	}

	if err := c.getToken(); err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	hash, err := c.genHashKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate hash key: %w", err)
	}
	c.hash = hash

	return c, nil
}

// func addReqAuthHeader(ctx context.Context, req *http.Request) error {
// 	// TBD
// 	return nil
// }

func (c *Client) getToken() error {
	reqBody := newJsonReaderMust(map[string]any{
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
	ret := unmarshalJsonRespBodyMust(resp.Body)
	c.token = ret["access_token"].(string)
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
	reqBody := newJsonReaderMust(map[string]any{
		"CANO":         c.account,
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
	ret := unmarshalJsonRespBodyMust(resp.Body)
	hash, ok := ret["HASH"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected response format: %v", ret)
	}
	return hash, nil
}
