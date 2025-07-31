// 국내주식 > 종목정보 > 국내주식 수익성 비율

package kinvest

import (
	"context"
	"fmt"

	"github.com/suapapa/go_kinvest/internal/oapi"
)

func (c *Client) GetDomesticProfitRatio(ctx context.Context, code string) ([]*DomesticFinanceProfitRatio, error) {
	if len(code) != 6 {
		return nil, fmt.Errorf("invalid item no: %s", code)
	}

	resp, err := c.oc.GetUapiDomesticStockV1FinanceProfitRatio(
		ctx,
		&oapi.GetUapiDomesticStockV1FinanceProfitRatioParams{
			FidCondMrktDivCode: ptr("J"), // 시장 구분 코드 (J: 주식)
			FidInputIscd:       ptr(code),
			FidDivClsCode:      ptr("1"), // 분기 결산 여부 (0: 연말, 1: 분기)
			TrId:               ptr("FHKST66430400"),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respData := &uapiDomesticStockV1FinanceProfitRatioResponse{}
	if err := unmarshalJsonBody(resp.Body, respData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return validateDomesticFinanceProfitRatio(respData)
}

type uapiDomesticStockV1FinanceProfitRatioResponse struct {
	Output []*DomesticFinanceProfitRatio `json:"output"`
	RtCd   string                        `json:"rt_cd"`
	MsgCd  string                        `json:"msg_cd"`
	Msg1   string                        `json:"msg1"`
}

type DomesticFinanceProfitRatio struct {
	StacYymm         string `json:"stac_yymm,omitempty" yaml:"결산년월,omitempty"`               // 결산 년월
	CptlNtinRate     string `json:"cptl_ntin_rate,omitempty" yaml:"총자본순이익율,omitempty"`       // 총자본 순이익율
	SelfCptlNtinInrt string `json:"self_cptl_ntin_inrt,omitempty" yaml:"자기자본순이익율,omitempty"` // 자기자본 순이익율
	SaleNtinRate     string `json:"sale_ntin_rate,omitempty" yaml:"매출액순이익율,omitempty"`       // 매출액 순이익율
	SaleTotlRate     string `json:"sale_totl_rate,omitempty" yaml:"매출액총이익율,omitempty"`       // 매출액 총이익율
}

func validateDomesticFinanceProfitRatio(data *uapiDomesticStockV1FinanceProfitRatioResponse) ([]*DomesticFinanceProfitRatio, error) {
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
