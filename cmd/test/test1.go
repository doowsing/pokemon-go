package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"
)

type Student struct {
	Name string
}

type class struct {
	name string
	s1   *Student
}

func NewClass() *class {
	c := cpool.Get().(*class)
	c.s1 = NewStudent()
	return c
}

func DropClass(c *class) {
	DropStudent(c.s1)
	cpool.Put(c)
}

func DropStudent(s *Student) {
	spool.Put(s)
}

func NewStudent() *Student {
	s := spool.Get().(*Student)
	s.Name = ""
	return s
}

var cpool = &sync.Pool{New: func() interface{} {
	return &class{}
}}

var spool = &sync.Pool{New: func() interface{} {
	return &Student{}
}}

func unmarshal(bytes []byte, r interface{}) error {
	return json.Unmarshal(bytes, r)
}

func GetTime() {
	time.Now()
}

func DeleteSlice(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}
func threeSum(nums []int) [][]int {
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})
	fmt.Println(nums)
	var result [][]int
	a, b, c := 0, 1, len(nums)-1
	for b < c && nums[a] <= 0 {

		for b < c {
			s := nums[a] + nums[b] + nums[c]
			if s == 0 {
				find := false
				if len(result) > 0 {
					l := len(result)
					if result[l-1][0] == nums[a] && result[l-1][1] == nums[b] && result[l-1][2] == nums[c] {
						find = true
					}
				}
				if !find {
					result = append(result, []int{nums[a], nums[b], nums[c]})
					fmt.Println(result)
				}
				if b-a > 1 {
					for a < b {
						if nums[a] == nums[a+1] {
							a += 1
						} else {
							break
						}
					}
				}
				b += 1
				c -= 1
			} else if s < 0 {
				b += 1
			} else {
				c -= 1
			}

		}
		a += 1
		b = a + 1
		c = len(nums) - 1
	}
	return result
}

func threeSumClosest(nums []int, target int) int {
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})
	a, b, c := 0, 1, len(nums)-1
	score := nums[a] + nums[b] + nums[c]
	deS := target - score
	if deS < 0 {
		deS = -deS
	}
	for b < c {
		for b < c {
			s := nums[a] + nums[b] + nums[c]
			de := target - s
			if de == 0 {
				return s
			} else {
				des := de
				if des < 0 {
					des = -des
				}
				if des < deS {
					deS = des
					score = s
				}
				if de < 0 {
					c -= 1
				} else {
					b += 1
				}
			}
		}
		a += 1
		b = a + 1
		c = len(nums) - 1
	}
	return score
}
func fourSum(nums []int, target int) [][]int {
	var result [][]int

	if len(nums) < 4 {
		return result
	}
	a, b, c, d := 0, 1, 2, len(nums)-1
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})
	for c < d {
		for c < d {
			for c < d {
				s := nums[a] + nums[b] + nums[c] + nums[d]
				if s > target {
					d -= 1
				} else if s < target {
					c += 1
				} else {
					find := false
					if len(result) > 0 {
						l := len(result)
						if result[l-1][0] == nums[a] && result[l-1][1] == nums[b] && result[l-1][2] == nums[c] && result[l-1][3] == nums[d] {
							find = true
						}
					}
					if !find {
						result = append(result, []int{nums[a], nums[b], nums[c], nums[d]})
					}
					c += 1
					d -= 1
				}
			}

			b += 1
			c = b + 1
			d = len(nums) - 1
		}
		a += 1
		b = a + 1
		c = b + 1
		d = len(nums) - 1
	}
	return result
}

func GetChatLoginString() {
	ptoken := "Nf+bAwEBClZlcmlmeUluZm8B/5wAAQMBAklkAQQAAQdBY2NvdW50AQwAAQVUb2tlbgEMAAAAF/+cAf4DBgEEbWFzawEIeu8MfopO01IA"
	data, _ := json.Marshal(&map[string]interface{}{
		"type": "login",
		"data": ptoken,
	})
	fmt.Printf("%s", data)
}

func RecordTime(taskName string) func() {
	now := time.Now()
	return func() {
		fmt.Printf("task %s spend time :%.3f s", taskName, time.Now().Sub(now).Seconds())
	}

}
func main() {
	GetA()
}
