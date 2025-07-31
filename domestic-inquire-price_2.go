// 국내주식 > 기본시세 > 주식현재가 시세2

package kinvest

import (
	"context"
	"fmt"

	"github.com/suapapa/go_kinvest/internal/oapi"
)

func (c *Client) GetDomesticInquirePrice2(ctx context.Context, code string) (*DomesticInquirePrice2, error) {
	if len(code) != 6 {
		return nil, fmt.Errorf("invalid item no: %s", code)
	}

	resp, err := c.oc.GetUapiDomesticStockV1QuotationsInquirePrice2(
		ctx,
		&oapi.GetUapiDomesticStockV1QuotationsInquirePrice2Params{
			FidCondMrktDivCode: ptr("J"), // 시장 구분 코드 (J: 주식)
			FidInputIscd:       ptr(code),
			TrId:               ptr("FHPST01010000"),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respData := &uapiDomesticStockV1QuotationsInquirePrice2Response{}
	if err := unmarshalJsonBody(resp.Body, respData); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	return validateDomesticInquirePrice2(respData)
}

type uapiDomesticStockV1QuotationsInquirePrice2Response struct {
	Output *DomesticInquirePrice2 `json:"output"`
	RtCd   string                 `json:"rt_cd"`
	MsgCd  string                 `json:"msg_cd"`
	Msg1   string                 `json:"msg1"`
}

type DomesticInquirePrice2 struct {
	RprsMrktKorName      string `json:"rprs_mrkt_kor_name,omitempty" yaml:"대표시장한글명,omitempty"`           // 대표 시장 한글 명
	NewHgprLwprClsCode   string `json:"new_hgpr_lwpr_cls_code,omitempty" yaml:"신고가저가구분코드,omitempty"`     // 신 고가 저가 구분 코드
	MxprLlamClsCode      string `json:"mxpr_llam_cls_code,omitempty" yaml:"상하한가구분코드,omitempty"`          // 상하한가 구분 코드
	CrdtAbleYn           string `json:"crdt_able_yn,omitempty" yaml:"신용가능여부,omitempty"`                  // 신용 가능 여부
	StckMxpr             string `json:"stck_mxpr,omitempty" yaml:"주식상한가,omitempty"`                      // 주식 상한가
	ElwPblcYn            string `json:"elw_pblc_yn,omitempty" yaml:"ELW발행여부,omitempty"`                  // ELW 발행 여부
	PrdyClprVrssOprcRate string `json:"prdy_clpr_vrss_oprc_rate,omitempty" yaml:"전일종가대비시가2비율,omitempty"` // 전일 종가 대비 시가2 비율
	CrdtRate             string `json:"crdt_rate,omitempty" yaml:"신용비율,omitempty"`                       // 신용 비율
	MargRate             string `json:"marg_rate,omitempty" yaml:"증거금비율,omitempty"`                      // 증거금 비율
	LwprVrssPrpr         string `json:"lwpr_vrss_prpr,omitempty" yaml:"최저가대비현재가,omitempty"`              // 최저가 대비 현재가
	LwprVrssPrprSign     string `json:"lwpr_vrss_prpr_sign,omitempty" yaml:"최저가대비현재가부호,omitempty"`       // 최저가 대비 현재가 부호
	PrdyClprVrssLwprRate string `json:"prdy_clpr_vrss_lwpr_rate,omitempty" yaml:"전일종가대비최저가비율,omitempty"` // 전일 종가 대비 최저가 비율
	StckLwpr             string `json:"stck_lwpr,omitempty" yaml:"주식최저가,omitempty"`                      // 주식 최저가
	HgprVrssPrpr         string `json:"hgpr_vrss_prpr,omitempty" yaml:"최고가대비현재가,omitempty"`              // 최고가 대비 현재가
	HgprVrssPrprSign     string `json:"hgpr_vrss_prpr_sign,omitempty" yaml:"최고가대비현재가부호,omitempty"`       // 최고가 대비 현재가 부호
	PrdyClprVrssHgprRate string `json:"prdy_clpr_vrss_hgpr_rate,omitempty" yaml:"전일종가대비최고가비율,omitempty"` // 전일 종가 대비 최고가 비율
	StckHgpr             string `json:"stck_hgpr,omitempty" yaml:"주식최고가,omitempty"`                      // 주식 최고가
	OprcVrssPrpr         string `json:"oprc_vrss_prpr,omitempty" yaml:"시가2대비현재가,omitempty"`              // 시가2 대비 현재가
	OprcVrssPrprSign     string `json:"oprc_vrss_prpr_sign,omitempty" yaml:"시가2대비현재가부호,omitempty"`       // 시가2 대비 현재가 부호
	MangIssuYn           string `json:"mang_issu_yn,omitempty" yaml:"관리종목여부,omitempty"`                  // 관리 종목 여부
	DiviAppClsCode       string `json:"divi_app_cls_code,omitempty" yaml:"동시호가배분처리코드,omitempty"`         // 동시호가배분처리코드
	ShortOverYn          string `json:"short_over_yn,omitempty" yaml:"단기과열여부,omitempty"`                 // 단기과열여부
	MrktWarnClsCode      string `json:"mrkt_warn_cls_code,omitempty" yaml:"시장경고코드,omitempty"`            // 시장경고코드
	InvtCafulYn          string `json:"invt_caful_yn,omitempty" yaml:"투자유의여부,omitempty"`                 // 투자유의여부
	StangeRunupYn        string `json:"stange_runup_yn,omitempty" yaml:"이상급등여부,omitempty"`               // 이상급등여부
	SstsHotYn            string `json:"ssts_hot_yn,omitempty" yaml:"공매도과열여부,omitempty"`                  // 공매도과열 여부
	LowCurrentYn         string `json:"low_current_yn,omitempty" yaml:"저유동성종목여부,omitempty"`              // 저유동성 종목 여부
	ViClsCode            string `json:"vi_cls_code,omitempty" yaml:"VI적용구분코드,omitempty"`                 // VI적용구분코드
	ShortOverClsCode     string `json:"short_over_cls_code,omitempty" yaml:"단기과열구분코드,omitempty"`         // 단기과열구분코드
	StckLlam             string `json:"stck_llam,omitempty" yaml:"주식하한가,omitempty"`                      // 주식 하한가
	NewLstnClsName       string `json:"new_lstn_cls_name,omitempty" yaml:"신규상장구분명,omitempty"`            // 신규 상장 구분 명
	VlntDealClsName      string `json:"vlnt_deal_cls_name,omitempty" yaml:"임의매매구분명,omitempty"`           // 임의 매매 구분 명
	FlngClsName          string `json:"flng_cls_name,omitempty" yaml:"락구분이름,omitempty"`                  // 락 구분 이름
	RevlIssuReasName     string `json:"revl_issu_reas_name,omitempty" yaml:"재평가종목사유명,omitempty"`         // 재평가 종목 사유 명
	MrktWarnClsName      string `json:"mrkt_warn_cls_name,omitempty" yaml:"시장경고구분명,omitempty"`           // 시장 경고 구분 명
	StckSdpr             string `json:"stck_sdpr,omitempty" yaml:"주식기준가,omitempty"`                      // 주식 기준가
	BstpClsCode          string `json:"bstp_cls_code,omitempty" yaml:"업종구분코드,omitempty"`                 // 업종 구분 코드
	StckPrdyClpr         string `json:"stck_prdy_clpr,omitempty" yaml:"주식전일종가,omitempty"`                // 주식 전일 종가
	InsnPbntYn           string `json:"insn_pbnt_yn,omitempty" yaml:"불성실공시여부,omitempty"`                 // 불성실 공시 여부
	FcamModClsName       string `json:"fcam_mod_cls_name,omitempty" yaml:"액면가변경구분명,omitempty"`           // 액면가 변경 구분 명
	StckPrpr             string `json:"stck_prpr,omitempty" yaml:"주식현재가,omitempty"`                      // 주식 현재가
	PrdyVrss             string `json:"prdy_vrss,omitempty" yaml:"전일대비,omitempty"`                       // 전일 대비
	PrdyVrssSign         string `json:"prdy_vrss_sign,omitempty" yaml:"전일대비부호,omitempty"`                // 전일 대비 부호
	PrdyCtrt             string `json:"prdy_ctrt,omitempty" yaml:"전일대비율,omitempty"`                      // 전일 대비율
	AcmlTrPbmn           string `json:"acml_tr_pbmn,omitempty" yaml:"누적거래대금,omitempty"`                  // 누적 거래 대금
	AcmlVol              string `json:"acml_vol,omitempty" yaml:"누적거래량,omitempty"`                       // 누적 거래량
	PrdyVrssVolRate      string `json:"prdy_vrss_vol_rate,omitempty" yaml:"전일대비거래량비율,omitempty"`         // 전일 대비 거래량 비율
	BstpKorIsnm          string `json:"bstp_kor_isnm,omitempty" yaml:"업종한글종목명,omitempty"`                // 업종 한글 종목명
	SltrYn               string `json:"sltr_yn,omitempty" yaml:"정리매매여부,omitempty"`                       // 정리매매 여부
	TrhtYn               string `json:"trht_yn,omitempty" yaml:"거래정지여부,omitempty"`                       // 거래정지 여부
	OprcRangContYn       string `json:"oprc_rang_cont_yn,omitempty" yaml:"시가범위연장여부,omitempty"`           // 시가 범위 연장 여부
	VlntFinClsCode       string `json:"vlnt_fin_cls_code,omitempty" yaml:"임의종료구분코드,omitempty"`           // 임의 종료 구분 코드
	StckOprc             string `json:"stck_oprc,omitempty" yaml:"주식시가2,omitempty"`                      // 주식 시가2
	PrdyVol              string `json:"prdy_vol,omitempty" yaml:"전일거래량,omitempty"`                       // 전일 거래량
}

func validateDomesticInquirePrice2(data *uapiDomesticStockV1QuotationsInquirePrice2Response) (*DomesticInquirePrice2, error) {
	if data == nil {
		return nil, fmt.Errorf("response data is nil")
	}

	if data.RtCd != "0" {
		return nil, fmt.Errorf("response error: %s (%s)", data.Msg1, data.MsgCd)
	}

	return data.Output, nil
}
