package main

import (
	"fmt"
	"sync"
)

var notifyA = make(chan struct{})
var notifyB = make(chan struct{})
var s = "abcdefg"
var wg sync.WaitGroup
var q = make(chan string, len(s))
var done = make(chan struct{})

func PrintA() {
	defer wg.Done()

	for v := range q {
		<-notifyA
		fmt.Println("A:", v)
		notifyB <- struct{}{}
	}
}

func PrintB() {
	defer wg.Done()
	for v := range q {
		<-notifyB
		fmt.Println("B:", v)
		notifyA <- struct{}{}
	}
}

func main() {
	wg.Add(2)
	for _, c := range s {
		q <- string(c)
	}
	close(q)
	go PrintA()
	go PrintB()
	notifyA <- struct{}{}
	wg.Wait()
}