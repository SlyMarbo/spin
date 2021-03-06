// Package spin provides a simple spinlock. Spinlocks should be used in
// cases where a mutex would be too much overhead and the lock will not
// be held for long.
//
// The most common use is for protecting a single variable that is
// changed briefly.
//
//		type Counter struct {
//			c int
//			l spin.Lock
//		}
//
//		func (c *Counter) Inc() {
//			c.l.Lock()
//			c.c++
//			c.l.Unlock()
//		}
//
package spin

import (
	"runtime"
	"sync/atomic"
)

// Lock is a simple spinlock.
// The key difference from sync/Mutex is
// that spinlock spins in a loop while
// waiting to lock, so it should only be
// used when the lock will be held briefly.
//
// The zero value of Lock is an unlocked
// Lock.
type Lock struct {
	state int32
}

// Lock locks l.
// If the Lock is already in use, the calling
// goroutine spins until the Lock is available.
func (l *Lock) Lock() {
	i := 0
	for !atomic.CompareAndSwapInt32(&l.state, 0, 1) {
		if i++; i == 1024 {
			runtime.Gosched()
			i = 0
		}
	}
}

// Lock unlocks l.
// It is a runtime error if l is not locked on
// entry to Unlock.
//
// A locked Lock is not associated with a
// particular goroutine. It is allowed for one
// goroutine to lock a Lock and then arrange for
// another goroutine to unlock it.
func (l *Lock) Unlock() {
	if !atomic.CompareAndSwapInt32(&l.state, 1, 0) {
		panic("spin: unlock of unlocked Lock")
	}
}
