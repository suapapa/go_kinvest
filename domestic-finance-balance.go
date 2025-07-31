// 국내주식 > 재무제표 > 국내주식 대차대조표

package kinvest

import (
	"context"
	"fmt"

	"github.com/suapapa/go_kinvest/internal/oapi"
)

func (c *Client) GetDomesticFinanceBalance(ctx context.Context, code string) ([]*DomesticFinanceBalance, error) {
	if len(code) != 6 {
		return nil, fmt.Errorf("invalid item no: %s", code)
	}

	resp, err := c.oc.GetUapiDomesticStockV1FinanceBalanceSheet(
		ctx,
		&oapi.GetUapiDomesticStockV1FinanceBalanceSheetParams{
			FidCondMrktDivCode: ptr("J"), // 시장 구분 코드 (J: 주식)
			FidInputIscd:       ptr(code),
			TrId:               ptr("FHKST66430100"),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respData := &uapiDomesticStockV1FinanceBalanceSheetResponse{}
	if err := unmarshalJsonBody(resp.Body, respData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return validateDomesticFinanceBalance(respData)
}

type uapiDomesticStockV1FinanceBalanceSheetResponse struct {
	Output []*DomesticFinanceBalance `json:"output"`
	RtCd   string                    `json:"rt_cd"`
	MsgCd  string                    `json:"msg_cd"`
	Msg1   string                    `json:"msg1"`
}

type DomesticFinanceBalance struct {
	StacYymm  string `json:"stac_yymm"`  // 결산 년월
	Cras      string `json:"cras"`       // 유동자산
	Fxas      string `json:"fxas"`       // 고정자산
	TotalAset string `json:"total_aset"` // 자산총계
	FlowLblt  string `json:"flow_lblt"`  // 유동부채
	FixLblt   string `json:"fix_lblt"`   // 고정부채
	TotalLblt string `json:"total_lblt"` // 부채총계
	Cpfn      string `json:"cpfn"`       // 자본금
	CfpSurp   string `json:"cfp_surp"`   // 자본 잉여금
	PrfiSurp  string `json:"prfi_surp"`  // 이익 잉여금
	TotalCptl string `json:"total_cptl"` // 자본총계
}

func validateDomesticFinanceBalance(data *uapiDomesticStockV1FinanceBalanceSheetResponse) ([]*DomesticFinanceBalance, error) {
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
