package main

import "fmt"

const base_time_format = "2006-01-02 15:04:05"

func main() {
	//passwd := "23692828"
	//ok, err := regexp.MatchString(`^[A-Za-z\d$@!%*?&.]{6, 16}`, passwd)
	//if err != nil {
	//	log.Printf("密码格式出错：%s\n", err)
	//	return
	//}
	//if !ok {
	//	log.Printf("密码不通过。\n")
	//} else {
	//	log.Printf("密码通过。\n")
	//}
	datas := make(map[int][]int)
	data, ok := datas[0]
	fmt.Print(ok, data == nil)

}
