package ibolt

import (
	"sync"
	"time"

	"go.etcd.io/bbolt"

	"github.com/starudream/creative-apartment/internal/ilog"
)

type Pool struct {
	connections   map[string]*Connection
	mu            sync.RWMutex
	removeTrigger chan struct{}
	quit          chan struct{}
}

func New() *Pool {
	p := &Pool{
		connections:   map[string]*Connection{},
		removeTrigger: make(chan struct{}, 1),
		quit:          make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-p.removeTrigger:
				select {
				case <-time.After(connectionExpire):
				case <-p.quit:
					return
				}
				p.mu.Lock()
				for _, c := range p.connections {
					c.mu.RLock()
					if !c.closeTime.IsZero() && c.closeTime.Before(time.Now()) {
						ilog.WrapError(c.remove(), "bolt")
					}
					c.mu.RUnlock()
				}
				p.mu.Unlock()
			case <-p.quit:
				return
			}
		}
	}()

	return p
}

func (p *Pool) Get(path string) (*Connection, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if c, ok := p.connections[path]; ok {
		c.mu.Lock()
		c.increment()
		c.mu.Unlock()
		return c, nil
	}

	db, err := bbolt.Open(path, fileMode, options)
	if err != nil {
		return nil, err
	}

	c := &Connection{
		DB:   db,
		path: path,
		pool: p,
	}
	c.mu.Lock()
	c.increment()
	p.connections[path] = c
	c.mu.Unlock()
	return c, nil
}

func (p *Pool) Has(path string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	_, ok := p.connections[path]
	return ok
}

func (p *Pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, c := range p.connections {
		ilog.WrapError(c.remove(), "bolt")
	}
	close(p.quit)
}

func (p *Pool) remove(path string) error {
	c, ok := p.connections[path]
	if !ok {
		return nil
	}
	delete(p.connections, path)
	return c.DB.Close()
}

type Connection struct {
	DB *bbolt.DB

	pool      *Pool
	path      string
	count     int64
	closeTime time.Time
	mu        sync.RWMutex
}

func (c *Connection) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.decrement()

	if c.count > 0 {
		return
	}

	c.closeTime = time.Now().Add(connectionExpire)
	select {
	case c.pool.removeTrigger <- struct{}{}:
	default:
	}
}

func (c *Connection) increment() {
	c.closeTime = time.Time{}
	c.count++
}

func (c *Connection) decrement() {
	c.count--
}

func (c *Connection) remove() error {
	return c.pool.remove(c.path)
}
