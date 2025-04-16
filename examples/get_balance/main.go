package main

import (
	"fmt"

	kinvest "github.com/suapapa/go_kinvest"
)

func main() {
	kc, err := kinvest.NewClient(nil)
	if err != nil {
		panic(err)
	}

	bal, err := kc.GetDomesticAccountBalance()
	if err != nil {
		panic(err)
	}

	fmt.Printf("balance: %v", bal)
}
