package main

import (
	"context"

	"github.com/goccy/go-yaml"
	kinvest "github.com/suapapa/go_kinvest"
)

func main() {
	kc, err := kinvest.NewClient(nil)
	if err != nil {
		panic(err)
	}

	res, err := kc.GetDomesticHoldings(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	y, err := yaml.Marshal(res)
	if err != nil {
		panic(err)
	}
	println(string(y))
}
