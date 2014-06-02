// GOMAXPROCS=10 go test
// GOMAXPROCS=10 go test -bench . -race

package spin

import (
	"runtime"
	"sync/atomic"
	"testing"
)

func HammerSpinlock(l *Lock, loops int, done chan bool) {
	for i := 0; i < loops; i++ {
		l.Lock()
		l.Unlock()
	}
	done <- true
}

func TestSpinlock(t *testing.T) {
	l := new(Lock)
	c := make(chan bool)
	for i := 0; i < 10; i++ {
		go HammerSpinlock(l, 1000, c)
	}
	for i := 0; i < 10; i++ {
		<-c
	}
}

func BenchmarkUncontendedSpinlock(b *testing.B) {
	l := new(Lock)
	HammerSpinlock(l, b.N, make(chan bool, 2))
}

func BenchmarkContendedSpinlock(b *testing.B) {
	b.StopTimer()
	l := new(Lock)
	c := make(chan bool)
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(2))
	b.StartTimer()

	go HammerSpinlock(l, b.N/2, c)
	go HammerSpinlock(l, b.N/2, c)
	<-c
	<-c
}

func TestSpinlockPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("unlock of unlocked spinlock did not panic")
		}
	}()

	var l Lock
	l.Lock()
	l.Unlock()
	l.Unlock()
}

func BenchmarkSpinlockUncontended(b *testing.B) {
	type PaddedSpinlock struct {
		Lock
		pad [128]uint8
	}
	const CallsPerSched = 1000
	procs := runtime.GOMAXPROCS(-1)
	N := int32(b.N / CallsPerSched)
	c := make(chan bool, procs)
	for p := 0; p < procs; p++ {
		go func() {
			var l PaddedSpinlock
			for atomic.AddInt32(&N, -1) >= 0 {
				runtime.Gosched()
				for g := 0; g < CallsPerSched; g++ {
					l.Lock.Lock()
					l.Unlock()
				}
			}
			c <- true
		}()
	}
	for p := 0; p < procs; p++ {
		<-c
	}
}

func benchmarkSpinlock(b *testing.B, slack, work bool) {
	const (
		CallsPerSched  = 1000
		LocalWork      = 100
		GoroutineSlack = 10
	)
	procs := runtime.GOMAXPROCS(-1)
	if slack {
		procs *= GoroutineSlack
	}
	N := int32(b.N / CallsPerSched)
	c := make(chan bool, procs)
	l := new(Lock)
	for p := 0; p < procs; p++ {
		go func() {
			foo := 0
			for atomic.AddInt32(&N, -1) >= 0 {
				runtime.Gosched()
				for g := 0; g < CallsPerSched; g++ {
					l.Lock()
					l.Unlock()
					if work {
						for i := 0; i < LocalWork; i++ {
							foo *= 2
							foo /= 2
						}
					}
				}
			}
			c <- foo == 42
		}()
	}
	for p := 0; p < procs; p++ {
		<-c
	}
}

func BenchmarkSpinlock(b *testing.B) {
	benchmarkSpinlock(b, false, false)
}

func BenchmarkSpinlockSlack(b *testing.B) {
	benchmarkSpinlock(b, true, false)
}

func BenchmarkSpinlockWork(b *testing.B) {
	benchmarkSpinlock(b, false, true)
}

func BenchmarkSpinlockWorkSlack(b *testing.B) {
	benchmarkSpinlock(b, true, true)
}
