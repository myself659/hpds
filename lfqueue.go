package hpds

import (
	"sync/atomic"
	"unsafe"
)

type node struct {
	value interface{}
	next  *node
}

type Queue struct {
	dummy *node
	tail  *node //指向尾节点
	len   int32 // 长度
}

func NewQueue() *Queue {
	q := new(Queue)

	q.dummy = new(node)
	q.tail = q.dummy /*  新建节点作为尾节点  */
	q.len = 0

	return q
}

func (q *Queue) Enqueue(v interface{}) {
	var oldTail, oldTailNext *node

	newNode := new(node)
	newNode.value = v

	newNodeAdd := false

	for !newNodeAdd {
		oldTail = q.tail
		oldTailNext = oldTail.next

		if q.tail != oldTail { //队列尾部有改变，返回继续，第一条件：队列尾部指针不被修改
			continue
		}

		if oldTailNext != nil { // 有其他新增节点已经添加，但是没有完成尾节点更新操作
			/*
				q.tail = oldTailNext;
				尝试更新尾节点
			*/
			atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)),
				unsafe.Pointer(oldTail),
				unsafe.Pointer(oldTailNext))
			continue
		}
		/*
			类似c语言操作 oldTail.next = newNode  新节点添加到尾节点
			这个操作必须保证成功
		*/
		newNodeAdd = atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&oldTail.next)),
			unsafe.Pointer(oldTailNext),
			unsafe.Pointer(newNode))

	}
	/*
		不判断返回值，能保证成功吗?  不需要保证成功，如果不成功，有其他节点操作，这时候该节点并不是尾节点
		q.tail = newNode;
		添加节点作为新尾部
		尝试更新尾节点
	*/

	atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)),
		unsafe.Pointer(oldTail),
		unsafe.Pointer(newNode))
	/* 长度加1 */
	atomic.AddInt32(&q.len, int32(1))
}

func (q *Queue) Dequeue() (interface{}, bool) {

	var temp interface{}
	var oldDummy, oldHead *node
	var ppDummy **node = &q.dummy
	removed := false

	for !removed {
		oldDummy = q.dummy // data race
		oldHead = oldDummy.next
		oldTail := q.tail

		if q.dummy != oldDummy { /* 头节点被修改 */
			continue
		}

		if oldHead == nil { /* 队列为空 */
			return nil, false
		}

		if oldTail == oldDummy {
			/*
				更新尾节点
				在只有一个节点条件下，其他两个线程一个添加节点，一个删除节点会出现这种情况
			*/
			atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)), unsafe.Pointer(oldTail), unsafe.Pointer(oldHead))
			continue
		}

		temp = oldHead.value
		/*
			从队列中删除首节点
			 q.dummy =  q.dummy.next
			 更新原来dummy指针的值
			 data race
		*/
		removed = atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(ppDummy)), unsafe.Pointer(oldDummy), unsafe.Pointer(oldHead))

	}
	/* 长度减1 */
	atomic.AddInt32(&q.len, int32(-1))

	return temp, true
}

func (q *Queue) iterate(c chan<- interface{}) {

	for {
		value, ok := q.Dequeue()
		if !ok {
			break
		}

		c <- value
	}

	close(c)
}

func (q *Queue) Iter() <-chan interface{} {
	c := make(chan interface{})
	go q.iterate(c)

	return c
}

func (q *Queue) Len() int32 {
	return atomic.LoadInt32(&q.len)
}
