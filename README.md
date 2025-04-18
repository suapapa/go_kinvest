# kinvest : A Go package for 한국투자증권 client

![logo](_img/go_kinvest.png)

## Usage

Install pkg:
```sh
go get github.com/suapapa/go_kinvest
```

Example usage:
```go
import kinvest "github.com/suapapa/go_kinvest"
// ...
kc, _ := kinvest.NewClient(nil)
bal, _ := kc.GetDomesticAccountBalance()
```

This package read following envs for the API calls:
- `KINVEST_ACCOUNT` : 계좌번호, XXXXXXXX-XX
- `KINVEST_APPKEY` : 한국투자증권 개발자센터에서 발급받은 appkey
- `KINVEST_APPSECRET` : 한국투자증권 개발자센터에서 발급받은 appsecret
- `KINVEST_TOKEN_PATH` : 발급받은 토큰을 저장하기 위한 경로. 설정하지 않으면 `./kinvest_access_token.json` 에 저장

## TODO
- [ ] /oauth2/Approval (post) : 웹소켓접속키발급
- [x] /oauth2/tokenP (post) : 토큰발급(선물옵션)
- [ ] /oauth2/revokeP (post) : 토큰폐기(선물옵션)
- [ ] /uapi/hashkey (post) : 해쉬키생성(선물옵션)
- [x] /uapi/domestic-stock/v1/trading/order-cash (post) : 주식주문(현금)
- [ ] /uapi/domestic-stock/v1/trading/order-credit (post) : 주식주문(신용)
- [ ] /uapi/domestic-stock/v1/trading/order-rvsecncl (post) : 주식주문(정정취소)
- [ ] /uapi/domestic-stock/v1/trading/inquire-psbl-rvsecncl (get) : 주식정정취소가능주문조회
- [ ] /uapi/domestic-stock/v1/trading/inquire-daily-ccld (get) : 주식일별주문체결조회
- [ ] /uapi/domestic-stock/v1/trading/inquire-balance (get) : 주식잔고조회
- [ ] /uapi/domestic-stock/v1/trading/inquire-psbl-order (get) : 매수가능조회
- [ ] /uapi/domestic-stock/v1/trading/order-resv (post) : 주식예약주문
- [ ] /uapi/domestic-stock/v1/trading/order-resv-rvsecncl (post) : 주식예약주문정정취소(정정)
- [ ] /uapi/domestic-stock/v1/trading/order-resv-ccnl (get) : 주식예약주문조회
- [ ] /uapi/domestic-stock/v1/trading/pension/inquire-present-balance (get) : 퇴직연금체결기준잔고
- [ ] /uapi/domestic-stock/v1/trading/pension/inquire-daily-ccld (get) : 퇴직연금 미체결내역
- [ ] /uapi/domestic-stock/v1/trading/pension/inquire-psbl-order (get) : 퇴직연금 매수가능조회
- [ ] /uapi/domestic-stock/v1/trading/pension/inquire-deposit (get) : 퇴직연금 예수금조회
- [ ] /uapi/domestic-stock/v1/trading/pension/inquire-balance (get) : 퇴직연금 잔고조회
- [ ] /uapi/domestic-stock/v1/trading/inquire-balance-rlz-pl (get) : 주식잔고조회_실현손익
- [ ] /uapi/domestic-stock/v1/trading/inquire-credit-psamount (get) : 신용매수가능조회
- [x] /uapi/domestic-stock/v1/trading/inquire-account-balance (get) : 투자계좌자산현황조회
- [ ] /uapi/domestic-stock/v1/trading/inquire-period-trade-profit (get) : 기간별매매손익현황조회
- [ ] /uapi/domestic-stock/v1/trading/inquire-period-profit (get) : 기간별손익일별합산조회
- [ ] /uapi/domestic-stock/v1/trading/inquire-psbl-sell (get) : 매도가능수량조회
- [ ] /uapi/domestic-stock/v1/quotations/inquire-price (get) : 주식현재가 시세
- [ ] /uapi/domestic-stock/v1/quotations/inquire-ccnl (get) : 주식현재가 체결(최근30건)
- [ ] /uapi/domestic-stock/v1/quotations/inquire-daily-price (get) : ELW 당일급변종목
- [ ] /uapi/domestic-stock/v1/quotations/inquire-asking-price-exp-ccn (get) : 주식현재가 호가 예상체결
- [ ] /uapi/domestic-stock/v1/quotations/inquire-investor (get) : 주식현재가 투자자
- [ ] /uapi/domestic-stock/v1/quotations/inquire-member (get) : 주식현재가 회원사
- [ ] /uapi/domestic-stock/v1/quotations/inquire-daily-itemchartprice (get) : 국내주식기간별시세(일/주/월/년)
- [ ] /uapi/domestic-stock/v1/quotations/inquire-time-itemconclusion (get) : 주식현재가 당일시간대별체결
- [ ] /uapi/domestic-stock/v1/quotations/inquire-time-overtimeconclusion (get) : 주식현재가 시간외 시간별체결
- [ ] /uapi/domestic-stock/v1/quotations/inquire-daily-overtimeprice (get) : 주식현재가 시간외 일자별주가
- [ ] /uapi/domestic-stock/v1/quotations/inquire-time-itemchartprice (get) : 주식당일분봉조회(주식)
- [ ] /uapi/domestic-stock/v1/quotations/inquire-daily-indexchartprice (get) : 국내주식업종기간별시세(일/주/월/년)
- [ ] /uapi/domestic-stock/v1/quotations/inquire-price-2 (get) : 주식현재가 시세2
- [ ] /uapi/etfetn/v1/quotations/inquire-price (get) : ETF/ETN현재가
- [ ] /uapi/etfetn/v1/quotations/nav-comparison-trend (get) : NAV 비교추이(종목)
- [ ] /uapi/etfetn/v1/quotations/nav-comparison-time-trend (get) : NAV 비교추이(분)
- [ ] /uapi/etfetn/v1/quotations/nav-comparison-daily-trend (get) : NAV 비교추이(일)
- [ ] /uapi/domestic-stock/v1/quotations/exp-closing-price (get) : 국내주식 장마감 예상체결가
- [ ] /uapi/etfetn/v1/quotations/inquire-component-stock-price (get) : ETF 구성종목시세
- [ ] /uapi/domestic-stock/v1/quotations/inquire-overtime-price (get) : 국내주식 시간외현재가
- [ ] /uapi/domestic-stock/v1/quotations/inquire-overtime-asking-price (get) : 국내주식 시간외호가
- [ ] /uapi/domestic-stock/v1/quotations/inquire-elw-price (get) : ELW현재가 시세
- [ ] /uapi/elw/v1/ranking/updown-rate (get) : ELW 상승률순위
- [ ] /uapi/elw/v1/quotations/newly-listed (get) : ELW 신규상장종목
- [ ] /uapi/elw/v1/quotations/indicator-trend-ccnl (get) : ELW 투자지표추이(체결)
- [ ] /uapi/elw/v1/quotations/indicator-trend-minute (get) : ELW 변동성추이(일별)
- [ ] /uapi/elw/v1/quotations/indicator-trend-daily (get) : ELW 투자지표추이(일별)
- [ ] /uapi/elw/v1/quotations/volatility-trend-tick (get) : ELW 변동성추이(틱)
- [ ] /uapi/elw/v1/quotations/volatility-trend-ccnl (get) : ELW 변동성 추이(체결)
- [ ] /uapi/elw/v1/quotations/volatility-trend-minute (get) : ELW 변동성 추이(분별)
- [ ] /uapi/elw/v1/quotations/sensitivity-trend-ccnl (get) : ELW 민감도 추이(체결)
- [ ] /uapi/elw/v1/quotations/sensitivity-trend-daily (get) : ELW 민감도 추이(일별)
- [ ] /uapi/elw/v1/quotations/volatility-trend-daily (get) : ELW 기초자산별 종목시세
- [ ] /uapi/elw/v1/quotations/lp-trade-trend (get) : ELW LP매매추이
- [ ] /uapi/elw/v1/quotations/compare-stocks (get) : ELW 비교대상종목조회
- [ ] /uapi/elw/v1/quotations/cond-search (get) : ELW 종목검색
- [ ] /uapi/elw/v1/quotations/udrl-asset-list (get) : ELW 기초자산 목록조회
- [ ] /uapi/elw/v1/quotations/expiration-stocks (get) : ELW 만기예정/만기종목
- [ ] /uapi/domestic-stock/v1/quotations/chk-holiday (get) : 국내휴장일조회
- [ ] /uapi/domestic-stock/v1/quotations/inquire-time-indexchartprice (get) : 업종분봉조회
- [ ] /uapi/domestic-stock/v1/quotations/inquire-vi-status (get) : 변동성완화장치(VI) 현황
- [ ] /uapi/domestic-stock/v1/quotations/inquire-index-tickprice (get) : 국내업종 시간별지수(초)
- [ ] /uapi/domestic-stock/v1/quotations/inquire-index-timeprice (get) : 국내업종 시간별지수(분)
- [ ] /uapi/domestic-stock/v1/quotations/exp-index-trend (get) : 국내주식 예상체결지수 추이
- [ ] /uapi/domestic-stock/v1/quotations/comp-interest (get) : 금리 종합(국내채권/금리)
- [ ] /uapi/domestic-stock/v1/quotations/news-title (get) : 종합 시황/공시(제목)
- [ ] /uapi/domestic-stock/v1/quotations/search-info (get) : 상품기본조회
- [ ] /uapi/domestic-stock/v1/quotations/search-stock-info (get) : 주식기본조회
- [ ] /uapi/domestic-stock/v1/finance/balance-sheet (get) : 국내주식 대차대조표
- [ ] /uapi/domestic-stock/v1/finance/income-statement (get) : 국내주식 손익계산서
- [ ] /uapi/domestic-stock/v1/finance/financial-ratio (get) : 국내주식 재무비율
- [ ] /uapi/domestic-stock/v1/finance/profit-ratio (get) : 국내주식 수익성비율
- [ ] /uapi/domestic-stock/v1/finance/other-major-ratios (get) : 국내주식 기타주요비율
- [ ] /uapi/domestic-stock/v1/finance/stability-ratio (get) : 국내주식 안정성비율
- [ ] /uapi/domestic-stock/v1/finance/growth-ratio (get) : 국내주식 성장성비율
- [ ] /uapi/domestic-stock/v1/quotations/credit-by-company (get) : 국내주식 당사 신용가능종목
- [ ] /uapi/domestic-stock/v1/ksdinfo/dividend (get) : 예탁원정보(배당일정)
- [ ] /uapi/domestic-stock/v1/ksdinfo/purreq (get) : 예탁원정보(주식매수청구일정)
- [ ] /uapi/domestic-stock/v1/ksdinfo/merger-split (get) : 예탁원정보(합병/분할일정)
- [ ] /uapi/domestic-stock/v1/ksdinfo/rev-split (get) : 예탁원정보(액면교체일정)
- [ ] /uapi/domestic-stock/v1/ksdinfo/cap-dcrs (get) : 예탁원정보(자본감소일정)
- [ ] /uapi/domestic-stock/v1/ksdinfo/list-info (get) : 예탁원정보(상장정보일정)
- [ ] /uapi/domestic-stock/v1/ranking/traded-by-company (get) : 국내주식 당사매매종목 상위
- [ ] /uapi/domestic-stock/v1/ksdinfo/paidin-capin (get) : 예탁원정보(유상증자일정)
- [ ] /uapi/domestic-stock/v1/ksdinfo/bonus-issue (get) : 예탁원정보(무상증자일정)
- [ ] /uapi/domestic-stock/v1/ksdinfo/sharehld-meet (get) : 예탁원정보(주주총회일정)
- [ ] /uapi/domestic-stock/v1/quotations/estimate-perform (get) : 국내주식 종목추정실적
- [ ] /uapi/domestic-stock/v1/quotations/lendable-by-company (get) : 당사 대주가능 종목
- [ ] /uapi/domestic-stock/v1/quotations/invest-opbysec (get) : 국내주식 증권사별 투자의견
- [ ] /uapi/domestic-stock/v1/quotations/invest-opinion (get) : 국내주식 종목투자의견
- [ ] /uapi/domestic-stock/v1/quotations/foreign-institution-total (get) : 국내기관_외국인 매매종목가집계
- [ ] /uapi/domestic-stock/v1/quotations/psearch-title (get) : 종목조건검색 목록조회
- [ ] /uapi/domestic-stock/v1/quotations/psearch-result (get) : 종목조건검색조회
- [ ] /uapi/domestic-stock/v1/quotations/program-trade-by-stock (get) : 종목별프로그램매매추이(체결)
- [ ] /uapi/domestic-stock/v1/quotations/program-trade-by-stock-daily (get) : 종목별 프로그램매매추이(일별)
- [ ] /uapi/domestic-stock/v1/quotations/investor-trend-estimate (get) : 종목별 외인기관 추정가집계
- [ ] /uapi/domestic-stock/v1/quotations/inquire-daily-trade-volume (get) : 종목별일별매수매도체결량
- [ ] /uapi/domestic-stock/v1/quotations/inquire-investor-time-by-market (get) : 시장별 투자자매매동향(시세)
- [ ] /uapi/domestic-stock/v1/quotations/inquire-investor-daily-by-market (get) : 시장별 투자자매매동향(일별)
- [ ] /uapi/domestic-stock/v1/quotations/daily-credit-balance (get) : 국내주식 신용잔고 일별추이
- [ ] /uapi/domestic-stock/v1/quotations/exp-price-trend (get) : 국내주식 예상체결가 추이
- [ ] /uapi/domestic-stock/v1/quotations/daily-short-sale (get) : 국내주식 공매도 일별추이
- [ ] /uapi/domestic-stock/v1/ranking/overtime-exp-trans-fluct (get) : 국내주식 시간외예상체결등락률
- [ ] /uapi/domestic-stock/v1/quotations/investor-program-trade-today (get) : 프로그램매매 투자자매매동향(당일)
- [ ] /uapi/domestic-stock/v1/quotations/comp-program-trade-today (get) : 프로그램매매 종합현황(시간)
- [ ] /uapi/domestic-stock/v1/quotations/comp-program-trade-daily (get) : 프로그램매매 종합현황(일별)
- [ ] /uapi/domestic-stock/v1/quotations/frgnmem-trade-estimate (get) : 외국계 매매종목 가집계
- [ ] /uapi/domestic-stock/v1/quotations/frgnmem-pchs-trend (get) : 종목별 외국계 순매수추이
- [ ] /uapi/domestic-stock/v1/quotations/tradprt-byamt (get) : 국내주식 체결금액별 매매비중
- [ ] /uapi/domestic-stock/v1/quotations/mktfunds (get) : 국내 증시자금 종합
- [ ] /uapi/domestic-stock/v1/quotations/intstock-grouplist (get) : 관심종목 그룹 조회
- [ ] /uapi/domestic-stock/v1/quotations/intstock-stocklist-by-group (get) : 관심종목 그룹별 종목조회
- [ ] /uapi/domestic-stock/v1/quotations/intstock-multprice (get) : 관심종목(멀티종목) 시세조회
- [ ] /uapi/domestic-stock/v1/quotations/capture-uplowprice (get) : 국내주식 상하한가 포착
- [ ] /uapi/domestic-stock/v1/quotations/frgnmem-trade-trend (get) : 회원사 실시간 매매동향(틱)
- [ ] /uapi/domestic-stock/v1/quotations/pbar-tratio (get) : 국내주식 매물대/거래비중
- [ ] /uapi/domestic-stock/v1/quotations/inquire-member-daily (get) : 주식현재가 회원사 종목매매동향
- [ ] /uapi/domestic-stock/v1/quotations/volume-rank (get) : 거래량순위
- [ ] /uapi/domestic-stock/v1/ranking/fluctuation (get) : 국내주식 등락률 순위
- [ ] /uapi/domestic-stock/v1/ranking/profit-asset-index (get) : 국내주식 수익자산지표 순위
- [ ] /uapi/domestic-stock/v1/ranking/market-cap (get) : 국내주식 시가총액 상위
- [ ] /uapi/domestic-stock/v1/ranking/finance-ratio (get) : 국내주식 재무비율 순위
- [ ] /uapi/domestic-stock/v1/ranking/after-hour-balance (get) : 국내주식 시간외잔량 순위
- [ ] /uapi/domestic-stock/v1/ranking/prefer-disparate-ratio (get) : 국내주식 우선주/괴리율 상위
- [ ] /uapi/domestic-stock/v1/ranking/quote-balance (get) : 국내주식 호가잔량 순위
- [ ] /uapi/domestic-stock/v1/ranking/disparity (get) : 국내주식 이격도 순위
- [ ] /uapi/domestic-stock/v1/ranking/market-value (get) : 국내주식 시장가치 순위
- [ ] /uapi/domestic-stock/v1/ranking/volume-power (get) : 국내주식 체결강도 상위
- [ ] /uapi/domestic-stock/v1/ranking/top-interest-stock (get) : 국내주식 관심종목등록 상위
- [ ] /uapi/domestic-stock/v1/ranking/exp-trans-updown (get) : 국내주식 예상체결 상승/하락상위
- [ ] /uapi/domestic-stock/v1/ranking/near-new-highlow (get) : 국내주식 신고/신저근접종목 상위
- [ ] /uapi/domestic-stock/v1/ranking/bulk-trans-num (get) : 국내주식 대량체결건수 상위
- [ ] /uapi/domestic-stock/v1/ranking/short-sale (get) : 국내주식 공매도 상위종목
- [ ] /uapi/domestic-stock/v1/ranking/credit-balance (get) : 국내주식 신용잔고 상위
- [ ] /uapi/domestic-stock/v1/ranking/dividend-rate (get) : 국내주식 배당률 상위
- [ ] /uapi/domestic-stock/v1/ranking/overtime-fluctuation (get) : 국내주식 시간외등락율순위
- [ ] /uapi/domestic-stock/v1/ranking/overtime-volume (get) : 국내주식 시간외거래량순위
- [ ] /uapi/domestic-futureoption/v1/trading/order (post) : 선물옵션 주문
- [ ] /uapi/domestic-futureoption/v1/trading/order-rvsecncl (post) : 선물옵션 정정취소주문
- [ ] /uapi/domestic-futureoption/v1/trading/inquire-ccnl (get) : 선물옵션 주문체결내역조회
- [ ] /uapi/domestic-futureoption/v1/trading/inquire-balance (get) : 선물옵션 잔고현황
- [ ] /uapi/domestic-futureoption/v1/trading/inquire-psbl-order (get) : 선물옵션 주문가능
- [ ] /uapi/domestic-futureoption/v1/trading/inquire-ngt-ccnl (get) : (야간)선물옵션 주문체결내역조회
- [ ] /uapi/domestic-futureoption/v1/trading/inquire-ngt-balance (get) : (야간)선물옵션 잔고현황
- [ ] /uapi/domestic-futureoption/v1/trading/inquire-psbl-ngt-order (get) : (야간)선물옵션 주문가능 조회
- [ ] /uapi/domestic-futureoption/v1/trading/inquire-balance-settlement-pl (get) : 선물옵션 잔고정산손익내역
- [ ] /uapi/domestic-futureoption/v1/trading/inquire-deposit (get) : 선물옵션 총자산현황
- [ ] /uapi/domestic-futureoption/v1/trading/inquire-balance-valuation-pl (get) : 선물옵션 잔고평가손익내역
- [ ] /uapi/domestic-futureoption/v1/trading/inquire-daily-amount-fee (get) : 선물옵션기간약정수수료일별
- [ ] /uapi/domestic-futureoption/v1/trading/inquire-ccnl-bstime (get) : 선물옵션 기준일체결내역
- [ ] /uapi/domestic-futureoption/v1/trading/ngt-margin-detail (get) : (야간)선물옵션증거금상세
- [ ] /uapi/domestic-futureoption/v1/quotations/inquire-price (get) : 선물옵션 시세
- [ ] /uapi/domestic-futureoption/v1/quotations/inquire-asking-price (get) : 선물옵션 시세호가
- [ ] /uapi/domestic-futureoption/v1/quotations/inquire-daily-fuopchartprice (get) : 선물옵션기간별시세(일/주/월/년)
- [ ] /uapi/domestic-futureoption/v1/quotations/inquire-time-fuopchartprice (get) : 선물옵션 분봉조회
- [ ] /uapi/domestic-futureoption/v1/quotations/exp-price-trend (get) : 선물옵션 일중예상체결추이
- [ ] /uapi/domestic-futureoption/v1/quotations/display-board-option-list (get) : 국내옵션전광판_옵션월물리스트
- [ ] /uapi/domestic-futureoption/v1/quotations/display-board-top (get) : 국내선물 기초자산시세
- [ ] /uapi/domestic-futureoption/v1/quotations/display-board-callput (get) : 국내옵션전광판_콜풋
- [ ] /uapi/domestic-futureoption/v1/quotations/display-board-futures (get) : 국내옵션전광판_선물
- [ ] /uapi/overseas-stock/v1/trading/order (post) : 해외주식 주문(베트남)
- [ ] /uapi/overseas-stock/v1/trading/order-rvsecncl (post) : 해외주식 정정취소주문(취소)
- [ ] /uapi/overseas-stock/v1/trading/order-resv (post) : 해외주식 예약주문접수(베트남)
- [ ] /uapi/overseas-stock/v1/trading/order-resv-ccnl (post) : 해외주식 예약주문접수 취소
- [ ] /uapi/overseas-stock/v1/trading/inquire-nccs (get) : 해외주식 미체결내역
- [ ] /uapi/overseas-stock/v1/trading/inquire-balance (get) : 해외주식 잔고
- [ ] /uapi/overseas-stock/v1/trading/inquire-ccnl (get) : 해외주식 주문체결내역
- [ ] /uapi/overseas-stock/v1/trading/inquire-present-balance (get) : 해외주식 체결기준 현재잔고
- [ ] /uapi/overseas-stock/v1/trading/order-resv-list (get) : 해외주식 예약주문조회
- [ ] /uapi/overseas-stock/v1/trading/inquire-psamount (get) : 해외주식 매수가능금액조회
- [ ] /uapi/overseas-stock/v1/trading/daytime-order (post) : 해외주식 미국주간주문(매도)
- [ ] /uapi/overseas-stock/v1/trading/daytime-order-rvsecncl (post) : 해외주식 미국주간정정취소(취소)
- [ ] /uapi/overseas-stock/v1/trading/inquire-period-profit (get) : 해외주식 기간손익
- [ ] /uapi/overseas-stock/v1/trading/foreign-margin (get) : 해외증거금 통화별조회
- [ ] /uapi/overseas-stock/v1/trading/inquire-period-trans (get) : 해외주식 일별거래내역
- [ ] /uapi/overseas-stock/v1/trading/inquire-paymt-stdr-balance (get) : 해외주식 결제기준잔고
- [ ] /uapi/overseas-price/v1/quotations/price (get) : 해외주식 현재체결가
- [ ] /uapi/overseas-price/v1/quotations/dailyprice (get) : 해외주식 기간별시세
- [ ] /uapi/overseas-price/v1/quotations/inquire-daily-chartprice (get) : 해외주식 종목/지수/환율기간별시세
- [ ] /uapi/overseas-price/v1/quotations/inquire-search (get) : 해외주식 조건검색
- [ ] /uapi/overseas-stock/v1/quotations/countries-holiday (get) : 해외결제일자조회
- [ ] /uapi/overseas-price/v1/quotations/price-detail (get) : 해외주식 현재가상세
- [ ] /uapi/overseas-price/v1/quotations/inquire-time-itemchartprice (get) : 해외주식분봉조회
- [ ] /uapi/overseas-price/v1/quotations/inquire-time-indexchartprice (get) : 해외지수분봉조회
- [ ] /uapi/overseas-price/v1/quotations/search-info (get) : 해외주식 상품기본정보
- [ ] /uapi/overseas-price/v1/quotations/inquire-asking-price (get) : 해외주식 현재가 10호가
- [ ] /uapi/overseas-price/v1/quotations/period-rights (get) : 해외주식 기간별권리조회
- [ ] /uapi/overseas-price/v1/quotations/news-title (get) : 해외뉴스종합(제목)
- [ ] /uapi/overseas-price/v1/quotations/rights-by-ice (get) : 해외주식 권리종합
- [ ] /uapi/overseas-price/v1/quotations/colable-by-company (get) : 당사 해외주식담보대출 가능 종목
- [ ] /uapi/overseas-price/v1/quotations/brknews-title (get) : 해외속보(제목)
- [ ] /uapi/overseas-futureoption/v1/trading/order (post) : 해외선물옵션 주문
- [ ] /uapi/overseas-futureoption/v1/trading/order-rvsecncl (post) : 해외선물옵션 정정취소주문(취소)
- [ ] /uapi/overseas-futureoption/v1/trading/inquire-ccld (get) : 해외선물옵션 당일주문내역
- [ ] /uapi/overseas-futureoption/v1/trading/inquire-unpd (get) : 해외선물옵션 미결제내역조회
- [ ] /uapi/overseas-futureoption/v1/trading/inquire-psamount (get) : 해외선물옵션 주문가능조회
- [ ] /uapi/overseas-futureoption/v1/trading/inquire-period-ccld (get) : 해외선물옵션 기간계좌손익 일별 Copy
- [ ] /uapi/overseas-futureoption/v1/trading/inquire-daily-ccld (get) : 해외선물옵션 일별체결내역
- [ ] /uapi/overseas-futureoption/v1/trading/inquire-deposit (get) : 해외선물옵션 예수금현황
- [ ] /uapi/overseas-futureoption/v1/trading/inquire-daily-order (get) : 해외선물옵션 일별 주문내역
- [ ] /uapi/overseas-futureoption/v1/trading/inquire-period-trans (get) : 해외선물옵션 기간계좌거래내역
- [ ] /uapi/overseas-futureoption/v1/trading/margin-detail (get) : 해외선물옵션 증거금상세
- [ ] /uapi/overseas-futureoption/v1/quotations/stock-detail (get) : 해외선물종목상세
- [ ] /uapi/overseas-futureoption/v1/quotations/inquire-price (get) : 해외선물종목현재가
- [ ] /uapi/overseas-futureoption/v1/quotations/inquire-time-futurechartprice (get) : 해외선물분봉조회
- [ ] /uapi/overseas-futureoption/v1/quotations/weekly-ccnl (get) : 해외선물 체결추이(주간)
- [ ] /uapi/overseas-futureoption/v1/quotations/daily-ccnl (get) : 해외선물 체결추이(일간)
- [ ] /uapi/overseas-futureoption/v1/quotations/tick-ccnl (get) : 해외선물 체결추이(틱)
- [ ] /uapi/overseas-futureoption/v1/quotations/monthly-ccnl (get) : 해외선물 체결추이(월간)
- [ ] /uapi/overseas-futureoption/v1/quotations/inquire-asking-price (get) : 해외선물 호가
- [ ] /uapi/overseas-futureoption/v1/quotations/search-contract-detail (get) : 해외선물 상품기본정보
- [ ] /uapi/overseas-futureoption/v1/quotations/market-time (get) : 해외선물 장운영시간
- [ ] /uapi/overseas-futureoption/v1/quotations/investor-unpd-trend (get) : 해외선물 미결제추이
- [ ] /uapi/overseas-futureoption/v1/quotations/opt-asking-price (get) : 해외옵션 호가
- [ ] /uapi/domestic-bond/v1/trading/buy (post) : 장내채권 매수주문
- [ ] /uapi/domestic-bond/v1/trading/sell (post) : 장내채권 매도주문
- [ ] /uapi/domestic-bond/v1/trading/order-rvsecncl (post) : 장내채권 정정취소주문
- [ ] /uapi/domestic-bond/v1/trading/inquire-psbl-rvsecncl (get) : 채권정정취소가능주문조회
- [ ] /uapi/domestic-bond/v1/trading/inquire-daily-ccld (get) : 장내채권 주문체결내역
- [ ] /uapi/domestic-bond/v1/trading/inquire-balance (get) : 장내채권 잔고조회
- [ ] /uapi/domestic-bond/v1/trading/inquire-psbl-order (get) : 장내채권 매수가능조회
- [ ] /uapi/domestic-bond/v1/quotations/issue-info (get) : 장내채권 발행정보
- [ ] /uapi/domestic-bond/v1/quotations/inquire-asking-price (get) : 장내채권현재가(호가)
- [ ] /uapi/domestic-bond/v1/quotations/avg-unit (get) : 장내채권 평균단가조회
- [ ] /uapi/domestic-bond/v1/quotations/inquire-daily-itemchartprice (get) : 장내채권 기간별시세(일)
- [ ] /uapi/domestic-bond/v1/quotations/inquire-price (get) : 장내채권현재가(시세)
- [ ] /uapi/domestic-bond/v1/quotations/inquire-ccnl (get) : 장내채권현재가(일별)
- [ ] /uapi/domestic-bond/v1/quotations/search-bond-info (get) : 장내채권 기본조회


## Reference
- [한국투자 OpenAPI](https://apiportal.koreainvestment.com/apiservice) - API문서
