package kinvest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	token     *accessToken
}

// NewClient creates a new Kinvest client
// It uses the provided config to set up the client
// If the config is nil, it will use the environment variables
// KINVEST_APPKEY, KINVEST_APPSECRET, KINVEST_ACCOUNT
// and KINVEST_TOKEN_PATH (optional) to save the access token
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

	fillHeader := func(ctx context.Context, req *http.Request) error {
		if c.token != nil {
			auth := c.token.Authorization()
			if auth != "" {
				req.Header.Set("authorization", c.token.Authorization())
			}
		}
		req.Header.Set("appkey", c.appKey)
		req.Header.Set("appsecret", c.appSecret)

		return nil
	}
	refreshToken := func(ctx context.Context, req *http.Request) error {
		switch {
		case strings.Contains(req.URL.Path, "/oauth2/revokeP"):
			return nil
		case strings.Contains(req.URL.Path, "/oauth2/tokenP"):
			return nil
		}

		return c.refreshToken(ctx)
	}
	c.oc, err = oapi.NewClient(
		prodAddr,
		oapi.WithRequestEditorFn(refreshToken),
		oapi.WithRequestEditorFn(fixCodeLen),
		oapi.WithRequestEditorFn(fillHeader),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create oapi client: %w", err)
	}

	return c, nil
}

func (c *Client) refreshToken(ctx context.Context) error {
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
			c.token, err = loadAccessToken(tokenPath)
			if err == nil {
				return nil
			}
		}
		c.token, err = c.getToken(ctx)
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

func (c *Client) getToken(ctx context.Context) (*accessToken, error) {
	resp, err := c.oc.PostOauth2TokenP(
		ctx,
		&oapi.PostOauth2TokenPParams{},
		oapi.PostOauth2TokenPJSONRequestBody{
			"grant_type": "client_credentials",
			"appkey":     c.appKey,
			"appsecret":  c.appSecret,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	data := mustUnmarshalJsonBody(resp.Body)

	return &accessToken{
		TokenType:   data["token_type"].(string),
		AccessToken: data["access_token"].(string),
		ExpiresIn:   time.Now().Add(time.Duration((data["expires_in"].(float64))) * time.Second),
	}, nil
}

func fixCodeLen(ctx context.Context, req *http.Request) error {
	codes := []string{"ACNT_PRDT_CD", "INQR_DVSN", "UNPR_DVSN", "PRCS_DVSN"}

	// find query parm "ACNT_PRDTCD" and change its value's len to 2. e. g. "0" -> "01"
	for _, code := range codes {
		val := req.URL.Query().Get(code)
		if val == "" {
			continue
		} else if len(val) == 1 {
			val = "0" + val
			query := req.URL.Query()
			query.Set(code, val)
			req.URL.RawQuery = query.Encode()
		}
	}

	// fix code len in body
	// from : {"ACNT_PRDT_CD":1,"CANO":"64632233","ORD_DVSN":"01","ORD_QTY":"1","ORD_UNPR":"0","PDNO":"005380"}
	// to : {"ACNT_PRDT_CD":"01","CANO":"64632233","ORD_DVSN":"01","ORD_QTY":"1","ORD_UNPR":"0","PDNO":"005380"}
	var bodyData map[string]any
	if req.Body == nil {
	} else {
		err := json.NewDecoder(req.Body).Decode(&bodyData)
		if err != nil {
			return fmt.Errorf("failed to decode request body: %w", err)
		}
		req.Body.Close()

		for k, v := range bodyData {
			for _, code := range codes {
				if k == code {
					if val, ok := v.(string); ok {
						switch len(val) {
						case 1:
							bodyData[k] = "0" + val
						case 2:
							bodyData[k] = val
						default:
							return fmt.Errorf("invalid code len for %s: %s", k, val)
						}
					} else if val, ok := v.(float64); ok && len(fmt.Sprintf("%d", int(val))) == 1 {
						bodyData[k] = fmt.Sprintf("%02d", int(val))
					} else {
						bodyData[k] = v
					}
				}
			}
		}

		bodyBytes, err := json.Marshal(bodyData)
		if err != nil {
			return fmt.Errorf("failed to encode request body: %w", err)
		}
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		req.ContentLength = int64(len(bodyBytes))
	}

	return nil
}
