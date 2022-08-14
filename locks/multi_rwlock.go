package locks

import (
	"sync"
)

//---------Multi RWlock

type MultiRWLock interface {
	Lock(int)
	Unlock(int)
	RLock(int)
	RUnlock(int)
}

type mappedMultiRWLock struct {
	lockmap sync.Map
}

func (mml *mappedMultiRWLock) getMutex(id int) *sync.RWMutex {
	lock, ok := mml.lockmap.Load(id)
	if !ok {
		lock = new(sync.RWMutex)
		mml.lockmap.Store(id, lock)
	}
	return lock.(*sync.RWMutex)
}

func (mml *mappedMultiRWLock) Lock(id int) {
	mml.getMutex(id).Lock()
}

func (mml *mappedMultiRWLock) Unlock(id int) {
	mml.getMutex(id).Unlock()
}

func (mml *mappedMultiRWLock) RLock(id int) {
	mml.getMutex(id).RLock()
}

func (mml *mappedMultiRWLock) RUnlock(id int) {
	mml.getMutex(id).RUnlock()
}

func NewMappedMultiRWMutex() MultiRWLock {
	return new(mappedMultiRWLock)
}