package main

import (
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

	res, err := kc.BuyDomesticStock(hmgCode, qty, nil)
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
