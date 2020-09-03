package utils

import (
	"bytes"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/unknwon/com"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

func ToGbk(str string) string {
	enc := mahonia.NewEncoder("gbk")
	return enc.ConvertString(str)
}
func ToUtf8(str string) string {
	//enc := mahonia.NewEncoder("utf-8")
	//return enc.ConvertString(str)
	I := bytes.NewReader([]byte(str))
	O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return ""
	}
	return string(d)
}

func ToQuote(str string) string {
	return strconv.Quote(str)
}

func IntToStr(num int) string {
	return strconv.Itoa(num)
}

func ToFloat64(str string) float64 {
	flag := false
	if strings.Index(str, "%") > -1 {
		str = strings.ReplaceAll(str, "%", "")
		flag = true
	}
	float := com.StrTo(str).MustFloat64()
	if flag {
		float /= 100.0
	}
	return float
}
func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}
func FloatToStr(num float64, accurate int) string {
	return fmt.Sprintf("%"+fmt.Sprintf(".%df", accurate), num)
}

func CzlStr(czl float64) string {
	return fmt.Sprintf("%g", com.StrTo(fmt.Sprintf("%.1f", czl)).MustFloat64())
}

func RandInt(start, end int) int {
	return rand.Intn(end-start+1) + start
}
