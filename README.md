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
- `KINVEST_CANO` : 계좌번호
- `KINVEST_APPKEY` : 한국투자증권 개발자센터에서 발급받은 appkey
- `KINVEST_APPSECRET` :  한국투자증권 개발자센터에서 발급받은 appkey
- `KINVEST_MAC` : 한국투자증권 API 서버와 통신하는데 사용되는 네트워크 인터페이스의 MAC 주소

## TODO
- [ ] Make member functions `.Sell`
- [ ] Make member functions `.Buy`
- [x] `.GetDomesticAccountBalance`
- [x] Make client struct `kinvest.Client`
- [x] Convert postman env. to Go code manually
- [x] Convert openapi yaml to Go code

## Reference
- [한국투자 OpenAPI](https://apiportal.koreainvestment.com/apiservice) - API문서