# spin

spin provides a simple spinlock for Go. Spinlocks should be used in
cases where a mutex would be too much overhead and the lock will not
be held for long.

The most common use is for protecting a single variable that is
changed briefly.

```go
type Counter struct {
	c int
	l spin.Lock
}

func (c *Counter) Inc() {
	c.l.Lock()
	c.c++
	c.l.Unlock()
}
```
