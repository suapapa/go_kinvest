# kinvest : A Go package for 한국투자증권 client

**!! WIP !!**

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
- `KINVEST_APPSECRET` :  한국투자증권 개발자센터에서 발급받은 appkey

## TODO
- [ ] Make member functions `.Sell`
- [ ] Make member functions `.Buy`
- [x] `.GetDomesticAccountBalance`
- [x] Make client struct `kinvest.Client`
- [x] Convert postman env. to Go code manually
- [x] Convert openapi yaml to Go code

## Reference
- [한국투자 OpenAPI](https://apiportal.koreainvestment.com/apiservice) - API문서