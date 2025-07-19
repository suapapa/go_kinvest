// 주식현재가 시세

package kinvest

import (
	"context"
	"fmt"

	"github.com/suapapa/go_kinvest/internal/oapi"
)

func (c *Client) GetDomesticInquirePrice(ctx context.Context, code string) (*DomesticInquirePrice, error) {
	if len(code) != 6 {
		return nil, fmt.Errorf("invalid item no: %s", code)
	}

	resp, err := c.oc.GetUapiDomesticStockV1QuotationsInquirePrice2(
		ctx,
		&oapi.GetUapiDomesticStockV1QuotationsInquirePrice2Params{
			FidCondMrktDivCode: ptr("J"), // 시장 구분 코드 (J: 주식)
			FidInputIscd:       ptr(code),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respData := &uapiDomesticStoecV1QuotationsInquirePrice2Response{}
	if err := unmarshalJsonBody(resp.Body, respData); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	return ValidateDomesticInquirePrice(respData)
}

type uapiDomesticStoecV1QuotationsInquirePrice2Response struct {
	Output *DomesticInquirePrice `json:"output"`
	RtCd   string                `json:"rt_cd"`
	MsgCd  string                `json:"msg_cd"`
	Msg1   string                `json:"msg1"`
}

type DomesticInquirePrice struct {
	RprsMrktKorName      string `json:"rprs_mrkt_kor_name"`       // 대표 시장 한글 명
	NewHgprLwprClsCode   string `json:"new_hgpr_lwpr_cls_code"`   // 신 고가 저가 구분 코드
	MxprLlamClsCode      string `json:"mxpr_llam_cls_code"`       // 상하한가 구분 코드
	CrdtAbleYn           string `json:"crdt_able_yn"`             // 신용 가능 여부
	StckMxpr             string `json:"stck_mxpr"`                // 주식 상한가
	ElwPblcYn            string `json:"elw_pblc_yn"`              // ELW 발행 여부
	PrdyClprVrssOprcRate string `json:"prdy_clpr_vrss_oprc_rate"` // 전일 종가 대비 시가2 비율
	CrdtRate             string `json:"crdt_rate"`                // 신용 비율
	MargRate             string `json:"marg_rate"`                // 증거금 비율
	LwprVrssPrpr         string `json:"lwpr_vrss_prpr"`           // 최저가 대비 현재가
	LwprVrssPrprSign     string `json:"lwpr_vrss_prpr_sign"`      // 최저가 대비 현재가 부호
	PrdyClprVrssLwprRate string `json:"prdy_clpr_vrss_lwpr_rate"` // 전일 종가 대비 최저가 비율
	StckLwpr             string `json:"stck_lwpr"`                // 주식 최저가
	HgprVrssPrpr         string `json:"hgpr_vrss_prpr"`           // 최고가 대비 현재가
	HgprVrssPrprSign     string `json:"hgpr_vrss_prpr_sign"`      // 최고가 대비 현재가 부호
	PrdyClprVrssHgprRate string `json:"prdy_clpr_vrss_hgpr_rate"` // 전일 종가 대비 최고가 비율
	StckHgpr             string `json:"stck_hgpr"`                // 주식 최고가
	OprcVrssPrpr         string `json:"oprc_vrss_prpr"`           // 시가2 대비 현재가
	OprcVrssPrprSign     string `json:"oprc_vrss_prpr_sign"`      // 시가2 대비 현재가 부호
	MangIssuYn           string `json:"mang_issu_yn"`             // 관리 종목 여부
	DiviAppClsCode       string `json:"divi_app_cls_code"`        // 동시호가배분처리코드
	ShortOverYn          string `json:"short_over_yn"`            // 단기과열여부
	MrktWarnClsCode      string `json:"mrkt_warn_cls_code"`       // 시장경고코드
	InvtCafulYn          string `json:"invt_caful_yn"`            // 투자유의여부
	StangeRunupYn        string `json:"stange_runup_yn"`          // 이상급등여부
	SstsHotYn            string `json:"ssts_hot_yn"`              // 공매도과열 여부
	LowCurrentYn         string `json:"low_current_yn"`           // 저유동성 종목 여부
	ViClsCode            string `json:"vi_cls_code"`              // VI적용구분코드
	ShortOverClsCode     string `json:"short_over_cls_code"`      // 단기과열구분코드
	StckLlam             string `json:"stck_llam"`                // 주식 하한가
	NewLstnClsName       string `json:"new_lstn_cls_name"`        // 신규 상장 구분 명
	VlntDealClsName      string `json:"vlnt_deal_cls_name"`       // 임의 매매 구분 명
	FlngClsName          string `json:"flng_cls_name"`            // 락 구분 이름
	RevlIssuReasName     string `json:"revl_issu_reas_name"`      // 재평가 종목 사유 명
	MrktWarnClsName      string `json:"mrkt_warn_cls_name"`       // 시장 경고 구분 명
	StckSdpr             string `json:"stck_sdpr"`                // 주식 기준가
	BstpClsCode          string `json:"bstp_cls_code"`            // 업종 구분 코드
	StckPrdyClpr         string `json:"stck_prdy_clpr"`           // 주식 전일 종가
	InsnPbntYn           string `json:"insn_pbnt_yn"`             // 불성실 공시 여부
	FcamModClsName       string `json:"fcam_mod_cls_name"`        // 액면가 변경 구분 명
	StckPrpr             string `json:"stck_prpr"`                // 주식 현재가
	PrdyVrss             string `json:"prdy_vrss"`                // 전일 대비
	PrdyVrssSign         string `json:"prdy_vrss_sign"`           // 전일 대비 부호
	PrdyCtrt             string `json:"prdy_ctrt"`                // 전일 대비율
	AcmlTrPbmn           string `json:"acml_tr_pbmn"`             // 누적 거래 대금
	AcmlVol              string `json:"acml_vol"`                 // 누적 거래량
	PrdyVrssVolRate      string `json:"prdy_vrss_vol_rate"`       // 전일 대비 거래량 비율
	BstpKorIsnm          string `json:"bstp_kor_isnm"`            // 업종 한글 종목명
	SltrYn               string `json:"sltr_yn"`                  // 정리매매 여부
	TrhtYn               string `json:"trht_yn"`                  // 거래정지 여부
	OprcRangContYn       string `json:"oprc_rang_cont_yn"`        // 시가 범위 연장 여부
	VlntFinClsCode       string `json:"vlnt_fin_cls_code"`        // 임의 종료 구분 코드
	StckOprc             string `json:"stck_oprc"`                // 주식 시가2
	PrdyVol              string `json:"prdy_vol"`                 // 전일 거래량
}

func ValidateDomesticInquirePrice(data *uapiDomesticStoecV1QuotationsInquirePrice2Response) (*DomesticInquirePrice, error) {
	if data == nil {
		return nil, fmt.Errorf("response data is nil")
	}

	if data.RtCd != "0" {
		return nil, fmt.Errorf("response error: %s (%s)", data.Msg1, data.MsgCd)
	}

	return data.Output, nil
}
