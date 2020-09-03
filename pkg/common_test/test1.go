package main

import (
	"fmt"
	"sync"
)

type student struct {
	name string
}

type class struct {
	name string
	s1   *student
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

func DropStudent(s *student) {
	spool.Put(s)
}

func NewStudent() *student {
	s := spool.Get().(*student)
	s.name = ""
	return s
}

var cpool = &sync.Pool{New: func() interface{} {
	return &class{}
}}

var spool = &sync.Pool{New: func() interface{} {
	return &student{}
}}

func main() {
	students := []*student{
		{name: "张瑞"}, {"张新"},
	}
	newStudents := []*student{}
	for _, v := range students {
		fmt.Printf("students:%p\n", v)
		newStudents = append(newStudents, v)
	}
	for _, v := range newStudents {
		fmt.Printf("newstudent:%p\n", v)
	}
}
