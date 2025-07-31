package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/goccy/go-yaml"
	kinvest "github.com/suapapa/go_kinvest"
)

var (
	anualFiscal = false
)

func main() {
	kc, err := kinvest.NewClient(nil)
	if err != nil {
		panic(err)
	}

	var ticker string
	if len(os.Args) > 1 {
		ticker = os.Args[1]
	} else {
		log.Println("ticker is not provided, using default ticker: 005380, 현대자동차")
		ticker = "005380" // 현대자동차
	}

	ctx := context.Background()

	fmt.Println("# 주식현재가시세")
	dp, err := kc.GetDomesticInquirePrice(ctx, ticker)
	if err != nil {
		panic(err)
	}
	printInYAML(dp)
	fmt.Println()

	fmt.Println("대차대조표:")
	bs, err := kc.GetDomesticFinanceBalanceSheet(ctx, ticker, anualFiscal)
	if err != nil {
		panic(err)
	}
	printInYAML(bs[:2])
	fmt.Println()

	fmt.Println("손익계산서:")
	is, err := kc.GetDomesticFinanceIncomeStatement(ctx, ticker, anualFiscal)
	if err != nil {
		panic(err)
	}
	printInYAML(is[:2])
	fmt.Println()

	fmt.Println("재무비율:")
	fr, err := kc.GetDomesticFinanceFinancialRatio(ctx, ticker, anualFiscal)
	if err != nil {
		panic(err)
	}
	printInYAML(fr[:2])
	fmt.Println()

	fmt.Println("수익성비율:")
	pr, err := kc.GetDomesticFinanceProfitRatio(ctx, ticker, anualFiscal)
	if err != nil {
		panic(err)
	}
	printInYAML(pr[:2])
	fmt.Println()

	fmt.Println("안정성비율:")
	sr, err := kc.GetDomesticFinanceStabilityRatio(ctx, ticker, anualFiscal)
	if err != nil {
		panic(err)
	}
	printInYAML(sr[:2])
	fmt.Println()

	fmt.Println("성장성비율:")
	gr, err := kc.GetDomesticFinanceGrowthRatio(ctx, ticker, anualFiscal)
	if err != nil {
		panic(err)
	}
	printInYAML(gr[:2])
	fmt.Println()

}

func printInYAML(v any) {
	y, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(y))
}
