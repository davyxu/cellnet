package util

import "sync"

type Allocator struct {
	pool             sync.Pool
	maxPoolAllocSize int
}

func (self *Allocator) Alloc(size int) []byte {

	if size > self.maxPoolAllocSize {
		return make([]byte, size)
	}

	return self.pool.Get().([]byte)[:size]
}

func (self *Allocator) Free(data []byte) {

	if cap(data) > self.maxPoolAllocSize {

		// GC
		return
	}

	self.pool.Put(data)
}

func NewAllocator(maxPoolAllocSize int) *Allocator {

	self := &Allocator{
		maxPoolAllocSize: maxPoolAllocSize,
	}

	self.pool.New = func() interface{} {
		return make([]byte, maxPoolAllocSize)
	}

	return self
}
