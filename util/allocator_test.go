package util

import "testing"

func TestAllocator(t *testing.T) {

	alloc := NewAllocator([...]int{5, 10, 20})
	if cap(alloc.Alloc(1)) != 5 {
		t.Log("level 1 alloc failed")
		t.FailNow()
	}

	if cap(alloc.Alloc(6)) != 10 {
		t.Log("level 2 alloc failed")
		t.FailNow()
	}

	if cap(alloc.Alloc(12)) != 20 {
		t.Log("level 3 alloc failed")
		t.FailNow()
	}

	if cap(alloc.Alloc(100)) != 100 {
		t.Log("level 3+ alloc failed")
		t.FailNow()
	}
}

func BenchmarkAllocator(b *testing.B) {

	alloc := NewAllocator([...]int{5, 10, 20})
	for i := 0; i < b.N; i++ {

		data := alloc.Alloc(11)

		alloc.Free(data)
	}

}

func BenchmarkAllocatorClassic(b *testing.B) {

	for i := 0; i < b.N; i++ {

		data := make([]byte, 20)

		data = data
	}

}
