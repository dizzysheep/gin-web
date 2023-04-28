package log

import (
	"regexp"
	"strings"
)

//敏感变量名(大小写不敏感)
var SensitiveWords = []string{
	"password",
	"mobile",
	"phone",
	"card_no",
}

func HideSensitiveInfo(content string) string {
	//替换常见敏感变量名(大小写不敏感)
	content = doHideSensitiveParam(SensitiveWords, content)

	//识别并替换手机号
	reg := regexp.MustCompile(`[^\d](1\d{10})`)
	content = doHideStrSelf(reg, content, 1)

	//识别并替换身份证号
	reg = regexp.MustCompile(`[^\d](\d{18}X?)([^\d]|$)`)
	content = doHideStrSelf(reg, content, 1)

	return content
}

func doHideSensitiveParam(SensitiveWords []string, content string) string {
	for _, sensitiveWord := range SensitiveWords {
		reg := regexp.MustCompile("(?i)" + sensitiveWord + `['"]?(=|\:"|\:'|:)` + `([^$&,"':]*?)` + `[$|&|,|"|']`)
		content = doHideStrSelf(reg, content, 2)
	}
	return content
}

func doHideStrSelf(reg *regexp.Regexp, content string, subMatchIndex int) string {
	if reg.MatchString(content) {
		allMatches := reg.FindAllStringSubmatch(content, -1)
		if len(allMatches) == 0 {
			return content
		}
		oriStrSlice := []string{}
		for _, sub := range allMatches {
			if len(sub) > subMatchIndex {
				oriStrSlice = append(oriStrSlice, sub[subMatchIndex])
			}
		}
		hideStrSlice := []string{}
		for _, oriStr := range oriStrSlice {
			hideStrSlice = append(hideStrSlice, getHideStr(oriStr))
		}
		if len(oriStrSlice) != len(hideStrSlice) {
			return content
		}
		for key, oriStr := range oriStrSlice {
			content = strings.Replace(content, oriStr, hideStrSlice[key], -1)
		}
	}
	return content
}

func getHideStr(str string) string {
	length := len(str)
	if length == 0 {
		return ""
	}
	hideLen := 0
	startIndex := 0
	if len(str) >= 8 {
		hideLen = len(str) - 6
		startIndex = 3
	} else {
		hideLen = len(str) - 2
		startIndex = 1
	}
	if hideLen < 0 {
		hideLen = 0
	}

	frontPart := append([]byte(str[:startIndex]), []byte(`\*\*\*\*`)...)
	return string(append(frontPart, []byte(str[startIndex+hideLen:])...))
}
