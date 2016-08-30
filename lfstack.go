package hpds

import (
	"sync/atomic"
	"unsafe"
)

type stacknode struct {
	value interface{}
	next  *stacknode
}

type Stack struct {
	head *stacknode
	len  int32 // 长度
}

func NewStack() *Stack {
	stack := new(Stack)
	stack.head = nil
	stack.len = 0
	return stack
}

func (self *Stack) Push(value interface{}) {
	newHead := new(stacknode)
	newHead.value = value
	ok := false
	for {
		oldHead := self.head
		newHead.next = oldHead
		ok = atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&self.head)),
			unsafe.Pointer(oldHead),
			unsafe.Pointer(newHead))
		if ok == true {
			atomic.AddInt32(&self.len, 1)
			return
		}
	}

}

func (self *Stack) Pop() interface{} {
	ok := false
	for {
		oldHead := self.head
		if oldHead == nil {
			return nil
		}
		newHead := oldHead.next
		ok = atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&self.head)),
			unsafe.Pointer(oldHead),
			unsafe.Pointer(newHead))
		if ok == true {
			atomic.AddInt32(&self.len, -1)
			return oldHead.value
		}

	}

}

func (self *Stack) Len() int32 {
	return atomic.LoadInt32(&self.len)
}
