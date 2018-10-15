package lfqueue

type Queue struct {
	capacity       uint32
	capacityMod    uint32
	getPos, putPos uint32
	entries        []queueEntry
}

type queueEntry struct {
	putPos uint32
	getPos uint32
	elem   interface{}
}

func NewQueue(capacity uint32) *Queue {
	q := &Queue{
		capacity: minQuantity(capacity),
		putPos:   0,
		getPos:   0,
	}

	q.capacityMod = q.capacity - 1
	q.entries = make([]queueEntry, q.capacity)
	for i, entry := range q.entries {
		entry.getPos = uint32(i)
		entry.putPos = uint32(i)
	}

	return q
}

func (q *Queue) Put(elem interface{}) (ok bool) {
	return false
}

func (q *Queue) Get() (elem interface{}, ok bool) {
	return nil, false
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
