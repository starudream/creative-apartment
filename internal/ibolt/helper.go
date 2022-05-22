package ibolt

import (
	"github.com/spf13/cast"
)

func CreateBuckets(bucketNames ...string) error {
	return Update(func(tx *Tx) error {
		for _, name := range bucketNames {
			_, err := tx.CreateBucketIfNotExists([]byte(name))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func Get(bucketName, key string) (v []byte, e error) {
	e = View(func(tx *Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		v = bucket.Get([]byte(key))
		return nil
	})
	return
}

func Put(bucketName, key string, value any) error {
	v, ve := cast.ToStringE(value)
	if ve != nil {
		return ve
	}
	return Update(func(tx *Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(key), []byte(v))
	})
}
