package itask

import (
	"context"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
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

	fCh := make(chan F)
	errCh := make(chan error, 1)
	errSig := int32(0)
	errOnce := sync.Once{}

	for i := 0; i < n; i++ {
		go func() {
			for {
				f, ok := <-fCh
				if !ok {
					return
				}
				if f == nil || atomic.LoadInt32(&errSig) != 0 {
					return
				}
				func() {
					defer func() {
						if ev := recover(); ev != nil {
							log.Error().CallerSkipFrame(2).Msgf("[task] panic, %v\n%s", ev, debug.Stack())
						}
						wg.Add(-1)
					}()
					err := f(t.ctx)
					if err != nil {
						errOnce.Do(func() {
							atomic.StoreInt32(&errSig, 1)
							errCh <- err
						})
					}
				}()
			}
		}()
	}
	defer t.Stop()

	for i := range t.fs {
		fCh <- t.fs[i]
	}

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
