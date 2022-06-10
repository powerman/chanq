// Package chanq provides outgoing queue for channel to use in select case.
package chanq

// Queue is an outgoing queue for channel to use in select case.
//
//   out := make(chan []byte) // Usually not buffered, with blocking send.
//   q := NewQueue(out)
//   q.Enqueue([]byte(`one`))
//   q.Enqueue([]byte(`two`))
//   for {
//     select {
//     case data := <-in: // E.g.: forward from in to out without blocking.
//       q.Enqueue(data)
//     case q.C <- q.Elem: // Works only when queue is not empty.
//       q.Dequeue()
//     }
//   }
type Queue[T any] struct {
	C     chan<- T // Use only this way: q.C <- q.Elem.
	Elem  T        // Use only this way: q.C <- q.Elem.
	Queue []T      // You can make any changes, but not concurrently with any other use.
	c     chan<- T // Unbounded send, won't be closed while Queue is in use.
}

// NewQueue creates new unrestricted queue for given out channel.
// This way you can handle channel like it has an unlimited buffer.
//
// You must not close the channel while you use the queue.
func NewQueue[T any](out chan<- T) *Queue[T] {
	return &Queue[T]{
		c: out,
	}
}

// Enqueue adds new elem to queue. If queue was empty this will result in
// unblocking next attempt to send q.Elem into q.C.
func (q *Queue[T]) Enqueue(elem T) {
	if len(q.Queue) == 0 {
		q.C = q.c
		q.Elem = elem
	}
	q.Queue = append(q.Queue, elem)
}

// Dequeue removes just sent element from queue. It must be called after
// each successful sent of q.Elem into q.C. If queue will became empty
// after this call it will block next attempt to send into q.C.
func (q *Queue[T]) Dequeue() {
	q.Queue = q.Queue[1:]

	if len(q.Queue) == 0 {
		q.C = nil
	} else {
		q.Elem = q.Queue[0]
	}
}
