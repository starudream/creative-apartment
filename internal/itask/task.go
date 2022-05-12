package itask

import (
	"context"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Task struct {
	ctx    context.Context
	cancel func()

	mu sync.Mutex
	fs []F
}

type F func(ctx context.Context) error

func New(ctx context.Context, timeout time.Duration) *Task {
	if ctx == nil {
		ctx = context.Background()
	}

	task := &Task{}

	if timeout > 0 {
		task.ctx, task.cancel = context.WithTimeout(ctx, timeout)
	} else {
		task.ctx, task.cancel = context.WithCancel(ctx)
	}

	return task
}

func (t *Task) Add(fs ...F) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.fs = append(t.fs, fs...)
}

func (t *Task) Run(n int) error {
	if n <= 0 {
		n = 2 * runtime.NumCPU()
	}

	wg := &sync.WaitGroup{}

	wg.Add(len(t.fs))

	errCh := make(chan error, 1)
	for i := 0; i < len(t.fs); i++ {
		f := t.fs[i]
		go func() {
			defer func() {
				if ev := recover(); ev != nil {
					log.Error().CallerSkipFrame(2).Msgf("[task] panic, %v\n%s", ev, debug.Stack())
				}
				wg.Add(-1)
			}()
			errCh <- f(t.ctx)
		}()
	}
	defer t.Stop()

	doneCh := make(chan struct{}, 1)

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case err := <-errCh:
			if err != nil {
				t.Stop()
				return err
			}
		case <-doneCh:
			t.Stop()
			return nil
		case <-t.ctx.Done():
			return t.ctx.Err()
		}
	}
}

func (t *Task) Stop() {
	t.cancel()
}
