// 국내주식 > 종목정보 > 국내주식 재무비율

package kinvest

import (
	"context"
	"fmt"

	"github.com/suapapa/go_kinvest/internal/oapi"
)

func (c *Client) GetDomesticFinancialRatio(ctx context.Context, code string) ([]*DomesticFinancialRatio, error) {
	if len(code) != 6 {
		return nil, fmt.Errorf("invalid item no: %s", code)
	}

	resp, err := c.oc.GetUapiDomesticStockV1FinanceFinancialRatio(
		ctx,
		&oapi.GetUapiDomesticStockV1FinanceFinancialRatioParams{
			FidCondMrktDivCode: ptr("J"), // 시장 구분 코드 (J: 주식)
			FidInputIscd:       ptr(code),
			FidDivClsCode:      ptr("1"), // 분기 결산 여부 (0: 연말, 1: 분기)
			TrId:               ptr("FHKST66430300"),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respData := &uapiDomesticStockV1FinanceFinancialRatioResponse{}
	if err := unmarshalJsonBody(resp.Body, respData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return validateDomesticFinanceRatio(respData)
}

type uapiDomesticStockV1FinanceFinancialRatioResponse struct {
	Output []*DomesticFinancialRatio `json:"output"`
	RtCd   string                    `json:"rt_cd"`
	MsgCd  string                    `json:"msg_cd"`
	Msg1   string                    `json:"msg1"`
}

type DomesticFinancialRatio struct {
	StacYymm     string `json:"stac_yymm,omitempty" yaml:"결산년월,omitempty"`         // 결산 년월
	Grs          string `json:"grs,omitempty" yaml:"매출액증가율,omitempty"`             // 매출액 증가율
	BsopPrfiInrt string `json:"bsop_prfi_inrt,omitempty" yaml:"영업이익증가율,omitempty"` // 영업 이익 증가율
	NtinInrt     string `json:"ntin_inrt,omitempty" yaml:"순이익증가율,omitempty"`       // 순이익 증가율
	RoeVal       string `json:"roe_val,omitempty" yaml:"ROE값,omitempty"`           // ROE 값
	Eps          string `json:"eps,omitempty" yaml:"EPS,omitempty"`                // EPS
	Sps          string `json:"sps,omitempty" yaml:"주당매출액,omitempty"`              // 주당매출액
	Bps          string `json:"bps,omitempty" yaml:"BPS,omitempty"`                // BPS
	RsrvRate     string `json:"rsrv_rate,omitempty" yaml:"유보비율,omitempty"`         // 유보 비율
	LbltRate     string `json:"lblt_rate,omitempty" yaml:"부채비율,omitempty"`         // 부채 비율
}

func validateDomesticFinanceRatio(data *uapiDomesticStockV1FinanceFinancialRatioResponse) ([]*DomesticFinancialRatio, error) {
	if data == nil {
		return nil, fmt.Errorf("response data is nil")
	}

	if data.RtCd != "0" {
		return nil, fmt.Errorf("response error: %s (%s)", data.Msg1, data.MsgCd)
	}

	if data.Output == nil {
		return nil, fmt.Errorf("no output data")
	}

	return data.Output, nil
}
