package utils

import (
	"errors"
	"time"
)

// DateLayout 日期（不含时间）的解析/格式化布局
const DateLayout = "2006-01-02"

// ParseDateToTimestamp 将 YYYY-MM-DD 格式的日期字符串解析为 Unix 秒时间戳。
// value 为空时返回 nil, nil，便于可选筛选参数直接透传；
// endOfDay 为 true 时取次日零点，用于表示区间上界（闭区间）。
func ParseDateToTimestamp(value string, endOfDay bool) (*int64, error) {
	if value == "" {
		return nil, nil
	}
	day, err := time.ParseInLocation(DateLayout, value, time.Local)
	if err != nil {
		return nil, errors.New("日期格式必须为 YYYY-MM-DD")
	}
	if endOfDay {
		day = day.AddDate(0, 0, 1)
	}
	timestamp := day.Unix()
	return &timestamp, nil
}
