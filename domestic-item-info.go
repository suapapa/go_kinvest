// 상품기본조회[v1_국내주식-029]

package kinvest

import (
	"context"
	"fmt"

	"github.com/suapapa/go_kinvest/internal/oapi"
)

func (c *Client) GetDomesticItemInfo(ctx context.Context, code string) (*ItemInfo, error) {
	if len(code) != 6 {
		return nil, fmt.Errorf("invalid item no: %s", code)
	}

	resp, err := c.oc.GetUapiDomesticStockV1QuotationsSearchInfo(
		ctx,
		&oapi.GetUapiDomesticStockV1QuotationsSearchInfoParams{
			PDNO:       ptr(code),
			PRDTTYPECD: ptr("300"), // 주식
			TrId:       ptr("CTPF1604R"),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respData := &uapiDomesticStockV1QuotationsSearchInfoResponse{}
	if err := unmarshalJsonBody(resp.Body, respData); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	return validateDomesticItemInfo(respData)
}

type uapiDomesticStockV1QuotationsSearchInfoResponse struct {
	Output *ItemInfo `json:"output"`
	RtCd   string    `json:"rt_cd"`
	MsgCd  string    `json:"msg_cd"`
	Msg1   string    `json:"msg1"`
}

type ItemInfo struct {
	Pdno               string `json:"pdno,omitempty" yaml:"상품번호,omitempty"`                        // 상품번호
	PrdtTypeCd         string `json:"prdt_type_cd,omitempty" yaml:"상품유형코드,omitempty"`              // 상품유형코드
	PrdtName           string `json:"prdt_name,omitempty" yaml:"상품명,omitempty"`                    // 상품명
	PrdtName120        string `json:"prdt_name120,omitempty" yaml:"상품명120,omitempty"`              // 상품명120
	PrdtAbrvName       string `json:"prdt_abrv_name,omitempty" yaml:"상품약어명,omitempty"`             // 상품약어명
	PrdtEngName        string `json:"prdt_eng_name,omitempty" yaml:"상품영문명,omitempty"`              // 상품영문명
	PrdtEngName120     string `json:"prdt_eng_name120,omitempty" yaml:"상품영문명120,omitempty"`        // 상품영문명120
	PrdtEngAbrvName    string `json:"prdt_eng_abrv_name,omitempty" yaml:"상품영문약어명,omitempty"`       // 상품영문약어명
	StdPdno            string `json:"std_pdno,omitempty" yaml:"표준상품번호,omitempty"`                  // 표준상품번호
	ShtnPdno           string `json:"shtn_pdno,omitempty" yaml:"단축상품번호,omitempty"`                 // 단축상품번호
	PrdtSaleStatCd     string `json:"prdt_sale_stat_cd,omitempty" yaml:"상품판매상태코드,omitempty"`       // 상품판매상태코드
	PrdtRiskGradCd     string `json:"prdt_risk_grad_cd,omitempty" yaml:"상품위험등급코드,omitempty"`       // 상품위험등급코드
	PrdtClsfCd         string `json:"prdt_clsf_cd,omitempty" yaml:"상품분류코드,omitempty"`              // 상품분류코드
	PrdtClsfName       string `json:"prdt_clsf_name,omitempty" yaml:"상품분류명,omitempty"`             // 상품분류명
	SaleStrtDt         string `json:"sale_strt_dt,omitempty" yaml:"판매시작일자,omitempty"`              // 판매시작일자
	SaleEndDt          string `json:"sale_end_dt,omitempty" yaml:"판매종료일자,omitempty"`               // 판매종료일자
	WrapAsstTypeCd     string `json:"wrap_asst_type_cd,omitempty" yaml:"랩어카운트자산유형코드,omitempty"`    // 랩어카운트자산유형코드
	IvstPrdtTypeCd     string `json:"ivst_prdt_type_cd,omitempty" yaml:"투자상품유형코드,omitempty"`       // 투자상품유형코드
	IvstPrdtTypeCdName string `json:"ivst_prdt_type_cd_name,omitempty" yaml:"투자상품유형코드명,omitempty"` // 투자상품유형코드명
	FrstErlmDt         string `json:"frst_erlm_dt,omitempty" yaml:"최초등록일자,omitempty"`              // 최초등록일자
}

func validateDomesticItemInfo(resp *uapiDomesticStockV1QuotationsSearchInfoResponse) (*ItemInfo, error) {
	if resp.RtCd != "0" {
		return nil, fmt.Errorf("error response: rt_cd=%s, msg_cd=%s, msg1=%s", resp.RtCd, resp.MsgCd, resp.Msg1)
	}

	if resp.Output == nil {
		return nil, fmt.Errorf("no output data")
	}

	return resp.Output, nil
}
