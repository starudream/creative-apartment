package itask

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/rs/zerolog/log"
)

type Task struct {
	pool *ants.Pool

	mu sync.Mutex
	fs []func() error
}

func New(sizes ...int) *Task {
	if len(sizes) == 0 || sizes[0] <= 0 {
		sizes = []int{2 * runtime.NumCPU()}
	}
	pool, _ := ants.NewPool(sizes[0], ants.WithLogger(&logger{log.Logger}))
	return &Task{pool: pool}
}

func (t *Task) Add(fs ...func() error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.fs = append(t.fs, fs...)
}

func (t *Task) Run(ctx context.Context, timeout time.Duration) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	errCh := make(chan error, 1)
	for i := range t.fs {
		f := t.fs[i]
		err := t.pool.Submit(func() { errCh <- f() })
		if err != nil {
			return err
		}
	}
	defer t.pool.Release()

	nc, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	tk := time.NewTicker(time.Second)
	for {
		select {
		case err := <-errCh:
			if err != nil {
				return err
			}
		case <-nc.Done():
			return nc.Err()
		case <-tk.C:
			if t.pool.Running() == 0 && t.pool.Waiting() == 0 {
				return nil
			}
		}
	}
}
