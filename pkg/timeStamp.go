package pkg

import "time"

func UTCTransLocal(utcTime string) string {
	t, _ := time.Parse("2006-01-02T15:04:05.000Z", utcTime)
	//  转化为北京时间
	loc, _ := time.LoadLocation("Asia/Shanghai")
	bjTime := t.In(loc)
	//  格式化输出
	return bjTime.Format("2006-01-02 15:04:05")
}
