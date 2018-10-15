package lfqueue

import "testing"

func TestQueueBasic(t *testing.T) {
	capacity := 100
	q := NewQueue(capacity)
	for i := 0; i < capacity; i++ {
		ok := q.Put(i)
		if !ok {
			t.Fatalf("put %d returns false", i)
		}
	}
	for i := 0; i < capacity; i++ {
		elem, ok := q.Get()
		if !ok {
			t.Fatalf("get %d returns false", i)
		}
		if elem.(int) != i {
			t.Fatalf("Get wrong value")
		}
	}
}
