package kinvest

import "time"

// IsKRXTradable checks if the given time is tradable on KRX (Korea Exchange).
func IsKRXTradable(t time.Time) (bool, string) {
	// 한국 시간으로 변환 (서울 기준)
	loc, _ := time.LoadLocation("Asia/Seoul")
	t = t.In(loc)

	// 요일 검사
	weekday := t.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false, "주말 - 거래 불가"
	}

	// 시간 파싱
	hour, min := t.Hour(), t.Minute()
	totalMin := hour*60 + min

	switch {
	case totalMin >= 510 && totalMin < 520:
		return true, "장전 시간외 종가매매 (08:30~08:40)"
	case totalMin >= 520 && totalMin < 540:
		return false, "동시호가 주문 접수 중 (08:40~09:00) - 체결은 09:00"
	case totalMin >= 540 && totalMin < 930:
		return true, "정규장 거래 가능 (09:00~15:30)"
	case totalMin >= 940 && totalMin < 960:
		return true, "장후 시간외 종가매매 (15:40~16:00)"
	case totalMin >= 960 && totalMin < 1080:
		return true, "시간외 단일가 매매 (16:00~18:00)"
	default:
		return false, "장외시간 - 거래 불가"
	}
}
