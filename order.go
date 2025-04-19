package kinvest

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/suapapa/go_kinvest/internal/oapi"
)

func (c *Client) SellDomesticStock(itemNo string, qty int, opt *OrderDomesticStockOptions) (*OrderResult, error) {
	if len(itemNo) != 6 {
		return nil, fmt.Errorf("invalid item no: %s", itemNo)
	}

	if qty <= 0 {
		return nil, fmt.Errorf("invalid qty: %d", qty)
	}

	if opt == nil {
		var err error
		opt, err = NewOrderDomesticStockOptions("시장가", 0)
		if err != nil {
			return nil, fmt.Errorf("create buy option failed: %w", err)
		}
	}

	cano, acntprdtcd, err := parseAccount(c.account)
	if err != nil {
		return nil, fmt.Errorf("parse account failed: %w", err)
	}

	reqBody := mustCreateJsonReader(oapi.PostUapiDomesticStockV1TradingOrderCashJSONRequestBody{
		"CANO":         cano,
		"ACNT_PRDT_CD": acntprdtcd,
		"PDNO":         itemNo,                       // 종목코드
		"ORD_DVSN":     opt.getDVSN(),                // 주문구분
		"ORD_QTY":      fmt.Sprintf("%d", qty),       // 주문수량
		"ORD_UNPR":     fmt.Sprintf("%d", opt.Price), // 주문단가 0: 시장가
		"SLL_TYPE":     opt.getSellTypeCode(),        // 매도유형
	})

	req, err := oapi.NewPostUapiDomesticStockV1TradingOrderCashRequestWithBody(
		c.oc.Server,
		&oapi.PostUapiDomesticStockV1TradingOrderCashParams{
			TrId: ptr("TTTC0801U"),
		},
		jsonContentType,
		reqBody,
	)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("request failed: %s", res.Status)
	}
	// io.Copy(os.Stdout, res.Body)
	return newOrderResult(mustUnmarshalJsonBody(res.Body))
}

func (c *Client) BuyDomesticStock(code string, qty int, opt *OrderDomesticStockOptions) (*OrderResult, error) {
	if len(code) != 6 {
		return nil, fmt.Errorf("invalid item no: %s", code)
	}

	if qty <= 0 {
		return nil, fmt.Errorf("invalid qty: %d", qty)
	}

	if opt == nil {
		var err error
		opt, err = NewOrderDomesticStockOptions("시장가", 0)
		if err != nil {
			return nil, fmt.Errorf("create buy option failed: %w", err)
		}
	}

	cano, acntprdtcd, err := parseAccount(c.account)
	if err != nil {
		return nil, fmt.Errorf("parse account failed: %w", err)
	}

	reqBody := mustCreateJsonReader(oapi.PostUapiDomesticStockV1TradingOrderCashJSONRequestBody{
		"CANO":         cano,
		"ACNT_PRDT_CD": acntprdtcd,
		"PDNO":         code,                         // 종목코드
		"ORD_DVSN":     opt.getDVSN(),                // 주문구분
		"ORD_QTY":      fmt.Sprintf("%d", qty),       // 주문수량
		"ORD_UNPR":     fmt.Sprintf("%d", opt.Price), // 주문단가 0: 시장가
	})
	req, err := oapi.NewPostUapiDomesticStockV1TradingOrderCashRequestWithBody(
		c.oc.Server,
		&oapi.PostUapiDomesticStockV1TradingOrderCashParams{
			// ContentType: &jsonContentType,
			// Authorization: c.token.Authorization(),
			// Appkey:        &c.appKey,
			// Appsecret:     &c.appSecret,
			TrId: ptr("TTTC0802U"),
		},
		jsonContentType,
		reqBody,
	)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("request failed: %s", res.Status)
	}
	// io.Copy(os.Stdout, res.Body)
	return newOrderResult(mustUnmarshalJsonBody(res.Body))
}

type OrderDomesticStockOptions struct {
	Type     string // 주문구분
	Price    int    // 주문단가
	SellType string // 매도구분 : 일반매도, 임의매도, 대차매도
}

// NewOrderDemesticStockOption creates a new BuytDomesticStockOption.
// if orderPrice is 0, orderType is somthing like the "시장가".
// orderType should be one of the following:
//
//	지정가, 시장가, 조건부지정가, 최유리지정가, 최우선지정가, 장전 시간외, 장후 시간외, 시간외 단일가,
//	경매매, 자기주식, 자기주식S-Option, 자기주식금전신탁, IOC지정가, IOC시장가, FOK시장가, IOC최유리,
//	FOK최유리, 장중대량, 장중바스켓, 장개시전 시간외대량, 장개시전 시간외바스켓, 장개시전 금전신탁자사주,
//	장개시전 자기주식, 시간외대량, 시간외자사주신탁, 시간외대량자기주식, 바스켓, 중간가, 스톱지정가,
//	중간가IOC, 중간가FOK
func NewOrderDomesticStockOptions(orderType string, orderPrice int) (*OrderDomesticStockOptions, error) {
	if _, ok := dvsnCodes[orderType]; !ok {
		var orderTypes []string
		for k := range dvsnCodes {
			orderTypes = append(orderTypes, k)
		}
		return nil, fmt.Errorf("invalid order type: %s, set one of the following: %s", orderType, strings.Join(orderTypes, ", "))
	}

	return &OrderDomesticStockOptions{
		Type:  orderType,
		Price: orderPrice,
	}, nil
}

func (o *OrderDomesticStockOptions) SetSellType(sellType string) {
	switch sellType {
	case "일반매도", "임의매도", "대차매도":
		o.SellType = sellType
	default:
		o.SellType = "일반매도"
	}
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

func (o *OrderDomesticStockOptions) getOrderPrice() (string, error) {
	if o.Price == 0 {
		return "0", nil
	}
	if o.Price < 0 {
		return "", fmt.Errorf("order price is negative: %d", o.Price)
	}
	return fmt.Sprintf("%d", o.Price), nil
}

func (o *OrderDomesticStockOptions) getSellTypeCode() string {
	sellTypeCode := map[string]string{
		"일반매도": "01",
		"임의매도": "02",
		"대차매도": "05",
	}
	if code, ok := sellTypeCode[o.SellType]; ok {
		return code
	}
	return sellTypeCode["일반매도"]
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

		ordTimeUnix, err := strconv.ParseInt(ordTimeStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid order time: %s", ordTimeStr)
		}
		ordTime := time.Unix(ordTimeUnix, 0)

		return &OrderResult{
			OrderNo:   ordNo,
			OrderedAt: ordTime,
			Venue:     venue,
		}, nil
	}

	return nil, fmt.Errorf("response output is not a map")
}

type OrderResult struct {
	OrderNo   string    `yaml:"주문번호"`
	OrderedAt time.Time `yaml:"주문시간"`
	Venue     string    `yaml:"거래소코드"`
}
