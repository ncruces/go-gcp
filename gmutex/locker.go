package gmutex

import "context"

type locker struct {
	*Mutex
}

func (m locker) Lock() {
	if err := m.LockContext(context.Background()); err != nil {
		panic(err)
	}
}

func (m locker) Unlock() {
	if err := m.UnlockContext(context.Background()); err != nil {
		panic(err)
	}
}
