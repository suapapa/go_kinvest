// 국내주식 > 기본시세 > 주식현재가 시세

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

	resp, err := c.oc.GetUapiDomesticStockV1QuotationsInquirePrice(
		ctx,
		&oapi.GetUapiDomesticStockV1QuotationsInquirePriceParams{
			FidCondMrktDivCode: ptr("J"), // 시장 구분 코드 (J: 주식)
			FidInputIscd:       ptr(code),
			TrId:               ptr("FHKST01010100"),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respData := &uapiDomesticStoecV1QuotationsInquirePriceResponse{}
	if err := unmarshalJsonBody(resp.Body, respData); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	return validateDomesticInquirePriceResp(respData)
}

type uapiDomesticStoecV1QuotationsInquirePriceResponse struct {
	Output *DomesticInquirePrice `json:"output"`
	RtCd   string                `json:"rt_cd"`
	MsgCd  string                `json:"msg_cd"`
	Msg1   string                `json:"msg1"`
}

type DomesticInquirePrice struct {
	IscdStatClsCode      string `json:"iscd_stat_cls_code,omitempty" yaml:"종목상태구분코드,omitempty"`             // 종목 상태 구분 코드
	MargRate             string `json:"marg_rate,omitempty" yaml:"증거금비율,omitempty"`                         // 증거금 비율
	RprsMrktKorName      string `json:"rprs_mrkt_kor_name,omitempty" yaml:"대표시장한글명,omitempty"`              // 대표 시장 한글 명
	NewHgprLwprClsCode   string `json:"new_hgpr_lwpr_cls_code,omitempty" yaml:"신고가저가구분코드,omitempty"`        // 신 고가 저가 구분 코드
	BstpKorIsnm          string `json:"bstp_kor_isnm,omitempty" yaml:"업종한글종목명,omitempty"`                   // 업종 한글 종목명
	TempStopYn           string `json:"temp_stop_yn,omitempty" yaml:"임시정지여부,omitempty"`                     // 임시 정지 여부
	OprcRangContYn       string `json:"oprc_rang_cont_yn,omitempty" yaml:"시가범위연장여부,omitempty"`              // 시가 범위 연장 여부
	ClprRangContYn       string `json:"clpr_rang_cont_yn,omitempty" yaml:"종가범위연장여부,omitempty"`              // 종가 범위 연장 여부
	CrdtAbleYn           string `json:"crdt_able_yn,omitempty" yaml:"신용가능여부,omitempty"`                     // 신용 가능 여부
	GrmnRateClsCode      string `json:"grmn_rate_cls_code,omitempty" yaml:"보증금비율구분코드,omitempty"`            // 보증금 비율 구분 코드
	ElwPblcYn            string `json:"elw_pblc_yn,omitempty" yaml:"ELW발행여부,omitempty"`                     // ELW 발행 여부
	StckPrpr             string `json:"stck_prpr,omitempty" yaml:"주식현재가,omitempty"`                         // 주식 현재가
	PrdyVrss             string `json:"prdy_vrss,omitempty" yaml:"전일대비,omitempty"`                          // 전일 대비
	PrdyVrssSign         string `json:"prdy_vrss_sign,omitempty" yaml:"전일대비부호,omitempty"`                   // 전일 대비 부호
	PrdyCtrt             string `json:"prdy_ctrt,omitempty" yaml:"전일대비율,omitempty"`                         // 전일 대비율
	AcmlTrPbmn           string `json:"acml_tr_pbmn,omitempty" yaml:"누적거래대금,omitempty"`                     // 누적 거래 대금
	AcmlVol              string `json:"acml_vol,omitempty" yaml:"누적거래량,omitempty"`                          // 누적 거래량
	PrdyVrssVolRate      string `json:"prdy_vrss_vol_rate,omitempty" yaml:"전일대비거래량비율,omitempty"`            // 전일 대비 거래량 비율
	StckOprc             string `json:"stck_oprc,omitempty" yaml:"주식시가2,omitempty"`                         // 주식 시가2
	StckHgpr             string `json:"stck_hgpr,omitempty" yaml:"주식최고가,omitempty"`                         // 주식 최고가
	StckLwpr             string `json:"stck_lwpr,omitempty" yaml:"주식최저가,omitempty"`                         // 주식 최저가
	StckMxpr             string `json:"stck_mxpr,omitempty" yaml:"주식상한가,omitempty"`                         // 주식 상한가
	StckLlam             string `json:"stck_llam,omitempty" yaml:"주식하한가,omitempty"`                         // 주식 하한가
	StckSdpr             string `json:"stck_sdpr,omitempty" yaml:"주식기준가,omitempty"`                         // 주식 기준가
	WghnAvrgStckPrc      string `json:"wghn_avrg_stck_prc,omitempty" yaml:"가중평균주식가격,omitempty"`             // 가중 평균 주식 가격
	HtsFrgnEhrt          string `json:"hts_frgn_ehrt,omitempty" yaml:"HTS외국인소진율,omitempty"`                 // HTS 외국인 소진율
	FrgnNtbyQty          string `json:"frgn_ntby_qty,omitempty" yaml:"외국인순매수수량,omitempty"`                  // 외국인 순매수 수량
	PgtrNtbyQty          string `json:"pgtr_ntby_qty,omitempty" yaml:"프로그램매매순매수수량,omitempty"`               // 프로그램매매 순매수 수량
	PvtScndDmrsPrc       string `json:"pvt_scnd_dmrs_prc,omitempty" yaml:"피벗2차디저항가격,omitempty"`             // 피벗 2차 디저항 가격
	PvtFrstDmrsPrc       string `json:"pvt_frst_dmrs_prc,omitempty" yaml:"피벗1차디저항가격,omitempty"`             // 피벗 1차 디저항 가격
	PvtPontVal           string `json:"pvt_pont_val,omitempty" yaml:"피벗포인트값,omitempty"`                     // 피벗 포인트 값
	PvtFrstDmspPrc       string `json:"pvt_frst_dmsp_prc,omitempty" yaml:"피벗1차디지지가격,omitempty"`             // 피벗 1차 디지지 가격
	PvtScndDmspPrc       string `json:"pvt_scnd_dmsp_prc,omitempty" yaml:"피벗2차디지지가격,omitempty"`             // 피벗 2차 디지지 가격
	DmrsVal              string `json:"dmrs_val,omitempty" yaml:"디저항값,omitempty"`                           // 디저항 값
	DmspVal              string `json:"dmsp_val,omitempty" yaml:"디지지값,omitempty"`                           // 디지지 값
	Cpfn                 string `json:"cpfn,omitempty" yaml:"자본금,omitempty"`                                // 자본금
	RstcWdthPrc          string `json:"rstc_wdth_prc,omitempty" yaml:"제한폭가격,omitempty"`                     // 제한 폭 가격
	StckFcam             string `json:"stck_fcam,omitempty" yaml:"주식액면가,omitempty"`                         // 주식 액면가
	StckSspr             string `json:"stck_sspr,omitempty" yaml:"주식대용가,omitempty"`                         // 주식 대용가
	AsprUnit             string `json:"aspr_unit,omitempty" yaml:"호가단위,omitempty"`                          // 호가단위
	HtsDealQtyUnitVal    string `json:"hts_deal_qty_unit_val,omitempty" yaml:"HTS매매수량단위값,omitempty"`        // HTS 매매 수량 단위 값
	LstnStcn             string `json:"lstn_stcn,omitempty" yaml:"상장주수,omitempty"`                          // 상장 주수
	HtsAvls              string `json:"hts_avls,omitempty" yaml:"HTS시가총액,omitempty"`                        // HTS 시가총액
	Per                  string `json:"per,omitempty" yaml:"PER,omitempty"`                                 // PER
	Pbr                  string `json:"pbr,omitempty" yaml:"PBR,omitempty"`                                 // PBR
	StacMonth            string `json:"stac_month,omitempty" yaml:"결산월,omitempty"`                          // 결산 월
	VolTnrt              string `json:"vol_tnrt,omitempty" yaml:"거래량회전율,omitempty"`                         // 거래량 회전율
	Eps                  string `json:"eps,omitempty" yaml:"EPS,omitempty"`                                 // EPS
	Bps                  string `json:"bps,omitempty" yaml:"BPS,omitempty"`                                 // BPS
	D250Hgpr             string `json:"d250_hgpr,omitempty" yaml:"250일최고가,omitempty"`                       // 250일 최고가
	D250HgprDate         string `json:"d250_hgpr_date,omitempty" yaml:"250일최고가일자,omitempty"`                // 250일 최고가 일자
	D250HgprVrssPrprRate string `json:"d250_hgpr_vrss_prpr_rate,omitempty" yaml:"250일최고가대비현재가비율,omitempty"` // 250일 최고가 대비 현재가 비율
	D250Lwpr             string `json:"d250_lwpr,omitempty" yaml:"250일최저가,omitempty"`                       // 250일 최저가
	D250LwprDate         string `json:"d250_lwpr_date,omitempty" yaml:"250일최저가일자,omitempty"`                // 250일 최저가 일자
	D250LwprVrssPrprRate string `json:"d250_lwpr_vrss_prpr_rate,omitempty" yaml:"250일최저가대비현재가비율,omitempty"` // 250일 최저가 대비 현재가 비율
	StckDryyHgpr         string `json:"stck_dryy_hgpr,omitempty" yaml:"주식연중최고가,omitempty"`                  // 주식 연중 최고가
	DryyHgprVrssPrprRate string `json:"dryy_hgpr_vrss_prpr_rate,omitempty" yaml:"연중최고가대비현재가비율,omitempty"`   // 연중 최고가 대비 현재가 비율
	DryyHgprDate         string `json:"dryy_hgpr_date,omitempty" yaml:"연중최고가일자,omitempty"`                  // 연중 최고가 일자
	StckDryyLwpr         string `json:"stck_dryy_lwpr,omitempty" yaml:"주식연중최저가,omitempty"`                  // 주식 연중 최저가
	DryyLwprVrssPrprRate string `json:"dryy_lwpr_vrss_prpr_rate,omitempty" yaml:"연중최저가대비현재가비율,omitempty"`   // 연중 최저가 대비 현재가 비율
	DryyLwprDate         string `json:"dryy_lwpr_date,omitempty" yaml:"연중최저가일자,omitempty"`                  // 연중 최저가 일자
	W52Hgpr              string `json:"w52_hgpr,omitempty" yaml:"52주일최고가,omitempty"`                        // 52주일 최고가
	W52HgprVrssPrprCtrt  string `json:"w52_hgpr_vrss_prpr_ctrt,omitempty" yaml:"52주일최고가대비현재가대비,omitempty"`  // 52주일 최고가 대비 현재가 대비
	W52HgprDate          string `json:"w52_hgpr_date,omitempty" yaml:"52주일최고가일자,omitempty"`                 // 52주일 최고가 일자
	W52Lwpr              string `json:"w52_lwpr,omitempty" yaml:"52주일최저가,omitempty"`                        // 52주일 최저가
	W52LwprVrssPrprCtrt  string `json:"w52_lwpr_vrss_prpr_ctrt,omitempty" yaml:"52주일최저가대비현재가대비,omitempty"`  // 52주일 최저가 대비 현재가 대비
	W52LwprDate          string `json:"w52_lwpr_date,omitempty" yaml:"52주일최저가일자,omitempty"`                 // 52주일 최저가 일자
	WholLoanRmndRate     string `json:"whol_loan_rmnd_rate,omitempty" yaml:"전체융자잔고비율,omitempty"`            // 전체 융자 잔고 비율
	SstsYn               string `json:"ssts_yn,omitempty" yaml:"공매도가능여부,omitempty"`                         // 공매도가능여부
	StckShrnIscd         string `json:"stck_shrn_iscd,omitempty" yaml:"주식단축종목코드,omitempty"`                 // 주식 단축 종목코드
	FcamCnnm             string `json:"fcam_cnnm,omitempty" yaml:"액면가통화명,omitempty"`                        // 액면가 통화명
	CpfnCnnm             string `json:"cpfn_cnnm,omitempty" yaml:"자본금통화명,omitempty"`                        // 자본금 통화명
	ApprchRate           string `json:"apprch_rate,omitempty" yaml:"접근도,omitempty"`                         // 접근도
	FrgnHldnQty          string `json:"frgn_hldn_qty,omitempty" yaml:"외국인보유수량,omitempty"`                   // 외국인 보유 수량
	ViClsCode            string `json:"vi_cls_code,omitempty" yaml:"VI적용구분코드,omitempty"`                    // VI적용구분코드
	OvtmViClsCode        string `json:"ovtm_vi_cls_code,omitempty" yaml:"시간외단일가VI적용구분코드,omitempty"`         // 시간외단일가VI적용구분코드
	LastSstsCntgQty      string `json:"last_ssts_cntg_qty,omitempty" yaml:"최종공매도체결수량,omitempty"`            // 최종 공매도 체결 수량
	InvtCafulYn          string `json:"invt_caful_yn,omitempty" yaml:"투자유의여부,omitempty"`                    // 투자유의여부
	MrktWarnClsCode      string `json:"mrkt_warn_cls_code,omitempty" yaml:"시장경고코드,omitempty"`               // 시장경고코드
	ShortOverYn          string `json:"short_over_yn,omitempty" yaml:"단기과열여부,omitempty"`                    // 단기과열여부
	SltrYn               string `json:"sltr_yn,omitempty" yaml:"정리매매여부,omitempty"`                          // 정리매매여부
	MangIssuClsCode      string `json:"mang_issu_cls_code,omitempty" yaml:"관리종목여부,omitempty"`               // 관리종목여부
}

func validateDomesticInquirePriceResp(resp *uapiDomesticStoecV1QuotationsInquirePriceResponse) (*DomesticInquirePrice, error) {
	if resp.RtCd != "0" {
		return nil, fmt.Errorf("error response: rt_cd=%s, msg_cd=%s, msg1=%s", resp.RtCd, resp.MsgCd, resp.Msg1)
	}

	if resp.Output == nil {
		return nil, fmt.Errorf("no output data")
	}

	return resp.Output, nil
}
