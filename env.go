package kinvest

import (
	"cmp"
	"encoding/json"
	"os"
	"strings"
)

var (
	// modified from `실전계좌_POSTMAN_환경변수.json`
	kinvestProdEnvJson = `[
  {
    "key": "VTS",
    "value": "https://openapivts.koreainvestment.com:29443",
  },
  {
    "key": "PROD",
    "value": "https://openapi.koreainvestment.com:9443",
  },
  {
    "key": "CANO",
    "value": "",
  },
  {
    "key": "APPKEY",
    "value": "",
  },
  {
    "key": "APPSECRET",
    "value": "",
  },
  {
    "key": "TOKEN",
    "value": "",
  },
  {
    "key": "HASH",
    "value": "",
  }
]
`
	apiEnvs = map[string]string{}
)

func init() {
	initApiEnvsMust()
}

func initApiEnvsMust() {
	jd := json.NewDecoder(strings.NewReader(kinvestProdEnvJson))
	envs := make(
		[]struct {
			Key     string `json:"key"`
			Value   string `json:"value"`
			Enabled bool   `json:"enabled"`
			Type    string `json:"type"`
		},
		0,
	)
	if err := jd.Decode(&envs); err != nil {
		panic(err)
	}
	for _, env := range envs {
		if env.Value != "" {
			apiEnvs[env.Key] = env.Value
		} else {
			val := cmp.Or(os.Getenv("KINVEST_"+env.Key), "")
			apiEnvs[env.Key] = val
		}
	}
}
