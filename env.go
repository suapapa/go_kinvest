package kinvest

import (
	"cmp"
	"os"
)

const (
	vtsAddr  = "https://openapivts.koreainvestment.com:29443"
	prodAddr = "https://openapi.koreainvestment.com:9443"
)

var (
	apiEnvs = map[string]string{
		"APPKEY":     "",
		"APPSECRET":  "",
		"ACCOUNT":    "",
		"TOKEN_PATH": "",
	}
)

func init() {
	initApiEnvsMust()
}

func initApiEnvsMust() {
	for k := range apiEnvs {
		val := cmp.Or(os.Getenv("KINVEST_"+k), "")
		apiEnvs[k] = val
	}
}
