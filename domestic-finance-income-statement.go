// 국내주식 > 종목정보 > 국내주식 손익계산서

package kinvest

import (
	"context"
	"fmt"

	"github.com/suapapa/go_kinvest/internal/oapi"
)

func (c *Client) GetDomesticFinanceIncomeStatement(ctx context.Context, code string) ([]*DomesticFinanceIncomeStatement, error) {
	if len(code) != 6 {
		return nil, fmt.Errorf("invalid item no: %s", code)
	}

	resp, err := c.oc.GetUapiDomesticStockV1FinanceIncomeStatement(
		ctx,
		&oapi.GetUapiDomesticStockV1FinanceIncomeStatementParams{
			FidCondMrktDivCode: ptr("J"), // 시장 구분 코드 (J: 주식)
			FidInputIscd:       ptr(code),
			FidDivClsCode:      ptr("1"), // 분기 결산 여부 (0: 연말, 1: 분기)
			TrId:               ptr("FHKST66430200"),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respData := &uapiDomesticStockV1FinanceIncomeStatementResponse{}
	if err := unmarshalJsonBody(resp.Body, respData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return validateDomesticFinanceIncomeStatement(respData)
}

type uapiDomesticStockV1FinanceIncomeStatementResponse struct {
	Output []*DomesticFinanceIncomeStatement `json:"output"`
	RtCd   string                            `json:"rt_cd"`
	MsgCd  string                            `json:"msg_cd"`
	Msg1   string                            `json:"msg1"`
}

type DomesticFinanceIncomeStatement struct {
	StacYymm     string `json:"stac_yymm" yaml:"결산년월,omitempty"`       // 결산 년월
	SaleAccount  string `json:"sale_account" yaml:"매출액,omitempty"`     // 매출액
	SaleCost     string `json:"sale_cost" yaml:"매출원가,omitempty"`       // 매출 원가
	SaleTotlPrfi string `json:"sale_totl_prfi" yaml:"매출총이익,omitempty"` // 매출 총 이익
	DeprCost     string `json:"depr_cost" yaml:"감가상각비,omitempty"`      // 감가상각비
	SellMang     string `json:"sell_mang" yaml:"판매및관리비,omitempty"`     // 판매 및 관리비
	BsopPrti     string `json:"bsop_prti" yaml:"영업이익,omitempty"`       // 영업 이익
	BsopNonErnn  string `json:"bsop_non_ernn" yaml:"영업외수익,omitempty"`  // 영업 외 수익
	BsopNonExpn  string `json:"bsop_non_expn" yaml:"영업외비용,omitempty"`  // 영업 외 비용
	OpPrfi       string `json:"op_prfi" yaml:"경상이익,omitempty"`         // 경상 이익
	SpecPrfi     string `json:"spec_prfi" yaml:"특별이익,omitempty"`       // 특별 이익
	SpecLoss     string `json:"spec_loss" yaml:"특별손실,omitempty"`       // 특별 손실
	ThtrNtin     string `json:"thtr_ntin" yaml:"당기순이익,omitempty"`      // 당기순이익
}

func validateDomesticFinanceIncomeStatement(data *uapiDomesticStockV1FinanceIncomeStatementResponse) ([]*DomesticFinanceIncomeStatement, error) {
	if data == nil {
		return nil, fmt.Errorf("response data is nil")
	}

	if data.RtCd != "0" {
		return nil, fmt.Errorf("response error: %s (%s)", data.Msg1, data.MsgCd)
	}

	return data.Output, nil
}
