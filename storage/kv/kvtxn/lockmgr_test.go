package kvtxn

import (
	"sync"
	"testing"
	"time"
)

func TestKeyLockManager(t *testing.T) {
	timedKeyLockManagerTest(t, NewInmemLockManager())
}

func timedKeyLockManagerTest(t *testing.T, klm KeyLockManager) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		klm.Lock("lock_key")

		// signify we've completed locking
		wg.Done()

		time.Sleep(100 * time.Millisecond)

		klm.Unlock("lock_key")
	}()

	// wait until we've completed locking
	wg.Wait()

	ch := make(chan struct{})
	go func() {
		klm.RLock("lock_key")

		// signify that we've completed the lock
		close(ch)
	}()

	select {
	case <-ch:
		t.Error("expected to be blocked, but returned early")
	case <-time.After(50 * time.Millisecond):
		// this timer should fire before the the first timer, above.
		// which means the lock for the key should still be locked
		// i.e. this case is expected behavior to exit the select
	}

	// wait until we're done with the second lock
	<-ch

	klm.RUnlock("lock_key")

	select {
	case <-ch:
		//
		// Expected behavior: we got through after unlock
	default:
		t.Error("expected second lock to have happened, but didn't")
	}
}
