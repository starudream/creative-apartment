package ibolt

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"go.etcd.io/bbolt"
)

var (
	options = &bbolt.Options{
		Timeout:        30 * time.Second,
		NoGrowSync:     false,
		NoFreelistSync: false,
		FreelistType:   bbolt.FreelistMapType,
	}

	fileMode = os.FileMode(0600)

	connectionExpire = time.Minute

	pool *Pool
	path string

	poolOnce sync.Once
)

func Init(xpath string) {
	poolOnce.Do(func() {
		pool = New()
		path = xpath
	})
}

func Close() {
	if pool != nil {
		pool.Close()
		log.Info().Msg("[bolt] db closed")
	}
}

type DB = bbolt.DB

func D(fn func(*DB) error) error {
	conn, err := pool.Get(path)
	if err != nil {
		return err
	}
	defer conn.Close()
	return fn(conn.DB)
}

type Tx = bbolt.Tx

func Update(fn func(*Tx) error) error {
	conn, err := pool.Get(path)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.DB.Update(fn)
}

type Bucket = bbolt.Bucket

func View(fn func(*Tx) error) error {
	conn, err := pool.Get(path)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.DB.View(fn)
}

func Batch(fn func(*Tx) error) error {
	conn, err := pool.Get(path)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.DB.Batch(fn)
}
