**This package is auto generated. DO NOT EDIT MANUALLY!**

1. Download Postman collection, `실전계좌_POSTMAN_샘플코드_v2.6.json` from [한국투자증권 Open Trading API](https://github.com/koreainvestment/open-trading-api) repository
2. Convert the Postman collection to OpenAPI schema using [P20](https://p2o.defcon007.com/)
3. Gen Go code with `oapi-codegen`

```sh
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
oapi-codegen --package=kinvest --generate types -o types.go openapi/kinvest_prod.yaml
oapi-codegen --package=kinvest --generate client -o client.go openapi/kinvest_prod.yaml
```