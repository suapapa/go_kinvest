package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
	kinvest "github.com/suapapa/go_kinvest"
)

var (
	anualFiscal = false
	stacCnt     = 1
)

// FinanceReport holds all financial data for a stock.
type FinanceReport struct {
	ItemInfo        *kinvest.ItemInfo                         `yaml:"종목정보"`
	InquirePrice    *kinvest.DomesticInquirePrice             `yaml:",inline"`
	BalanceSheet    []*kinvest.DomesticFinanceBalanceSheet    `yaml:"대차대조표"`
	IncomeStatement []*kinvest.DomesticFinanceIncomeStatement `yaml:"손익계산서"`
	FinancialRatio  []*kinvest.DomesticFinanceFinancialRatio  `yaml:"재무비율"`
	ProfitRatio     []*kinvest.DomesticFinanceProfitRatio     `yaml:"수익성비율"`
	StabilityRatio  []*kinvest.DomesticFinanceStabilityRatio  `yaml:"안정성비율"`
	GrowthRatio     []*kinvest.DomesticFinanceGrowthRatio     `yaml:"성장성비율"`
}

// CheckListItem represents a single check in the checklist.
type CheckListItem struct {
	Value interface{} `yaml:"value"`
	OkIf  string      `yaml:"ok_if"`
	IsOk  bool        `yaml:"is_ok"`
}

// StockChecklist represents the entire checklist for a stock.
type StockChecklist struct {
	Ticker     string                    `yaml:"ticker"`
	Name       string                    `yaml:"name"`
	CurrPrice  int64                     `yaml:"curr_price"`
	Market     string                    `yaml:"market"`
	CheckLists map[string]*CheckListItem `yaml:"check_lists"`
	OkCnt      int                       `yaml:"ok_cnt"`
}

func main() {
	// 1. Get ticker
	var ticker string
	if len(os.Args) > 1 {
		ticker = os.Args[1]
	} else {
		log.Println("ticker is not provided, using default ticker: 005380, 현대자동차")
		ticker = "005380" // 현대자동차
	}

	// 2. Create kinvest client
	kc, err := kinvest.NewClient(nil)
	if err != nil {
		log.Fatalf("failed to create kinvest client: %v", err)
	}

	ctx := context.Background()

	// 3. Fetch financial data
	report, err := fetchFinanceData(ctx, kc, ticker)
	if err != nil {
		log.Fatalf("failed to fetch finance data: %v", err)
	}

	// 4. Write finance report to YAML file
	reportBytes, err := yaml.Marshal(report)
	if err != nil {
		log.Fatalf("failed to marshal finance report: %v", err)
	}
	reportFilename := fmt.Sprintf("%s_finance_report.yaml", ticker)
	if err := os.WriteFile(reportFilename, reportBytes, 0644); err != nil {
		log.Fatalf("failed to write finance report to %s: %v", reportFilename, err)
	}
	fmt.Printf("Successfully generated %s\n", reportFilename)

	// 5. Generate checklist
	checklist, err := generateChecklist(ticker, report)
	if err != nil {
		log.Fatalf("failed to generate checklist: %v", err)
	}

	// 6. Write checklist to YAML file
	checklistBytes, err := yaml.Marshal(checklist)
	if err != nil {
		log.Fatalf("failed to marshal checklist: %v", err)
	}
	checklistFilename := fmt.Sprintf("%s_checklist.yaml", ticker)
	if err := os.WriteFile(checklistFilename, checklistBytes, 0644); err != nil {
		log.Fatalf("failed to write checklist to %s: %v", checklistFilename, err)
	}
	fmt.Printf("Successfully generated %s\n", checklistFilename)
}

func fetchFinanceData(ctx context.Context, kc *kinvest.Client, ticker string) (*FinanceReport, error) {
	report := &FinanceReport{}
	var err error

	fmt.Println("Fetching 종목정보...")
	report.ItemInfo, err = kc.GetDomesticItemInfo(ctx, ticker)
	if err != nil {
		return nil, fmt.Errorf("GetDomesticItemInfo failed: %w", err)
	}

	fmt.Println("Fetching 주식현재가시세...")
	report.InquirePrice, err = kc.GetDomesticInquirePrice(ctx, ticker)
	if err != nil {
		return nil, fmt.Errorf("GetDomesticInquirePrice failed: %w", err)
	}

	fmt.Println("Fetching 대차대조표...")
	bs, err := kc.GetDomesticFinanceBalanceSheet(ctx, ticker, anualFiscal)
	if err != nil {
		return nil, fmt.Errorf("GetDomesticFinanceBalanceSheet failed: %w", err)
	}
	if len(bs) > 2 {
		report.BalanceSheet = bs[:stacCnt]
	} else {
		report.BalanceSheet = bs
	}

	fmt.Println("Fetching 손익계산서...")
	is, err := kc.GetDomesticFinanceIncomeStatement(ctx, ticker, anualFiscal)
	if err != nil {
		return nil, fmt.Errorf("GetDomesticFinanceIncomeStatement failed: %w", err)
	}
	if len(is) > 2 {
		report.IncomeStatement = is[:stacCnt]
	} else {
		report.IncomeStatement = is
	}

	fmt.Println("Fetching 재무비율...")
	fr, err := kc.GetDomesticFinanceFinancialRatio(ctx, ticker, anualFiscal)
	if err != nil {
		return nil, fmt.Errorf("GetDomesticFinanceFinancialRatio failed: %w", err)
	}
	if len(fr) > 2 {
		report.FinancialRatio = fr[:stacCnt]
	} else {
		report.FinancialRatio = fr
	}

	fmt.Println("Fetching 수익성비율...")
	pr, err := kc.GetDomesticFinanceProfitRatio(ctx, ticker, anualFiscal)
	if err != nil {
		return nil, fmt.Errorf("GetDomesticFinanceProfitRatio failed: %w", err)
	}
	if len(pr) > 2 {
		report.ProfitRatio = pr[:stacCnt]
	} else {
		report.ProfitRatio = pr
	}

	fmt.Println("Fetching 안정성비율...")
	sr, err := kc.GetDomesticFinanceStabilityRatio(ctx, ticker, anualFiscal)
	if err != nil {
		return nil, fmt.Errorf("GetDomesticFinanceStabilityRatio failed: %w", err)
	}
	if len(sr) > 2 {
		report.StabilityRatio = sr[:stacCnt]
	} else {
		report.StabilityRatio = sr
	}

	fmt.Println("Fetching 성장성비율...")
	gr, err := kc.GetDomesticFinanceGrowthRatio(ctx, ticker, anualFiscal)
	if err != nil {
		return nil, fmt.Errorf("GetDomesticFinanceGrowthRatio failed: %w", err)
	}
	if len(gr) > 2 {
		report.GrowthRatio = gr[:stacCnt]
	} else {
		report.GrowthRatio = gr
	}

	return report, nil
}

func generateChecklist(ticker string, report *FinanceReport) (*StockChecklist, error) {
	// Load the template
	templateBytes, err := os.ReadFile("stock_checklist.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read stock_checklist.yaml: %w", err)
	}

	var checklist StockChecklist
	if err := yaml.Unmarshal(templateBytes, &checklist); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stock_checklist.yaml: %w", err)
	}

	// Fill in the values
	checklist.Ticker = ticker
	checklist.Name = report.ItemInfo.PrdtName
	currPrice, _ := strconv.ParseInt(report.InquirePrice.StckPrpr, 10, 64)
	checklist.CurrPrice = currPrice
	checklist.Market = report.InquirePrice.RprsMrktKorName

	// Helper function to safely parse float
	parseFloat := func(s string) float64 {
		f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
		return f
	}

	// Fill check list items
	// 시가총액
	mktCap, _ := strconv.ParseInt(report.InquirePrice.HtsAvls, 10, 64)
	mktCap *= 100_000_000 // 억 단위
	checklist.CheckLists["시가총액"].Value = mktCap

	// 증거금률
	checklist.CheckLists["증거금률"].Value = parseFloat(report.InquirePrice.MargRate)

	// 52주일최저가대비현재가대비
	checklist.CheckLists["52주일최저가대비현재가대비"].Value = parseFloat(report.InquirePrice.W52LwprVrssPrprCtrt)

	// 52주일최고가대비현재가대비
	checklist.CheckLists["52주일최고가대비현재가대비"].Value = parseFloat(report.InquirePrice.W52HgprVrssPrprCtrt)

	// PER
	checklist.CheckLists["PER"].Value = parseFloat(report.InquirePrice.Per)

	// PBR
	checklist.CheckLists["PBR"].Value = parseFloat(report.InquirePrice.Pbr)

	// PSR
	sps := parseFloat(report.FinancialRatio[0].Sps)
	var psr float64
	if sps > 0 {
		psr = float64(currPrice) / sps
	}
	checklist.CheckLists["PSR"].Value = psr

	// PCR - Data not available in the report, leave as 0
	checklist.CheckLists["PCR"].Value = 0.0

	// PEG
	per := parseFloat(report.InquirePrice.Per)
	niGrowthRate := parseFloat(report.FinancialRatio[0].NtinInrt)
	var peg float64
	if niGrowthRate > 0 {
		peg = per / niGrowthRate
	}
	checklist.CheckLists["PEG"].Value = peg

	// ROE
	checklist.CheckLists["ROE"].Value = parseFloat(report.FinancialRatio[0].RoeVal)

	// ROA
	checklist.CheckLists["ROA"].Value = parseFloat(report.ProfitRatio[0].CptlNtinRate)

	// 영업이익률
	sales := parseFloat(report.IncomeStatement[0].SaleAccount)
	op := parseFloat(report.IncomeStatement[0].BsopPrti)
	var opMargin float64
	if sales > 0 {
		opMargin = (op / sales) * 100
	}
	checklist.CheckLists["영업이익률"].Value = opMargin

	// 순이익률
	checklist.CheckLists["순이익률"].Value = parseFloat(report.ProfitRatio[0].SaleNtinRate)

	// 매출액증가율
	checklist.CheckLists["매출액증가율"].Value = parseFloat(report.FinancialRatio[0].Grs)

	// 순이익증가율
	checklist.CheckLists["순이익증가율"].Value = parseFloat(report.FinancialRatio[0].NtinInrt)

	// 부채비율
	checklist.CheckLists["부채비율"].Value = parseFloat(report.FinancialRatio[0].LbltRate)

	// 유동비율
	currentAssets := parseFloat(report.BalanceSheet[0].Cras)
	currentLiabilities := parseFloat(report.BalanceSheet[0].FlowLblt)
	var currentRatio float64
	if currentLiabilities > 0 {
		currentRatio = (currentAssets / currentLiabilities) * 100
	}
	checklist.CheckLists["유동비율"].Value = currentRatio

	// 유보율
	checklist.CheckLists["유보율"].Value = parseFloat(report.FinancialRatio[0].RsrvRate)

	// 현금배당수익률 - Data not available, leave as 0
	checklist.CheckLists["현금배당수익률"].Value = 0.0

	// 외국인지분율
	checklist.CheckLists["외국인지분율"].Value = parseFloat(report.InquirePrice.HtsFrgnEhrt)

	// Evaluate is_ok and count
	okCnt := 0
	for _, item := range checklist.CheckLists {
		item.IsOk = evaluateOkIf(item.Value, item.OkIf)
		if item.IsOk {
			okCnt++
		}
	}
	checklist.OkCnt = okCnt

	return &checklist, nil
}

func evaluateOkIf(value interface{}, condition string) bool {
	// cleanup condition string
	condition = strings.ReplaceAll(condition, "_", "")
	condition = strings.ReplaceAll(condition, "%", "")

	var operator string
	var limitStr string

	if strings.HasPrefix(condition, ">=") {
		operator = ">="
		limitStr = strings.TrimSpace(condition[2:])
	} else if strings.HasPrefix(condition, "<=") {
		operator = "<="
		limitStr = strings.TrimSpace(condition[2:])
	} else if strings.HasPrefix(condition, ">") {
		operator = ">"
		limitStr = strings.TrimSpace(condition[1:])
	} else if strings.HasPrefix(condition, "<") {
		operator = "<"
		limitStr = strings.TrimSpace(condition[1:])
	} else {
		return false // unsupported condition
	}

	limit, err := strconv.ParseFloat(limitStr, 64)
	if err != nil {
		return false // failed to parse limit
	}

	var valFloat float64
	switch v := value.(type) {
	case float64:
		valFloat = v
	case float32:
		valFloat = float64(v)
	case int64:
		valFloat = float64(v)
	case int32:
		valFloat = float64(v)
	case int:
		valFloat = float64(v)
	default:
		return false // unsupported value type
	}

	switch operator {
	case ">=":
		return valFloat >= limit
	case "<=":
		return valFloat <= limit
	case ">":
		return valFloat > limit
	case "<":
		return valFloat < limit
	}

	return false
}
