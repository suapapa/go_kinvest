package kinvest

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/suapapa/go_kinvest/internal/oapi"
)

// SellDomesticStock sells domestic(KRX) stock.
func (c *Client) SellDomesticStock(ctx context.Context, code string, qty int, opt *OrderDomesticStockOptions) (*OrderResult, error) {
	if len(code) != 6 {
		return nil, fmt.Errorf("invalid item no: %s", code)
	}

	if qty <= 0 {
		return nil, fmt.Errorf("invalid qty: %d", qty)
	}

	if opt == nil {
		var err error
		opt, err = NewSellOrderDomesticStockOptions("시장가", "일반매도", 0)
		if err != nil {
			return nil, fmt.Errorf("create buy option failed: %w", err)
		}
	}

	cano, acntprdtcd, err := parseAccount(c.account)
	if err != nil {
		return nil, fmt.Errorf("parse account failed: %w", err)
	}

	res, err := c.oc.PostUapiDomesticStockV1TradingOrderCash(
		ctx,
		&oapi.PostUapiDomesticStockV1TradingOrderCashParams{
			TrId: ptr("TTTC0801U"),
		},
		oapi.PostUapiDomesticStockV1TradingOrderCashJSONRequestBody{
			"CANO":         *cano,
			"ACNT_PRDT_CD": fmt.Sprintf("%d", *acntprdtcd),
			"PDNO":         code,                         // 종목코드
			"ORD_DVSN":     opt.getDVSN(),                // 주문구분
			"ORD_QTY":      fmt.Sprintf("%d", qty),       // 주문수량
			"ORD_UNPR":     fmt.Sprintf("%d", opt.Price), // 주문단가 0: 시장가
			"SLL_TYPE":     opt.getSellTypeCode(),        // 매도유형
		},
	)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	return newOrderResult(mustUnmarshalJsonBody(res.Body))
}

// BuyDomesticStock buys domestic(KRX) stock.
func (c *Client) BuyDomesticStock(ctx context.Context, code string, qty int, opt *OrderDomesticStockOptions) (*OrderResult, error) {
	if len(code) != 6 {
		return nil, fmt.Errorf("invalid item no: %s", code)
	}

	if qty <= 0 {
		return nil, fmt.Errorf("invalid qty: %d", qty)
	}

	if opt == nil {
		var err error
		opt, err = NewBuyOrderDomesticStockOptions("시장가", 0)
		if err != nil {
			return nil, fmt.Errorf("create buy option failed: %w", err)
		}
	}

	cano, acntprdtcd, err := parseAccount(c.account)
	if err != nil {
		return nil, fmt.Errorf("parse account failed: %w", err)
	}

	res, err := c.oc.PostUapiDomesticStockV1TradingOrderCash(
		ctx,
		&oapi.PostUapiDomesticStockV1TradingOrderCashParams{
			TrId: ptr("TTTC0802U"),
		},
		oapi.PostUapiDomesticStockV1TradingOrderCashJSONRequestBody{
			"CANO":         *cano,
			"ACNT_PRDT_CD": fmt.Sprintf("%d", *acntprdtcd),
			"PDNO":         code,                         // 종목코드
			"ORD_DVSN":     opt.getDVSN(),                // 주문구분
			"ORD_QTY":      fmt.Sprintf("%d", qty),       // 주문수량
			"ORD_UNPR":     fmt.Sprintf("%d", opt.Price), // 주문단가 0: 시장가
		},
	)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	return newOrderResult(mustUnmarshalJsonBody(res.Body))
}

// OrderDomesticStockOptions is the options for domestic stock order.
// Type is one of following:
//
//	지정가, 시장가, 조건부지정가, 최유리지정가, 최우선지정가, 장전 시간외, 장후 시간외, 시간외 단일가,
//	경매매, 자기주식, 자기주식S-Option, 자기주식금전신탁, IOC지정가, IOC시장가, FOK시장가, IOC최유리, FOK최유리,
//	장중대량, 장중바스켓, 장개시전 시간외대량, 장개시전 시간외바스켓, 장개시전 금전신탁자사주, 장개시전 자기주식,
//	시간외대량, 시간외자사주신탁, 시간외대량자기주식, 바스켓, 중간가, 스톱지정가, 중간가IOC, 중간가FOK
//
// SellType will be ignored when buy and should be set one of following when sell:
//
//	일반매도, 임의매도, 대차매도
type OrderDomesticStockOptions struct {
	Type     string // 주문구분
	Price    int    // 주문단가
	SellType string // 매도구분 : 일반매도, 임의매도, 대차매도
}

// NewBuyOrderDomesticStockOptions creates a new BuyDomesticStock.
func NewBuyOrderDomesticStockOptions(orderType string, orderPrice int) (*OrderDomesticStockOptions, error) {
	opt, err := newOrderDomesticStockOptions(orderType, "", orderPrice)
	if err != nil {
		return nil, fmt.Errorf("create buy option failed: %w", err)
	}
	return opt, nil
}

// NewSellOrderDomesticStockOptions creates a new SellDomesticStock.
func NewSellOrderDomesticStockOptions(orderType, sellType string, orderPrice int) (*OrderDomesticStockOptions, error) {
	opt, err := newOrderDomesticStockOptions(orderType, sellType, orderPrice)
	if err != nil {
		return nil, fmt.Errorf("create sell option failed: %w", err)
	}
	return opt, nil
}

func newOrderDomesticStockOptions(orderType, sellType string, orderPrice int) (*OrderDomesticStockOptions, error) {
	if orderPrice < 0 {
		return nil, fmt.Errorf("invalid order price: %d", orderPrice)
	}

	if _, ok := dvsnCodes[orderType]; !ok {
		var orderTypes []string
		for k := range dvsnCodes {
			orderTypes = append(orderTypes, k)
		}
		return nil, fmt.Errorf("invalid order type: %s, set one of the following: %s", orderType, strings.Join(orderTypes, ", "))
	}

	if sellType != "" {
		if _, ok := sellTypeCode[sellType]; !ok {
			var sellTypeCodes []string
			for k := range sellTypeCode {
				sellTypeCodes = append(sellTypeCodes, k)
			}
			return nil, fmt.Errorf("invalid sell type: %s, set one of the following: %s", sellType, strings.Join(sellTypeCodes, ", "))
		}
	}

	return &OrderDomesticStockOptions{
		Type:     orderType,
		Price:    orderPrice,
		SellType: sellType,
	}, nil
}

func (o *OrderDomesticStockOptions) getDVSN() string {
	if code, ok := dvsnCodes[o.Type]; ok {
		return code
	} else {
		return dvsnCodes["지정가"]
	}
}

var dvsnCodes = orderType{
	"지정가":          "00",
	"시장가":          "01",
	"조건부지정가":       "02",
	"최유리지정가":       "03",
	"최우선지정가":       "04",
	"장전 시간외":       "05",
	"장후 시간외":       "06",
	"시간외 단일가":      "07",
	"경매매":          "65",
	"자기주식":         "08",
	"자기주식S-Option": "09",
	"자기주식금전신탁":     "10",
	"IOC지정가":       "11", // 즉시체결,잔량취소
	"IOC시장가":       "13", // 즉시체결,잔량취소
	"FOK시장가":       "14", // 즉시체결,잔량취소
	"IOC최유리":       "15", // 즉시체결,잔량취소
	"FOK최유리":       "16", // 즉시체결,잔량취소
	"장중대량":         "51",
	"장중바스켓":        "52",
	"장개시전 시간외대량":   "62",
	"장개시전 시간외바스켓":  "63",
	"장개시전 금전신탁자사주": "67",
	"장개시전 자기주식":    "69",
	"시간외대량":        "72",
	"시간외자사주신탁":     "77",
	"시간외대량자기주식":    "79",
	"바스켓":          "80",
	"중간가":          "21",
	"스톱지정가":        "22",
	"중간가IOC":       "23",
	"중간가FOK":       "24",
}

func (o *OrderDomesticStockOptions) getSellTypeCode() string {
	if code, ok := sellTypeCode[o.SellType]; ok {
		return code
	}
	return sellTypeCode["일반매도"]
}

var sellTypeCode = map[string]string{
	"일반매도": "01",
	"임의매도": "02",
	"대차매도": "05",
}

type orderType map[string]string

func (o orderType) OrdDvsnCode(ot string) (string, error) {
	if code, ok := o[ot]; ok {
		return code, nil
	}
	for k, v := range o {
		if strings.Contains(k, ot) {
			return v, nil
		}
	}

	return "", fmt.Errorf("order type not found: %s", ot)
}

func newOrderResult(data map[string]any) (*OrderResult, error) {
	if data == nil {
		return nil, fmt.Errorf("response is nil")
	}
	if data["rt_cd"] != "0" {
		return nil, fmt.Errorf("response error: %s (%s)", data["msg1"], data["msg_cd"])
	}

	if output, ok := data["output"].(map[string]any); ok {
		ordNo := output["ODNO"].(string)
		ordTimeStr := output["ORD_TMD"].(string)
		venue := output["KRX_FWDG_ORD_ORGNO"].(string)
		if ordNo == "" || ordTimeStr == "" || venue == "" {
			return nil, fmt.Errorf("response output is nil")
		}

		ordTime, err := hhmmssToTime(ordTimeStr)
		if err != nil {
			return nil, fmt.Errorf("convert order time failed: %w", err)
		}

		return &OrderResult{
			OrderNo:   ordNo,
			OrderedAt: ordTime,
			Venue:     venue,
		}, nil
	}

	return nil, fmt.Errorf("response output is not a map")
}

// OrderResult is the result of the order.
type OrderResult struct {
	OrderNo   string    `yaml:"주문번호"`
	OrderedAt time.Time `yaml:"주문시간"`
	Venue     string    `yaml:"거래소코드"`
}
