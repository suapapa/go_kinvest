// 주식잔고조회

package kinvest

import (
	"context"
	"fmt"

	"github.com/suapapa/go_kinvest/internal/oapi"
)

// GetDomesticAccountBalance retrieves the balance of the domestic account
func (c *Client) GetDomesticAccountBalance(ctx context.Context) (*DomesticAccountBalance, error) {
	cano, acntprdtcd, err := parseAccount(c.account)
	if err != nil {
		return nil, fmt.Errorf("parse account failed: %w", err)
	}

	resp, err := c.oc.GetUapiDomesticStockV1TradingInquireAccountBalance(
		ctx,
		&oapi.GetUapiDomesticStockV1TradingInquireAccountBalanceParams{
			CANO:           cano,
			ACNTPRDTCD:     acntprdtcd,
			INQRDVSN1:      ptr(""),
			BSPRBFDTAPLYYN: ptr(""),
			TrId:           ptr("CTRP6548R"),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respData := &uapiDomesticStockV1TradingInquireAccountBalanceResp{}
	if err := unmarshalJsonBody(resp.Body, respData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return NewDomesticAccountBalance(respData)
}

type DomesticAccountBalanceItem struct {
	PchsAmt     int     `yaml:"매입금매,omitempty"`
	EvluAmt     int     `yaml:"평가금액,omitempty"`
	EvluPflsAmt int     `yaml:"평가손익금액,omitempty"`
	CrdtLndAmt  int     `yaml:"신용대출금액,omitempty"`
	RealNassAmt int     `yaml:"실현손익금액,omitempty"`
	WholWeitRt  float64 `yaml:"전체비중율,omitempty"`
}

// NewDomesticAccountBalanceItem creates a new DomesticAccountBalanceItem from the response data
func NewDomesticAccountBalanceItem(data *output1) (*DomesticAccountBalanceItem, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	item := &DomesticAccountBalanceItem{
		PchsAmt:     toInt(data.PchsAmt),
		EvluAmt:     toInt(data.EvluAmt),
		EvluPflsAmt: toInt(data.EvluPflsAmt),
		CrdtLndAmt:  toInt(data.CrdtLndAmt),
		RealNassAmt: toInt(data.RealNassAmt),
		WholWeitRt:  toFloat(data.WholWeitRt),
	}

	return item, nil
}

// empty checks if the item is empty
func (i *DomesticAccountBalanceItem) empty() bool {
	return i == nil || (i.PchsAmt == 0 && i.EvluAmt == 0 && i.EvluPflsAmt == 0 && i.CrdtLndAmt == 0 && i.RealNassAmt == 0)
}

// DomesticAccountBalance represents the balance of a domestic account
type DomesticAccountBalance struct {
	Items                  map[string]*DomesticAccountBalanceItem `yaml:"계좌항목,omitempty"`
	PchsAmtSmtl            int                                    `yaml:"매입금액합계,omitempty"`
	NassTotAmt             int                                    `yaml:"순자산총금액,omitempty"`
	LoanAmtSmtl            int                                    `yaml:"대출금액합계,omitempty"`
	EvluPflsAmtSmtl        int                                    `yaml:"평가손익금액합계,omitempty"`
	EvluAmtSmtl            int                                    `yaml:"평가금액합계,omitempty"`
	TotAsstAmt             int                                    `yaml:"총자산금액,omitempty"`
	TotLndaTotUlstLnda     int                                    `yaml:"총대출금액총융자대출금액,omitempty"`
	CmaAutoLoanAmt         int                                    `yaml:"CMA자동대출금액,omitempty"`
	TotMglnAmt             int                                    `yaml:"총담보대출금액,omitempty"`
	StlnEvluAmt            int                                    `yaml:"대주평가금액,omitempty"`
	CrdtFncgAmt            int                                    `yaml:"신용융자금융,omitempty"`
	OclAplLoanAmt          int                                    `yaml:"OCL_APL대출금액,omitempty"`
	PldgStupAmt            int                                    `yaml:"질권설정금액,omitempty"`
	FrcrEvluTota           int                                    `yaml:"외화평가총액,omitempty"`
	TotDnclAmt             int                                    `yaml:"총예수금액,omitempty"`
	CmaEvluAmt             int                                    `yaml:"CMA평가금액,omitempty"`
	DnclAmt                int                                    `yaml:"예수금액,omitempty"`
	TotSbstAmt             int                                    `yaml:"총대용금액,omitempty"`
	ThdtRcvbAmt            int                                    `yaml:"당일미수금액,omitempty"`
	OvrsStckEvluAmt1       int                                    `yaml:"해외주식평가금액1,omitempty"`
	OvrsBondEvluAmt        int                                    `yaml:"해외채권평가금액,omitempty"`
	MmfCmaMggeLoanAmt      int                                    `yaml:"MMFCMA담보대출금액,omitempty"`
	SbscDnclAmt            int                                    `yaml:"청약예수금액,omitempty"`
	PbstSbscFndsLoanUseAmt int                                    `yaml:"공모주청약자금대출사용금액,omitempty"`
	EtprCrdtGrntLoanAmt    int                                    `yaml:"기업신용공예대출금액,omitempty"`
}

// NewDomesticAccountBalance creates a new DomesticAccountBalance from the response data
func NewDomesticAccountBalance(data *uapiDomesticStockV1TradingInquireAccountBalanceResp) (*DomesticAccountBalance, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	bItems := make(map[string]*DomesticAccountBalanceItem)
	if len(data.Output1) == 19 {
		for i, key := range output1Items19 {
			item, err := NewDomesticAccountBalanceItem(&data.Output1[i])
			if err != nil {
				return nil, fmt.Errorf("failed to create item %d: %w", i, err)
			}
			if !item.empty() {
				bItems[key] = item
			}
		}
	} else if len(data.Output1) == 16 {
		for i, key := range output1Items16 {
			item, err := NewDomesticAccountBalanceItem(&data.Output1[i])
			if err != nil {
				return nil, fmt.Errorf("failed to create item %d: %w", i, err)
			}
			if !item.empty() {
				bItems[key] = item
			}
		}
	} else if len(data.Output1) == 0 {
		bItems = nil
	} else {
		return nil, fmt.Errorf("unexpected number of items: %d", len(data.Output1))
	}

	return &DomesticAccountBalance{
		Items:                  bItems,
		PchsAmtSmtl:            toInt(data.Output2.PchsAmtSmtl),
		NassTotAmt:             toInt(data.Output2.NassTotAmt),
		LoanAmtSmtl:            toInt(data.Output2.LoanAmtSmtl),
		EvluPflsAmtSmtl:        toInt(data.Output2.EvluPflsAmtSmtl),
		EvluAmtSmtl:            toInt(data.Output2.EvluAmtSmtl),
		TotAsstAmt:             toInt(data.Output2.TotAsstAmt),
		TotLndaTotUlstLnda:     toInt(data.Output2.TotLndaTotUlstLnda),
		CmaAutoLoanAmt:         toInt(data.Output2.CmaAutoLoanAmt),
		TotMglnAmt:             toInt(data.Output2.TotMglnAmt),
		StlnEvluAmt:            toInt(data.Output2.StlnEvluAmt),
		CrdtFncgAmt:            toInt(data.Output2.CrdtFncgAmt),
		OclAplLoanAmt:          toInt(data.Output2.OclAplLoanAmt),
		PldgStupAmt:            toInt(data.Output2.PldgStupAmt),
		FrcrEvluTota:           toInt(data.Output2.FrcrEvluTota),
		TotDnclAmt:             toInt(data.Output2.TotDnclAmt),
		CmaEvluAmt:             toInt(data.Output2.CmaEvluAmt),
		DnclAmt:                toInt(data.Output2.DnclAmt),
		TotSbstAmt:             toInt(data.Output2.TotSbstAmt),
		ThdtRcvbAmt:            toInt(data.Output2.ThdtRcvbAmt),
		OvrsStckEvluAmt1:       toInt(data.Output2.OvrsStckEvluAmt1),
		OvrsBondEvluAmt:        toInt(data.Output2.OvrsBondEvluAmt),
		MmfCmaMggeLoanAmt:      toInt(data.Output2.MmfCmaMggeLoanAmt),
		SbscDnclAmt:            toInt(data.Output2.SbscDnclAmt),
		PbstSbscFndsLoanUseAmt: toInt(data.Output2.PbstSbscFndsLoanUseAmt),
		EtprCrdtGrntLoanAmt:    toInt(data.Output2.EtprCrdtGrntLoanAmt),
	}, nil
}

var (
	output1Items19 = []string{
		"주식",
		"펀드/MMW",
		"채권",
		"ELS/DLS",
		"WRAP",
		"신탁/퇴직연금/외화신탁",
		"RP/발행어음",
		"해외주식",
		"해외채권",
		"금현물",
		"CD/CP",
		"단기사채",
		"타사상품",
		"외화단기사채",
		"외화 ELS/DLS",
		"외화",
		"예수금+CMA",
		"청약자예수금",
		"합계",
	}

	output1Items16 = []string{
		"수익증권",
		"채권",
		"ELS/DLS",
		"wrap",
		"신탁",
		"rp",
		"외화rp",
		"해외채권",
		"CD/CP",
		"전자단기사채",
		"외화전자단기사채",
		"외화ELS/DLS",
		"외화평가금액",
		"예수금+cma",
		"청약자예수금",
		"합계",
	}
)

type uapiDomesticStockV1TradingInquireAccountBalanceResp struct {
	Output1 []output1 `json:"output1"` // 응답상세1
	Output2 output2   `json:"output2"` // 응답상세2
	RtCd    string    `json:"rt_cd"`
	MsgCd   string    `json:"msg_cd"`
	Msg1    string    `json:"msg1"`
}

type output1 struct {
	PchsAmt     string `json:"pchs_amt"`      // 매입금매
	EvluAmt     string `json:"evlu_amt"`      // 평가금액
	EvluPflsAmt string `json:"evlu_pfls_amt"` // 평가손익금액
	CrdtLndAmt  string `json:"crdt_lnd_amt"`  // 신용대출금액
	RealNassAmt string `json:"real_nass_amt"` // 실현손익금액
	WholWeitRt  string `json:"whol_weit_rt"`  // 전체비중율
}

type output2 struct {
	PchsAmtSmtl            string `json:"pchs_amt_smtl"`               // 매입금액합계
	NassTotAmt             string `json:"nass_tot_amt"`                // 순자산총금액
	LoanAmtSmtl            string `json:"loan_amt_smtl"`               // 대출금액합계
	EvluPflsAmtSmtl        string `json:"evlu_pfls_amt_smtl"`          // 평가손익금액합계
	EvluAmtSmtl            string `json:"evlu_amt_smtl"`               // 평가금액합계
	TotAsstAmt             string `json:"tot_asst_amt"`                // 총자산금액
	TotLndaTotUlstLnda     string `json:"tot_lnda_tot_ulst_lnda"`      // 총대출금액총융자대출금액
	CmaAutoLoanAmt         string `json:"cma_auto_loan_amt"`           // CMA자동대출금액
	TotMglnAmt             string `json:"tot_mgln_amt"`                // 총담보대출금액
	StlnEvluAmt            string `json:"stln_evlu_amt"`               // 대주평가금액
	CrdtFncgAmt            string `json:"crdt_fncg_amt"`               // 신용융자금융
	OclAplLoanAmt          string `json:"ocl_apl_loan_amt"`            // OCL_APL대출금액
	PldgStupAmt            string `json:"pldg_stup_amt"`               // 질권설정금액
	FrcrEvluTota           string `json:"frcr_evlu_tota"`              // 외화평가총액
	TotDnclAmt             string `json:"tot_dncl_amt"`                // 총예수금액
	CmaEvluAmt             string `json:"cma_evlu_amt"`                // CMA평가금액
	DnclAmt                string `json:"dncl_amt"`                    // 예수금액
	TotSbstAmt             string `json:"tot_sbst_amt"`                // 총대용금액
	ThdtRcvbAmt            string `json:"thdt_rcvb_amt"`               // 당일미수금액
	OvrsStckEvluAmt1       string `json:"ovrs_stck_evlu_amt1"`         // 해외주식평가금액1
	OvrsBondEvluAmt        string `json:"ovrs_bond_evlu_amt"`          // 해외채권평가금액
	MmfCmaMggeLoanAmt      string `json:"mmf_cma_mgge_loan_amt"`       // MMFCMA담보대출금액
	SbscDnclAmt            string `json:"sbsc_dncl_amt"`               // 청약예수금액
	PbstSbscFndsLoanUseAmt string `json:"pbst_sbsc_fnds_loan_use_amt"` // 공모주청약자금대출사용금액
	EtprCrdtGrntLoanAmt    string `json:"etpr_crdt_grnt_loan_amt"`     // 기업신용공예대출금액
}
