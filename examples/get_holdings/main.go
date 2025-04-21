package main

import (
	"context"
	"log"

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

	for {
		res, err = res.GetNext(context.Background())
		if err != nil {
			log.Println("error:", err)
			break
		}
		if res == nil {
			log.Println("no more data")
			break
		}
		y, err := yaml.Marshal(res)
		if err != nil {
			panic(err)
		}
		println(string(y))
	}

}
