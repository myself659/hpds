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
var len = flag.Int("len", 1000, "queue len")
var cpus = flag.Int("cpus", 1, "use cpu number")

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

func DoEnqueue(q *hpds.Queue, i int) {
	q.Enqueue(i)
}

func DoDequue(q *hpds.Queue) {
	n := *len
	for i := 0; i < n; i++ {
		value, ok := q.Dequeue()
		if ok == false {
			ch <- i

			return
		}
		//fmt.Println("dequeue:", value)
		value = value
	}
	ch <- n
}
func DoNew() *hpds.Queue {
	q := hpds.NewQueue()
	return q
}

func main() {
	flag.Parse()
	q := hpds.NewQueue()
	runtime.GOMAXPROCS(*cpus)

	for i := 0; i < *len; i++ {
		q.Enqueue(i)
	}

	fmt.Println("len:", q.Len())

	for g := 0; g < 10; g++ {
		go DoDequue(q)
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
