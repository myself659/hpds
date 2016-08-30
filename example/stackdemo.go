package main

import (
	"flag"
	"fmt"
	"github.com/myself659/hpds"
	_ "math/rand"
	_ "os"
	"runtime"
	_ "sync"
	"time"
)

var ch = make(chan int)
var len = flag.Int("len", 1000, "stack len")
var cpus = flag.Int("cpus", 2, "use cpu number")

func Counter() {
	sum := 0
	for {
		c, ok := <-ch
		if ok == true {
			sum += c
		} else {
			fmt.Println("sum:", sum)
			return
		}

	}
}

func DoPush(stack *hpds.Stack, i int) {
	stack.Push(i)
}

func DoPop(stack *hpds.Stack) {
	n := *len
	for i := 0; i < n; i++ {
		value := stack.Pop()
		if value == nil {
			ch <- i

			return
		}
		fmt.Println("Pop:", value)
		value = value
	}
	ch <- n
}
func DoNew() *hpds.Stack {
	stack := hpds.NewStack()
	return stack
}

func main() {
	flag.Parse()
	stack := hpds.NewStack()
	runtime.GOMAXPROCS(*cpus)

	for i := 0; i < *len; i++ {
		stack.Push(i)
	}

	fmt.Println("len:", stack.Len())

	for g := 0; g < 10; g++ {
		go DoPop(stack)
	}

	sum := 0
	for g := 0; g < 10; g++ {
		c, ok := <-ch
		if ok == true {
			sum += c
		}

	}
	fmt.Println("sum:", sum)

	<-time.After(1 * time.Second)

}
