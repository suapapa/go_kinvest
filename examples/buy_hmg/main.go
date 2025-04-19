package main

import (
	"context"
	"fmt"

	"github.com/goccy/go-yaml"
	kinvest "github.com/suapapa/go_kinvest"
)

func main() {
	kc, err := kinvest.NewClient(nil)
	if err != nil {
		panic(err)
	}

	hmgCode := "005380" // 현대자동차
	qty := 1

	fmt.Println("엔터를 누르면 '현대자동차'를 1주 매수합니다. 중지하려면 Ctrl+C를 누르세요.")
	fmt.Scanln()

	res, err := kc.BuyDomesticStock(context.Background(), hmgCode, qty, nil)
	if err != nil {
		// for unwrappedErr := err; unwrappedErr != nil; unwrappedErr = errors.Unwrap(unwrappedErr) {
		// 	fmt.Println("Error:", unwrappedErr)
		// }
		panic(err)
	}
	y, err := yaml.Marshal(res)
	if err != nil {
		panic(err)
	}
	println(string(y))
}
