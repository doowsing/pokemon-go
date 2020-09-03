package main

import "fmt"

func main() {
	ids := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ids <- i
		}
		close(ids)
	}()
	for id := range ids {
		fmt.Printf("id:%d\n", id)
	}
}
