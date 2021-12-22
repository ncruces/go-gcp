package gmutex

import "context"

type locker struct {
	*Mutex
}

func (m locker) Lock() {
	if err := m.Mutex.Lock(context.Background()); err != nil {
		panic(err)
	}
}

func (m locker) Unlock() {
	if err := m.Mutex.Unlock(context.Background()); err != nil {
		panic(err)
	}
}
