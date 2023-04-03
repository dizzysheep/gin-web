package xtime

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FMT_TYPE_NOMAL
const (
	DATE_TIME_FMT = "2006-01-02 15:04:05"

	DATE_MINUTE_FMT = "2006-01-02 15:04"

	DATE_MONTH_FMT = "2006-01"

	DATE_TIME_FMT_ISO8601 = "2006-01-02T15:04:05Z"

	DATE_FMT = "2006-01-02"

	TIME_FMT = "15:04:05"

	DATE_TIME_FMT_CN = "2006年01月02日 15时04分05秒"

	DATE_FMT_CN = "2006年01月02日"

	TIME_FMT_CN = "15时04分05秒"

	DAY_FMT_CN = "01月02日"

	DAY_FMT = "15:04"

	DEFAULE_TIME_FMT = "1111-11-11 11:11:11"
)

const SecondInNano = 1000 * 1000 * 1000

//return 1441006057 in sec
func GetTimestamp() int64 {
	return time.Now().Unix()
}

//return 1441006057 in sec
func GetTimestampString() string {
	return strconv.FormatInt(GetTimestamp(), 10)
}

// return 1441007112776 in millisecond
func GetTimestampInMilli() int64 {
	return int64(time.Now().UnixNano() / (1000 * 1000)) // ms
}

// return 1441007112776 in millisecond
func GetTimestampInMilliString() string {
	return strconv.FormatInt(GetTimestampInMilli(), 10)
}

//微秒
func GetTimestampInMicro() int64 {
	return int64(time.Now().UnixNano() / 1000) // ms
}

// 微秒
func GetTimestampInMicroString() string {
	return strconv.FormatInt(GetTimestampInMicro(), 10)
}

//format
func GetCurrentTimeFormat(format string) string {
	return GetTimeFormat(GetTimestamp(), format)
}

//
func GetTimeFormat(second int64, format string) string {
	return time.Unix(second, 0).Format(format)
}

// Timing the cost of function call, unix nano was returned
func Elapse(f func()) int64 {
	now := time.Now().UnixNano()
	f()
	return time.Now().UnixNano() - now
}

// Timing the cost of function call, unix nano was returned
func ElapseString(f func()) string {
	return strconv.FormatInt(Elapse(f), 10)
}

// GetMonthDays return days of the month/year
func GetMonthDays(year, month int) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if IsLeapYear(year) {
			return 29
		}
		return 28
	default:
		panic(fmt.Sprintf("Illegal month:%d", month))
	}
}

// IsLeapYear check whether a year is leay
func IsLeapYear(year int) bool {
	if year%100 == 0 {
		return year%400 == 0
	}

	return year%4 == 0
}

// ToDayStart .
func ToDayStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

// ToDayEnd .
func ToDayEnd(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999, time.Local)
}

//GetBirthday brith 2019-11-07
func Birthday(birth string) (int, error) {
	birthSlice := strings.Split(birth, "-")

	if len(birthSlice) != 3 {
		return 0, errors.New("出生日期格式错误")
	}
	birYear, _ := strconv.Atoi(birthSlice[0])
	birMonth, _ := strconv.Atoi(birthSlice[1])
	birDay, _ := strconv.Atoi(birthSlice[2])
	age := time.Now().Year() - birYear
	if int(time.Now().Month()) < birMonth { //月份小于出生月份，岁数减一
		age--
	}
	nowDay, _ := strconv.Atoi(time.Now().Format("02"))
	if int(time.Now().Month()) == birMonth && nowDay < birDay { //月份与出生月份相同，日期小于出生日期，岁数减一
		age--
	}

	return age, nil
}

//时间转换
//t time.Time 要转换的时间
//返回结果：
//1、近三天内，显示“今天”、“昨天”、“前天”+时间，例如“前天17:59”；
//2、今年内显示日期，例如“11月5日”；
//3、去年及以前则显示年月日，例如“2018-10-21”
func StrTime(t time.Time) string {
	//默认时间返回空
	if t.Format(DATE_TIME_FMT) == DEFAULE_TIME_FMT {
		return ""
	}
	day := t.Format(DATE_FMT)
	//今天
	today := time.Now().AddDate(0, 0, 0).Format(DATE_FMT)
	if today == day {
		return "今天" + t.Format(DAY_FMT)
	}
	//昨天
	yesterDay := time.Now().AddDate(0, 0, -1).Format(DATE_FMT)
	if yesterDay == day {
		return "昨天" + t.Format(DAY_FMT)
	}
	//前天
	beforeYesterDay := time.Now().AddDate(0, 0, -2).Format(DATE_FMT)
	if beforeYesterDay == day {
		return "前天" + t.Format(DAY_FMT)
	}
	//今年
	year := time.Now().Year()
	if year == t.Year() {
		return t.Format(DAY_FMT_CN)
	}
	//去年以前
	lastYear := time.Now().Year() - 1
	if lastYear >= t.Year() {
		return t.Format(DATE_FMT)
	}
	return ""
}

//时间转换
//timeStr string  时间戳字符串
//返回结果：
//1、近三天内，显示 今天/昨天/前天
//2、非近三天，则返回空
func WithinThreeDays(timeStr string) string {
	// 计算今天最大时间戳 2020-03-03 23:59:59
	todayDate := time.Now().Format(DATE_FMT)
	todayTime, _ := time.Parse(DATE_FMT, todayDate)
	todayUnix := todayTime.Unix() - 8*60*60 + 60*60*24 - 1
	// 获取参数字符串
	timeParamsInt64, _ := strconv.ParseInt(timeStr, 10, 64)
	// 判断时间差
	timeDiff := (todayUnix - timeParamsInt64) / 86400
	if timeDiff == 0 {
		return "今天"
	} else if timeDiff == 1 {
		return "昨天"
	} else if timeDiff == 2 {
		return "前天"
	} else {
		return ""
	}
}
