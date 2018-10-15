package lfqueue

import (
	"runtime"
	"sync/atomic"
)

type Queue struct {
	capacityMod    uint32
	capacity       uint32
	getPos uint32
    _1 [64]byte
    putPos uint32
    // _2 [64]byte
	entries        []queueEntry
}

type queueEntry struct {
	putPos uint32
	getPos uint32
	elem   interface{}
}

func NewQueue(capacity int) *Queue {
	q := &Queue{
		capacity: minQuantity(uint32(capacity + 1)),
		putPos:   0,
		getPos:   0,
	}

	q.capacityMod = q.capacity - 1
	q.entries = make([]queueEntry, q.capacity)
	for i := range q.entries {
		entry := &q.entries[i]
		entry.getPos = uint32(i)
		entry.putPos = uint32(i)
	}

	return q
}

func (q *Queue) Put(elem interface{}) (ok bool) {
	capacityMod := q.capacityMod

retry:
	getPos := atomic.LoadUint32(&q.getPos)
	putPos := atomic.LoadUint32(&q.putPos)

	var cnt uint32
	if putPos >= getPos {
		cnt = putPos - getPos
	} else { // putPos < getPos when putPos exceed uint32 boundary
		cnt = (putPos - getPos) + capacityMod
	}

	if cnt >= capacityMod {
		return false
	}

	if !atomic.CompareAndSwapUint32(&q.putPos, putPos, putPos+1) {
		runtime.Gosched()
		goto retry
	}

	entry := &q.entries[putPos&capacityMod]
	for {
		entryGetPos := atomic.LoadUint32(&entry.getPos)
		entryPutPos := atomic.LoadUint32(&entry.putPos)
		if putPos == entryPutPos && entryPutPos == entryGetPos {
			entry.elem = elem
			atomic.AddUint32(&entry.putPos, q.capacity)
			return true
		} else {
			runtime.Gosched()
		}
	}
}

func (q *Queue) Get() (elem interface{}, ok bool) {
	capacity := q.capacity
	capacityMod := q.capacityMod
retry:
	getPos := atomic.LoadUint32(&q.getPos)
	putPos := atomic.LoadUint32(&q.putPos)

	var cnt uint32
	if putPos >= getPos {
		cnt = putPos - getPos
	} else { // putPos < getPos when putPos exceed uint32 boundary
		cnt = (putPos - getPos) + capacityMod
	}

	if cnt <= 0 {
		return nil, false
	}

	if !atomic.CompareAndSwapUint32(&q.getPos, getPos, getPos+1) {
		runtime.Gosched()
		goto retry
	}

	entry := &q.entries[getPos&capacityMod]
	for {
		entryGetPos := atomic.LoadUint32(&entry.getPos)
		entryPutPos := atomic.LoadUint32(&entry.putPos)
		if getPos == entryGetPos && entryGetPos == entryPutPos-capacity {
			elem := entry.elem
			entry.elem = nil
			atomic.AddUint32(&entry.getPos, capacity)
			return elem, true
		} else {
			runtime.Gosched()
		}
	}
}

func (q *Queue) Size() int {
	getPos := atomic.LoadUint32(&q.getPos)
	putPos := atomic.LoadUint32(&q.putPos)

	var cnt uint32
	if putPos >= getPos {
		cnt = putPos - getPos
	} else { // putPos < getPos when putPos exceed uint32 boundary
		cnt = (putPos - getPos) + q.capacityMod
	}
	return int(cnt)
}

func minQuantity(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}
