package strtrim

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 截取字符串
func Substr(s string, start, length int) string {
	bt := []rune(s)
	if start < 0 {
		start = 0
	}
	if start > len(bt) {
		start = start % len(bt)
	}
	var end int
	if (start + length) > (len(bt) - 1) {
		end = len(bt)
	} else {
		end = start + length
	}
	return string(bt[start:end])
}

// 字符串长度
func RuneLen(s string) int {
	bt := []rune(s)
	return len(bt)
}

// 获取summary
func GetTextSummary(s string, length int) string {
	summary := Substr(s, 0, length)
	if RuneLen(s) > length {
		summary += "..."
	}
	return summary
}

// 获取html文本
func GetHtmlText(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	return doc.Text()
}
