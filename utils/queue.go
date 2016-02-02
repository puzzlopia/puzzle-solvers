package utils

import "fmt"

type T interface{}

type PriorityQueue interface {
	PushBack(T)
	PopFront() T
}

type node struct {
	prev_  *node
	next_  *node
	value_ T
}

type Queue struct {
	sttNode_ *node
	endNode_ *node
	size_    int
}

func (q *Queue) PushBack(v T) {
	n := &node{q.endNode_, nil, v}

	if q.sttNode_ == nil {
		q.sttNode_ = n
	}

	if q.endNode_ == nil {
		q.endNode_ = n
	} else {
		q.endNode_.next_ = n
		q.endNode_ = n
	}
	q.size_++
}

func (q *Queue) PopFront() T {
	if q.sttNode_ != nil {
		s := q.sttNode_.value_

		q.sttNode_ = q.sttNode_.next_
		if q.sttNode_ != nil {
			q.sttNode_.prev_ = nil
		}

		// What about the last elem pointer?
		if q.sttNode_ == nil {
			q.endNode_ = nil
		}
		if s != nil {
			q.size_--
		}
		return s
	}
	return nil
}
func (q *Queue) Size() int {
	return q.size_
}
func (q *Queue) Print() {
	fmt.Printf("<")

	if q.sttNode_ != nil {
		q.subPrint(q.sttNode_)
	}

	fmt.Printf(">")
}

func (q *Queue) subPrint(x *node) {
	fmt.Printf("%v ", x.value_)

	if x.next_ != nil {
		q.subPrint(x.next_)
	}

}
