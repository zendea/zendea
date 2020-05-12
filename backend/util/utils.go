package util

import (
	"errors"
	"regexp"
	"strings"

	"github.com/spf13/viper"

	"zendea/model"
	"zendea/util/markdown"
	"zendea/util/strtrim"
)

// 是否是正式环境
func IsProd() bool {
	return viper.GetString("mode") == "release"
}

// index of
func IndexOf(userIds []int64, userId int64) int {
	if userIds == nil || len(userIds) == 0 {
		return -1
	}
	for i, v := range userIds {
		if v == userId {
			return i
		}
	}
	return -1
}

func GetSummary(contentType string, content string) (summary string) {
	if contentType == model.ContentTypeMarkdown {
		mr := markdown.NewMd().Run(content)
		summary = mr.SummaryText
	} else {
		summary = strtrim.GetTextSummary(strtrim.GetHtmlText(content), 256)
	}
	return
}

// 截取markdown摘要
func GetMarkdownSummary(text string) string {
	if len(text) == 0 {
		return ""
	}
	mdResult := markdown.NewMd().Run(text)
	return mdResult.SummaryText
}

// 获取html内容摘要
func GetHtmlSummary(html string) string {
	if len(html) == 0 {
		return ""
	}
	text := strtrim.GetHtmlText(html)
	return strtrim.GetTextSummary(text, 256)
}

// 验证用户名合法性，用户名必须由5-12位(数字、字母、_、-)组成，且必须以字母开头。
func IsValidateUsername(username string) error {
	if len(username) == 0 {
		return errors.New("请输入用户名")
	}
	matched, err := regexp.MatchString("^[0-9a-zA-Z_-]{5,12}$", username)
	if err != nil || !matched {
		return errors.New("用户名必须由5-12位(数字、字母、_、-)组成，且必须以字母开头")
	}
	matched, err = regexp.MatchString("^[a-zA-Z]", username)
	if err != nil || !matched {
		return errors.New("用户名必须由5-12位(数字、字母、_、-)组成，且必须以字母开头")
	}
	return nil
}

// 验证是否是合法的邮箱
func IsValidateEmail(email string) (err error) {
	if len(email) == 0 {
		err = errors.New("邮箱格式不符合规范")
		return
	}
	pattern := `^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$`
	matched, _ := regexp.MatchString(pattern, email)
	if !matched {
		err = errors.New("邮箱格式不符合规范")
	}
	return
}

// 是否是合法的密码
func IsValidatePassword(password, rePassword string) error {
	if len(password) == 0 {
		return errors.New("请输入密码")
	}
	if strtrim.RuneLen(password) < 6 {
		return errors.New("密码过于简单")
	}
	if password != rePassword {
		return errors.New("两次输入密码不匹配")
	}
	return nil
}

// 是否是合法的URL
func IsValidateUrl(url string) error {
	if len(url) == 0 {
		return errors.New("URL格式错误")
	}
	indexOfHttp := strings.Index(url, "http://")
	indexOfHttps := strings.Index(url, "https://")
	if indexOfHttp == 0 || indexOfHttps == 0 {
		return nil
	}
	return errors.New("URL格式错误")
}

// 是否是内部图片
func IsInternalImage(imageUrl string) bool {
	// TODO @ 2019/12/31 这个地方硬编码了，要修改
	return strings.Contains(imageUrl, "file.zendea.com") || strings.Contains(imageUrl, "static.zendea.com")
}

// 应用图片样式
func ApplyImageStyle(imageUrl, styleName string) string {
	if !IsInternalImage(imageUrl) {
		return imageUrl
	}
	splitterIndex := strings.LastIndex(imageUrl, "!")
	if splitterIndex >= 0 {
		imageUrl = imageUrl[:splitterIndex]
	}
	return imageUrl + "!" + styleName
}

// 处理标签
func ParseTagsToArray(tags string) []string {
	if len(tags) == 0 {
		return nil
	}
	ss := strings.Split(tags, ",")
	if len(ss) == 0 {
		return nil
	}
	var ret []string
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			continue
		}
		ret = append(ret, s)
	}
	return ret
}
