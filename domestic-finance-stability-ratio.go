// 국내주식 > 종목정보 > 국내주식 안정성 비율

package kinvest

import (
	"context"
	"fmt"

	"github.com/suapapa/go_kinvest/internal/oapi"
)

func (c *Client) GetDomesticFinanceStabilityRatio(ctx context.Context, code string, anualFiscal bool) ([]*DomesticFinanceStabilityRatio, error) {
	if len(code) != 6 {
		return nil, fmt.Errorf("invalid item no: %s", code)
	}

	resp, err := c.oc.GetUapiDomesticStockV1FinanceStabilityRatio(
		ctx,
		&oapi.GetUapiDomesticStockV1FinanceStabilityRatioParams{
			FidCondMrktDivCode: ptr("J"), // 시장 구분 코드 (J: 주식)
			FidInputIscd:       ptr(code),
			FidDivClsCode:      getFiscalPeriodCode(anualFiscal),
			TrId:               ptr("FHKST66430500"),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respData := &uapiDomesticStockV1FinanceStabilityRatioResponse{}
	if err := unmarshalJsonBody(resp.Body, respData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return validateDomesticFinanceStabilityRatio(respData)
}

type uapiDomesticStockV1FinanceStabilityRatioResponse struct {
	Output []*DomesticFinanceStabilityRatio `json:"output"`
	RtCd   string                           `json:"rt_cd"`
	MsgCd  string                           `json:"msg_cd"`
	Msg1   string                           `json:"msg1"`
}

type DomesticFinanceStabilityRatio struct {
	StacYymm string `json:"stac_yymm,omitempty" yaml:"결산년월,omitempty"`   // 결산 년월
	LbltRate string `json:"lblt_rate,omitempty" yaml:"부채비율,omitempty"`   // 부채 비율
	BramDepn string `json:"bram_depn,omitempty" yaml:"차입금의존도,omitempty"` // 차입금 의존도
	CrntRate string `json:"crnt_rate,omitempty" yaml:"유동비율,omitempty"`   // 유동 비율
	QuckRate string `json:"quck_rate,omitempty" yaml:"당좌비율,omitempty"`   // 당좌 비율
}

func validateDomesticFinanceStabilityRatio(data *uapiDomesticStockV1FinanceStabilityRatioResponse) ([]*DomesticFinanceStabilityRatio, error) {
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
