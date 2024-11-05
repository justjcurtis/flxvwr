package utils

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestDebounce(t *testing.T) {
	var count int32

	fn := func() {
		atomic.AddInt32(&count, 1)
	}

	debounceDelay := 100 * time.Millisecond
	debouncedFn := Debounce(fn, debounceDelay)

	for i := 0; i < 10; i++ {
		debouncedFn()
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(2 * debounceDelay)

	if atomic.LoadInt32(&count) != 1 {
		t.Fatalf("Expected fn to be called once, but got %d", count)
	}
}
func TestDebounce_WithNoCallsAfter(t *testing.T) {
	var count int32
	fn := func() {
		atomic.AddInt32(&count, 1)
	}

	debounceDelay := 100 * time.Millisecond
	debouncedFn := Debounce(fn, debounceDelay)

	for i := 0; i < 5; i++ {
		debouncedFn()
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(2 * debounceDelay)

	if atomic.LoadInt32(&count) != 1 {
		t.Fatalf("Expected fn to be called once, but got %d", count)
	}
}
func TestDebounce_CalledOnceIfNotInterrupted(t *testing.T) {
	var count int32
	fn := func() {
		atomic.AddInt32(&count, 1)
	}

	debounceDelay := 200 * time.Millisecond
	debouncedFn := Debounce(fn, debounceDelay)

	debouncedFn()

	time.Sleep(2 * debounceDelay)

	if atomic.LoadInt32(&count) != 1 {
		t.Fatalf("Expected fn to be called once, but got %d", count)
	}
}
