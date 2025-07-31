// 국내주식 > 종목정보 > 국내주식 성장성비율

package kinvest

import (
	"context"
	"fmt"

	"github.com/suapapa/go_kinvest/internal/oapi"
)

func (c *Client) GetDomesticFinanceGrowthRatio(ctx context.Context, code string, anualFiscal bool) ([]*DomesticFinanceGrowthRatio, error) {
	if len(code) != 6 {
		return nil, fmt.Errorf("invalid item no: %s", code)
	}

	resp, err := c.oc.GetUapiDomesticStockV1FinanceGrowthRatio(
		ctx,
		&oapi.GetUapiDomesticStockV1FinanceGrowthRatioParams{
			FidCondMrktDivCode: ptr("J"), // 시장 구분 코드 (J: 주식)
			FidInputIscd:       ptr(code),
			FidDivClsCode:      getFiscalPeriodCode(anualFiscal),
			TrId:               ptr("FHKST66430800"),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respData := &uapiDomesticStockV1FinanceGrowthRatioResponse{}
	if err := unmarshalJsonBody(resp.Body, respData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return validateDomesticFinanceGrowthRatio(respData)
}

type uapiDomesticStockV1FinanceGrowthRatioResponse struct {
	Output []*DomesticFinanceGrowthRatio `json:"output"`
	RtCd   string                        `json:"rt_cd"`
	MsgCd  string                        `json:"msg_cd"`
	Msg1   string                        `json:"msg1"`
}

type DomesticFinanceGrowthRatio struct {
	StacYymm     string `json:"stac_yymm,omitempty" yaml:"결산년월,omitempty"`         // 결산 년월
	Grs          string `json:"grs,omitempty" yaml:"매출액증가율,omitempty"`             // 매출액 증가율
	BsopPrfiInrt string `json:"bsop_prfi_inrt,omitempty" yaml:"영업이익증가율,omitempty"` // 영업 이익 증가율
	EqutInrt     string `json:"equt_inrt,omitempty" yaml:"자기자본증가율,omitempty"`      // 자기자본 증가율
	TotlAsetInrt string `json:"totl_aset_inrt,omitempty" yaml:"총자산증가율,omitempty"`  // 총자산 증가율
}

func validateDomesticFinanceGrowthRatio(data *uapiDomesticStockV1FinanceGrowthRatioResponse) ([]*DomesticFinanceGrowthRatio, error) {
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
