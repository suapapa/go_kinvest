package main

import (
	"context"
	"fmt"

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

	fmt.Printf("%s\n", bal)
}
