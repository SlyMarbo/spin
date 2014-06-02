package spin

import (
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
	for {
		if atomic.CompareAndSwapInt32(&l.state, 0, 1) {
			return
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
