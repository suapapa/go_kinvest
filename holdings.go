package kinvest

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/suapapa/go_kinvest/internal/oapi"
)

func (c *Client) GetDomesticHoldings(ctx context.Context, opt *GetDomesticHoldingsOptions) (*GetDomesticHoldingsResult, error) {
	if opt == nil {
		var err error
		opt, err = NewGetDomesticHoldingsOptions("기본", "종목별")
		if err != nil {
			return nil, fmt.Errorf("create get domestic holdings option failed: %w", err)
		}
	}

	cano, acntprdtcd, err := parseAccount(c.account)
	if err != nil {
		return nil, fmt.Errorf("parse account failed: %w", err)
	}

	resp, err := c.oc.GetUapiDomesticStockV1TradingInquireBalance(
		ctx,
		&oapi.GetUapiDomesticStockV1TradingInquireBalanceParams{
			CANO:              cano,
			ACNTPRDTCD:        acntprdtcd,
			AFHRFLPRYN:        opt.tradingSessionTypeCode(),          // 시간외단일가, 거래소여부
			OFLYN:             ptr(""),                               // -
			INQRDVSN:          opt.queryTypeCode(),                   // 01: 대출일별, 02: 종목별
			UNPRDVSN:          ptr(1),                                // -
			FUNDSTTLICLDYN:    ptr(toStr(opt.IncludeFundSettlement)), // 펀드 평가금액 포함
			FNCGAMTAUTORDPTYN: ptr("N"),                              // -
			PRCSDVSN:          opt.includePrevTradingCode(),          // 00: 전일매매포함, 01: 전일매매미포함
			CTXAREAFK100:      ptr(opt.CtxAreaFK),                    // 이전 조회 CTX_AREA_MK100
			CTXAREANK100:      ptr(opt.CtxAreaNK),                    // 이전 조회 CTX_AREA_MK100
			TrId:              ptr("TTTC8434R"),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	return newGetDomesticHoldingsResult(c, opt, mustUnmarshalJsonBody(resp.Body))
}

type GetDomesticHoldingsOptions struct {
	TradingSessionType    string `yaml:"거래세션유형,omitempty"` // 기본, 시간외단일가, NXT정규장
	QueryType             string `yaml:"조회구분,omitempty"`   // 대출일별, 종목별
	IncludeFundSettlement bool   `yaml:"편드결제분포함여부,omitempty"`
	IncludePrevTrading    bool   `yaml:"전일매매포함,omitempty"`

	CtxAreaFK, CtxAreaNK string
}

func NewGetDomesticHoldingsOptions(tradingSessionType, queryType string) (*GetDomesticHoldingsOptions, error) {
	if _, ok := tradingSessionTypeCode[tradingSessionType]; !ok {
		var tradingSessionTypes []string
		for k := range tradingSessionTypeCode {
			tradingSessionTypes = append(tradingSessionTypes, k)
		}
		return nil, fmt.Errorf("invalid trading session type: %s, set one of the following: %s", tradingSessionType, strings.Join(tradingSessionTypes, ", "))
	}

	if _, ok := queryTypeCode[queryType]; !ok {
		var queryTypes []string
		for k := range queryTypeCode {
			queryTypes = append(queryTypes, k)
		}
		return nil, fmt.Errorf("invalid query type: %s, set one of the following: %s", queryType, strings.Join(queryTypes, ", "))
	}

	return &GetDomesticHoldingsOptions{
		TradingSessionType: tradingSessionType,
		QueryType:          queryType,
	}, nil
}

func (o *GetDomesticHoldingsOptions) tradingSessionTypeCode() *string {
	if code, ok := tradingSessionTypeCode[o.TradingSessionType]; ok {
		return code
	}
	return tradingSessionTypeCode["기본"]
}

var tradingSessionTypeCode = map[string]*string{
	"기본":     ptr("N"),
	"시간외단일가": ptr("Y"),
	"NXT정규장": ptr("X"),
}

func (o *GetDomesticHoldingsOptions) queryTypeCode() *int {
	if code, ok := queryTypeCode[o.QueryType]; ok {
		return code
	}
	return queryTypeCode["종목별"]
}

func (o *GetDomesticHoldingsOptions) includePrevTradingCode() *int {
	if o.IncludePrevTrading {
		return prevTradingCode[true]
	}
	return prevTradingCode[false]
}

var queryTypeCode = map[string]*int{
	"대출일별": ptr(01),
	"종목별":  ptr(02),
}

var prevTradingCode = map[bool]*int{
	true:  ptr(00), // 전일매매포함"
	false: ptr(01), // 전일매매미포함
}

func newGetDomesticHoldingsResult(c *Client, opt *GetDomesticHoldingsOptions, data map[string]any) (*GetDomesticHoldingsResult, error) {
	if data == nil {
		return nil, fmt.Errorf("response is nil")
	}
	if data["rt_cd"] != "0" {
		return nil, fmt.Errorf("response error: %s (%s)", data["msg1"], data["msg_cd"])
	}

	ret := &GetDomesticHoldingsResult{
		c:         c,
		opt:       opt,
		ctxAreaFK: data["ctx_area_fk100"].(string),
		ctxAreaNK: data["ctx_area_nk100"].(string),
	}
	if output, ok := data["output1"].(map[string]any); ok {
		s := &Stock{
			Code:              output["pdno"].(string),
			Name:              output["prdt_name"].(string),
			OrderType:         output["trad_dvsn_name"].(string), // 매수매도구분
			PrevBuyQty:        toInt(output["bfdy_buy_qty"]),
			PrevSellQty:       toInt(output["bfdy_sll_qty"]),
			TodayBuyQty:       toInt(output["thdt_buyqty"]),
			TodaySellQty:      toInt(output["thdt_sll_qty"]),
			HoldingQty:        toInt(output["hldg_qty"]),
			OrdPossibleQty:    toInt(output["ord_psbl_qty"]),
			PurchaseAvgPrice:  toFloat(output["pchs_avg_pric"]),
			PurchaseAmount:    toInt(output["pchs_amt"]),
			CurrPrice:         toInt(output["prpr"]),
			EvalAmount:        toInt(output["evlu_amt"]),
			EvalProfitAmount:  toInt(output["evlu_pfls_amt"]),
			EvalProfitRate:    toFloat(output["evlu_pfls_rt"]),
			LoanDate:          toTime(output["loan_dt"]),
			LoanAmount:        toInt(output["loan_amt"]),
			ShortSellAmount:   toInt(output["stln_slng_chgs"]),
			ExpiredDate:       toTime(output["expd_dt"]),
			ChangeRate:        toFloat(output["fltt_rt"]),
			PriceDiffFromPrev: toInt(output["bfdy_cprs_icdc"]),
			MarginRate:        output["item_mgna_rt_name"].(string),
			GuaranteeRate:     output["grta_rt_name"].(string),
			SubstitutePrice:   toInt(output["sbst_pric"]),
			LoanPrice:         toFloat(output["stck_loan_unpr"]),
		}
		ret.Holdings = append(ret.Holdings, s)
	}
	if output, ok := data["output2"].(map[string]any); ok {
		b := &Balance{
			TotalDeposit:                  toInt(output["dnca_tot_amt"]),
			NextSettlementAmount:          toInt(output["nxdy_excc_amt"]),
			TempSettlementAmount:          toInt(output["prvs_rcdl_excc_amt"]),
			CMAValuationAmount:            toInt(output["cma_evlu_amt"]),
			PrevBuyAmount:                 toInt(output["bfdy_buy_amt"]),
			TodayBuyAmount:                toInt(output["thdt_buy_amt"]),
			NextAutoRepaymentAmount:       toInt(output["nxdy_auto_rdpt_amt"]),
			PrevSellAmount:                toInt(output["bfdy_sll_amt"]),
			TodaySellAmount:               toInt(output["thdt_sll_amt"]),
			D2AutoRepaymentAmount:         toInt(output["d2_auto_rdpt_amt"]),
			PrevFeeAmount:                 toInt(output["bfdy_tlex_amt"]),
			TodayFeeAmount:                toInt(output["thdt_tlex_amt"]),
			TotalLoanAmount:               toInt(output["tot_loan_amt"]),
			SecuritiesValuationAmount:     toInt(output["scts_evlu_amt"]),
			TotalValuationAmount:          toInt(output["tot_evlu_amt"]),
			NetAssetAmount:                toInt(output["nass_amt"]),
			IsAutoRepaymentForLoan:        output["fncg_gld_auto_rdpt_yn"].(string) == "Y",
			TotalPurchaseAmount:           toInt(output["pchs_amt_smtl_amt"]),
			TotalValuationSum:             toInt(output["evlu_amt_smtl_amt"]),
			TotalUnrealizedPnL:            toInt(output["evlu_pfls_smtl_amt"]),
			TotalShortSellProceeds:        toInt(output["tot_stln_slng_chgs"]),
			PrevTotalAssetValuationAmount: toInt(output["bfdy_tot_asst_evlu_amt"]),
			AssetChangeAmount:             toInt(output["asst_icdc_amt"]),
			AssetChangeReturnRate:         toFloat(output["asst_icdc_erng_rt"]),
		}
		ret.Balances = append(ret.Balances, b)
	}

	if len(ret.Holdings) == 0 && len(ret.Balances) == 0 {
		return nil, fmt.Errorf("response output is empty")
	}
	return ret, nil
}

type GetDomesticHoldingsResult struct {
	c                    *Client
	opt                  *GetDomesticHoldingsOptions
	ctxAreaFK, ctxAreaNK string
	Holdings             []*Stock   `yaml:"holdings,omitempty"`
	Balances             []*Balance `yaml:"balances,omitempty"`
}

func (r *GetDomesticHoldingsResult) GetNext(ctx context.Context) (*GetDomesticHoldingsResult, error) {
	if r.ctxAreaFK == "" || r.ctxAreaNK == "" {
		return nil, fmt.Errorf("no next page")
	}
	opt := *r.opt
	opt.CtxAreaFK = r.ctxAreaFK
	opt.CtxAreaNK = r.ctxAreaNK

	return r.c.GetDomesticHoldings(ctx, &opt)
}

type Stock struct {
	Code              string    `yaml:"종목번호"`
	Name              string    `yaml:"종목명"`
	OrderType         string    `yaml:"매매구분,omitempty"`
	PrevBuyQty        int       `yaml:"전 일 매수수량,omitempty"`
	PrevSellQty       int       `yaml:"전 일 매도수량,omitempty"`
	TodayBuyQty       int       `yaml:"금 일 매수수량,omitempty"`
	TodaySellQty      int       `yaml:"금 일 매도수량,omitempty"`
	HoldingQty        int       `yaml:"보유수량,omitempty"`
	OrdPossibleQty    int       `yaml:"주문가능수량,omitempty"`
	PurchaseAvgPrice  float64   `yaml:"매입평균가격,omitempty"` // 매입금액 / 보유수량
	PurchaseAmount    int       `yaml:"매입금액,omitempty"`
	CurrPrice         int       `yaml:"현재가,omitempty"`
	EvalAmount        int       `yaml:"평가금액,omitempty"`
	EvalProfitAmount  int       `yaml:"평가손익금액,omitempty"` // 평가금액 - 매입금액
	EvalProfitRate    float64   `yaml:"평가손익률,omitempty"`
	LoanDate          time.Time `yaml:"대출일자,omitempty"`
	LoanAmount        int       `yaml:"대출금액,omitempty"`
	ShortSellAmount   int       `yaml:"대주매각대금,omitempty"` // 공매도
	ExpiredDate       time.Time `yaml:"만기일자,omitempty"`
	ChangeRate        float64   `yaml:"등락률,omitempty"`
	PriceDiffFromPrev int       `yaml:"전일대비증감,omitempty"`
	MarginRate        string    `yaml:"종목증거금율명,omitempty"`
	GuaranteeRate     string    `yaml:"보증금율명,omitempty"`
	SubstitutePrice   int       `yaml:"대용가격,omitempty"` // 증권매매의 위탁보증금으로서 현금 대신에 사용되는 유가증권 가격
	LoanPrice         float64   `yaml:"주식대출단가,omitempty"`
}

type Balance struct {
	TotalDeposit                  int     `yaml:"예수금 총액,omitempty"`
	NextSettlementAmount          int     `yaml:"익일정산금액,omitempty"`  // D+1 예수금
	TempSettlementAmount          int     `yaml:"가수도정산금액,omitempty"` // D+2 예수금
	CMAValuationAmount            int     `yaml:"CMA평가금액,omitempty"`
	PrevBuyAmount                 int     `yaml:"전일매수금액,omitempty"`
	TodayBuyAmount                int     `yaml:"금일매수금액,omitempty"`
	NextAutoRepaymentAmount       int     `yaml:"익일자동상환금액,omitempty"`
	PrevSellAmount                int     `yaml:"전일매도금액,omitempty"`
	TodaySellAmount               int     `yaml:"금일매도금액,omitempty"`
	D2AutoRepaymentAmount         int     `yaml:"D+2자동상환금액,omitempty"`
	PrevFeeAmount                 int     `yaml:"전일제비용금액,omitempty"`
	TodayFeeAmount                int     `yaml:"금일제비용금액,omitempty"`
	TotalLoanAmount               int     `yaml:"총대출금액,omitempty"`
	SecuritiesValuationAmount     int     `yaml:"유가평가금액,omitempty"`
	TotalValuationAmount          int     `yaml:"총평가금액,omitempty"` // 유가증권 평가금액 합계금액 + D+2 예수금
	NetAssetAmount                int     `yaml:"순자산금액,omitempty"`
	IsAutoRepaymentForLoan        bool    `yaml:"융자금자동상환여부,omitempty"` //보유현금에 대한 융자금만 차감여부
	TotalPurchaseAmount           int     `yaml:"매입금액합계금액,omitempty"`
	TotalValuationSum             int     `yaml:"평가금액합계금액,omitempty"`
	TotalUnrealizedPnL            int     `yaml:"평가손익합계금액,omitempty"`
	TotalShortSellProceeds        int     `yaml:"총대주매각대금,omitempty"`
	PrevTotalAssetValuationAmount int     `yaml:"전일총자산평가금액,omitempty"`
	AssetChangeAmount             int     `yaml:"자산증감액,omitempty"`
	AssetChangeReturnRate         float64 `yaml:"자산증감수익율,omitempty"`
}
