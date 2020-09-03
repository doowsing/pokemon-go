package captcha

import (
	"github.com/mojocn/base64Captcha"
)

// 返回：ID，答案，base64_str，error
func NewCaptcha() (string, string, string, error) {
	id, _, a := base64Captcha.DefaultDriverDigit.GenerateIdQuestionAnswer()
	item, err := base64Captcha.DefaultDriverDigit.DrawCaptcha(a)
	if err != nil {
		return "", "", "", err
	}
	return id, a, item.EncodeB64string(), nil
}
