package timeutil

import (
	"strconv"
	"time"
)

var NowFunc = time.Now

// YYYYMM 形式の文字列を fiscal year 形式にする
func ConvertToFiscalYear(yyyymm string) string {
	year, _ := strconv.Atoi(yyyymm[:4])
	month, _ := strconv.Atoi(yyyymm[4:])
	if month >= 1 && month <= 3 {
		return strconv.Itoa(year - 1)
	}
	return strconv.Itoa(year)
}
