package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/unknwon/com"
	"io/ioutil"
	"net/http"
	"strings"
)

type Info struct {
	Id       int
	Name     string
	Children []Info
}

func GetData() ([]byte, error) {
	// 获取json数据
	data, err := http.Get("https://job.xiyanghui.com/api/q1/json")
	if err != nil {
		return nil, err
	}
	defer data.Body.Close()
	b, err := ioutil.ReadAll(data.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func IterInfos(id int, infoSlice []Info, saveSlice *[]string) bool {
	// 对数据进行多层迭代
	for _, info := range infoSlice {
		//fmt.Printf("name :%s\n", info.Name)
		if id == info.Id || IterInfos(id, info.Children, saveSlice) {
			*saveSlice = append(*saveSlice, info.Name)
			return true
		}
	}
	return false
}

func PrintStructNames(id int) {
	data, err := GetData()
	if err != nil {
		fmt.Printf("获取数据失败！err:%s\n", err.Error())
	}
	var infos []Info
	if err = json.Unmarshal(data, &infos); err != nil {
		fmt.Printf("格式化数据失败！！err:%s\n", err.Error())
	}
	structNames := []string{}

	// 填充结果切片
	IterInfos(id, infos, &structNames)

	// 对结果切片进行倒置
	for i := 0; i < len(structNames)/2; i++ {
		structNames[i], structNames[len(structNames)-i-1] = structNames[len(structNames)-i-1], structNames[i]
	}
	fmt.Printf("层级结构:%s\n", strings.Join(structNames, " > "))

}

func GetRateData() ([]byte, error) {
	// 获取json数据
	data, err := http.Get("https://app-cdn.2q10.com/api/v2/currency")
	if err != nil {
		return nil, err
	}
	defer data.Body.Close()
	b, err := ioutil.ReadAll(data.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

type Currency struct {
	Name string
	Code string
	Rate float64 // base on some currency
}

var currencyEnum = []Currency{
	{Name: "USD", Code: "$", Rate: 0},
	{Name: "GBP", Code: "£", Rate: 0},
	{Name: "EUR", Code: "€", Rate: 0},
	{Name: "HKD", Code: "HK$", Rate: 0},
	{Name: "JPY", Code: "¥", Rate: 0},
}

func ExchangeCurrencyString(money string) (float64, error) {
	data, err := GetRateData()
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Get rate data failed, error:%s\n", err.Error()))
	}
	rateData := &struct {
		Rates        map[string]float64
		Base         string
		Last_data_at int
	}{}
	err = json.Unmarshal(data, rateData)
	if err != nil {

		return 0, errors.New(fmt.Sprintf("Decode rate data failed, error:%s\n", err.Error()))
	}
	var baseRateUSD2RMB float64 = 0
	for name, rate := range rateData.Rates {
		if name == "CNH" {
			baseRateUSD2RMB = rate
		}
		for i, _ := range currencyEnum {
			if currencyEnum[i].Name == name {
				currencyEnum[i].Rate = rate
			}
		}
	}
	if baseRateUSD2RMB == 0 {
		return 0, errors.New(fmt.Sprint("Never find CNH Rate!"))
	}
	for _, currency := range currencyEnum {
		finded := false
		if strings.Index(money, currency.Name) > -1 {
			money = money[len(currency.Name):]
			finded = true
		} else if strings.Index(money, currency.Code) > -1 {
			money = money[len(currency.Code):]
			finded = true
		}
		if finded {
			if currency.Rate == 0 {
				return 0, errors.New(fmt.Sprintf("Never find %s Rate!", currency.Name))
			}
			money = strings.ReplaceAll(money, ",", "")
			moneyData, err := com.StrTo(money).Float64()
			fmt.Printf("获取汇率：%f, 输入钱数%f\n", currency.Rate, moneyData)
			if err != nil {
				return 0, errors.New(fmt.Sprint("The input is not number!"))
			}
			return moneyData / currency.Rate * baseRateUSD2RMB, nil
		}

	}
	return 0, errors.New("Get unknown currency!")
}

func main1() {
	var money string
	fmt.Println("请输入外币符号与数理，按回车键结束:")
	fmt.Scanf("%s", &money)
	if money2CHY, err := ExchangeCurrencyString(money); err != nil {
		panic(err)
	} else {
		fmt.Printf("转换为人民币数量为:%.2f CHY", money2CHY)
	}
}

//func main() {
//	var id int
//	fmt.Println("请输入数据的整数id，按回车键结束:")
//	_, err := fmt.Scanf("%d", &id)
//	if err != nil {
//		fmt.Printf("输入数据非整数！")
//	}
//	PrintStructNames(id)
//}
