# kinvest : A Go package for 한국투자증권 client

**!! WIP !!**

Install:
```sh
go get github.com/suapapa/go_kinvest
```

Example usage:
```go
// TBD
```

This package read following envs for the API calls:
- `KINVEST_CANO` : 계좌번호
- `KINVEST_APPKEY`
- `KINVEST_APPSECRET`
- `KINVEST_TOKEN`
- `KINVEST_HASH`

## TODO
- [ ] Make essensial member functions `.Balance`, `.Buy`, `Sell`
- [ ] Make client struct `kinvest.Client`
- [x] Convert postman env. to Go code manually

## Reference
- [한국투자 OpenAPI](https://apiportal.koreainvestment.com/apiservice) - API문서