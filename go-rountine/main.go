package main

import (
	"fmt"
	"time"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}

	c <- sum // send to channel
}

func main() {
	/*
	go say("world")
	say("hello")
	*/

	/*
	s := []int{7, 2, 8, -9, 4, 0}

	c := make(chan int)
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)

	x, y := <-c, <-c // receive from channel

	fmt.Println(x, y, x+y)
	*/

	// ANONYMUS FUNCTION
	for _, v := range []int{1, 2, 3, 4, 5} {
		go func(i int) {
			fmt.Println(i)
		}(v)
	}
}