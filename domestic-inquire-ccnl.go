// 주식현재가 체결

package kinvest

import (
	"context"
	"fmt"

	"github.com/suapapa/go_kinvest/internal/oapi"
)

func (c *Client) GetDomesticInquireCcnl(ctx context.Context, code string) ([]*DomesticInquireCcnl, error) {
	if len(code) != 6 {
		return nil, fmt.Errorf("invalid item no: %s", code)
	}

	resp, err := c.oc.GetUapiDomesticStockV1QuotationsInquireCcnl(
		ctx,
		&oapi.GetUapiDomesticStockV1QuotationsInquireCcnlParams{
			FidCondMrktDivCode: ptr("J"), // 시장 구분 코드 (J: 주식)
			FidInputIscd:       ptr(code),
			TrId:               ptr("FHKST01010300"),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respData := &uapiDomesticStockV1QuotationsInquireCcnlResponse{}
	if err := unmarshalJsonBody(resp.Body, respData); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	return validateDomesticInquireCcnlResp(respData)
}

type uapiDomesticStockV1QuotationsInquireCcnlResponse struct {
	Output []*DomesticInquireCcnl `json:"output"`
	RtCd   string                 `json:"rt_cd"`
	MsgCd  string                 `json:"msg_cd"`
	Msg1   string                 `json:"msg1"`
}

type DomesticInquireCcnl struct {
	StckCntgHour string `json:"stck_cntg_hour,omitempty" yaml:"주식체결시간,omitempty"` // 주식 체결 시간
	StckPrpr     string `json:"stck_prpr,omitempty" yaml:"주식현재가,omitempty"`       // 주식 현재가
	PrdyVrss     string `json:"prdy_vrss,omitempty" yaml:"전일대비,omitempty"`        // 전일 대비
	PrdyVrssSign string `json:"prdy_vrss_sign,omitempty" yaml:"전일대비부호,omitempty"` // 전일 대비 부호
	CntgVol      string `json:"cntg_vol,omitempty" yaml:"체결거래량,omitempty"`        // 체결 거래량
	TdayRltv     string `json:"tday_rltv,omitempty" yaml:"당일체결강도,omitempty"`      // 당일 체결강도
	PrdyCtrt     string `json:"prdy_ctrt,omitempty" yaml:"전일대비율,omitempty"`       // 전일 대비율
}

func validateDomesticInquireCcnlResp(resp *uapiDomesticStockV1QuotationsInquireCcnlResponse) ([]*DomesticInquireCcnl, error) {
	if resp.RtCd != "0" {
		return nil, fmt.Errorf("error response: rt_cd=%s, msg_cd=%s, msg1=%s", resp.RtCd, resp.MsgCd, resp.Msg1)
	}

	if resp.Output == nil {
		return nil, fmt.Errorf("no output data")
	}

	return resp.Output, nil
}
