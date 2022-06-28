package waitfree

import (
	"sync/atomic"
	"unsafe"
)

/*
Simple, Fast, and Practical Non-Blocking and Blocking Concurrent Queue Algorithms
by Maged M. Michael and Michael L. Scott
https://cs.rochester.edu/u/scott/papers/1996_PODC_queues.pdf
*/
type LinkedListQueue[T any] struct {
	front unsafe.Pointer
	back unsafe.Pointer
}

type node[T any] struct {
	val T
	next unsafe.Pointer
}

/*
structure pointer t fptr: pointer to node t, count: unsigned integerg
structure node t fvalue: data type, next: pointer tg
structure queue t fHead: pointer t, Tail: pointer tg

initialize(Q: pointer to queue t)
node = new node() 					# Allocate a free node
node–>next.ptr = NULL 				# Make it the only node in the linked list
Q–>Head = Q–>Tail = node 			# Both Head and Tail point to it
*/
func NewLinkedListQueue[T any]() *LinkedListQueue[T] {
	p := unsafe.Pointer(new(node[T]))
	return &LinkedListQueue[T]{p, p}
}

/*
enqueue(Q: pointer to queue t, value: data type)
E1: node = new node() 													# Allocate a new node from the free list
E2: node–>value = value 												# Copy enqueued value into node
E3: node–>next.ptr = NULL 												# Set next pointer of node to NULL
E4: loop 																# Keep trying until Enqueue is done
E5: 	tail = Q–>Tail 													# Read Tail.ptr and Tail.count together
E6: 	next = tail.ptr–>next 											# Read next ptr and count fields together
E7: 	if tail == Q–>Tail 												# Are tail and next consistent?
E8: 		if next.ptr == NULL 										# Was Tail pointing to the last node?
E9: 			if CAS(&tail.ptr–>next, next, <node, next.count+1>) 	# Try to link node at the end of the linked list
E10: 				break 												# Enqueue is done. Exit loop
E11: 			endif
E12: 		else 														# Tail was not pointing to the last node
E13: 			CAS(&Q–>Tail, tail, <next.ptr, tail.count+1>) 			# Try to swing Tail to the next node
E14: 		endif
E15: 	endif
E16: endloop
E17: CAS(&Q–>Tail, tail, <node, tail.count+1>) 							# Enqueue is done. Try to swing Tail to the inserted node
*/
func (q *LinkedListQueue[T]) Enqueue(element T) {
	new_node := unsafe.Pointer(&node[T]{element, nil})
	var old_back unsafe.Pointer
	for {
		old_back = atomic.LoadPointer(&q.back)
		next := atomic.LoadPointer(&(*node[T])(old_back).next)
		if old_back == atomic.LoadPointer(&q.back) {
			if next == nil {
				if atomic.CompareAndSwapPointer(&(*node[T])(old_back).next, nil, new_node) {
					break
				}
			} else {
				new_next := atomic.LoadPointer(&(*node[T])(old_back).next)
				atomic.CompareAndSwapPointer(&q.back, old_back, new_next)
			}
		}
	}
	atomic.CompareAndSwapPointer(&q.back, old_back, new_node)
}

/*
dequeue(Q: pointer to queue t, pvalue: pointer to data type): boolean
D1: loop 																# Keep trying until Dequeue is done
D2: 	head = Q–>Head 													# Read Head
D3: 	tail = Q–>Tail 													# Read Tail
D4: 	next = head–>next 												# Read Head.ptr–>next
D5: 	if head == Q–>Head 												# Are head, tail, and next consistent?
D6: 		if head.ptr == tail.ptr 									# Is queue empty or Tail falling behind?
D7: 			if next.ptr == NULL 									# Is queue empty?
D8: 				return FALSE 										# Queue is empty, couldn’t dequeue
D9: 			endif
D10: 			CAS(&Q–>Tail, tail, <next.ptr, tail.count+1>) 			# Tail is falling behind. Try to advance it
D11: 		else 														# No need to deal with Tail
																		# Read value before CAS, otherwise another dequeue might free the next node
D12: 			*pvalue = next.ptr–>value
D13: 			if CAS(&Q–>Head, head, <next.ptr, head.count+1>) 		# Try to swing Head to the next node
D14: 				break 												# Dequeue is done. Exit loop
D15: 			endif
D16: 		endif
D17: 	endif
D18: endloop
D19: free(head.ptr) 													# It is safe now to free the old dummy node
D20: return TRUE 														# Queue was not empty, dequeue succeeded
*/
func (q *LinkedListQueue[T]) Dequeue() (T, bool) {
	var element T
	for {
		front := atomic.LoadPointer(&q.front)
		back := atomic.LoadPointer(&q.back)
		next := atomic.LoadPointer(&(*node[T])(front).next)
		if front == atomic.LoadPointer(&q.front) {
			if front == back {
				if next == nil {
					var empty T
					return empty, false
				} 
				atomic.CompareAndSwapPointer(&q.back, back, next)
			} else {
				element = (*node[T])(next).val
				if atomic.CompareAndSwapPointer(&q.front, front, next) {
					break
				}
			}
		}
	}
	return element, true
}