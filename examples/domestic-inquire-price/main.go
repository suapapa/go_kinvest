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
	domesticInquirePrice, err := kc.GetDomesticInquirePrice(context.Background(), hmgCode)
	if err != nil {
		panic(err)
	}
	y, err := yaml.Marshal(domesticInquirePrice)
	if err != nil {
		panic(err)
	}
	fmt.Println("현대자동차 현재가 정보:")
	fmt.Println(string(y))
}
