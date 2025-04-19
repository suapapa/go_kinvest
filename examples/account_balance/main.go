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

	bal, err := kc.GetDomesticAccountBalance(context.Background())
	if err != nil {
		panic(err)
	}

	y, err := yaml.Marshal(bal)
	if err != nil {
		panic(err)
	}
	println(string(y))
}
