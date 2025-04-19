# kinvest : A Go package for 한국투자증권 client

![logo](_img/go_kinvest.png)

## Usage

Install pkg:
```sh
go get github.com/suapapa/go_kinvest
```

Example usage:
```go
import kinvest "github.com/suapapa/go_kinvest"
// ...
kc, _ := kinvest.NewClient(nil)
bal, _ := kc.GetDomesticAccountBalance()
```

This package read following envs for the API calls:
- `KINVEST_ACCOUNT` : 계좌번호, XXXXXXXX-XX
- `KINVEST_APPKEY` : 한국투자증권 개발자센터에서 발급받은 appkey
- `KINVEST_APPSECRET` : 한국투자증권 개발자센터에서 발급받은 appsecret
- `KINVEST_TOKEN_PATH` : 발급받은 토큰을 저장하기 위한 경로. 설정하지 않으면 `./kinvest_access_token.json` 에 저장

## Reference
- [한국투자 OpenAPI](https://apiportal.koreainvestment.com/apiservice) - API문서
