package batch

import (
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {
	g := &errgroup.Group{}
	mu := &sync.Mutex{}

	// Limit working goroutines to the pool constraint
	g.SetLimit(int(pool))

	for i := int64(0); i < n; i++ {
		id := i
		g.Go(func() error {
			user := getOne(id)
			mu.Lock() // Lock while we write to []user
			res = append(res, user)
			mu.Unlock()
			return nil
		})
	}

	// Wait for group to finish, if first non-nil error appears finish our task with no res,
	// but as far as getOne() returns no errors it will be ok
	if err := g.Wait(); err != nil {
		return nil
	}

	return res
}
