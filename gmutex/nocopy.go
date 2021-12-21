package gmutex

import "sync"

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

var _ sync.Locker = (*noCopy)(nil)
