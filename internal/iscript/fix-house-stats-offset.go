package iscript

import (
	"bytes"
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.etcd.io/bbolt"

	"github.com/starudream/creative-apartment/internal/ibolt"
	"github.com/starudream/creative-apartment/internal/ilog"
)

func FixHouseStatsOffset(context.Context) error {
	return ibolt.Update(func(tx *ibolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("config"))
		if !ilog.WrapError(err) {
			return err
		}

		if len(bucket.Get([]byte("fix-house-stats-offset"))) > 0 {
			return nil
		}

		err = tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
			if !bytes.Contains(name, []byte("_house_stats_")) {
				return nil
			}

			log.Debug().Str("bucket", string(name)).Msg("fixing house stats offset")

			m := map[string]string{}

			err = b.ForEach(func(k, v []byte) error { m[string(k)] = string(v); return nil })
			if !ilog.WrapError(err) {
				return err
			}

			for k, v := range m {
				t, te := time.ParseInLocation("0102", k[:4], time.Local)
				if !ilog.WrapError(te) {
					return te
				}

				k = t.AddDate(0, 0, -1).Format("0102") + k[4:]

				pe := b.Put([]byte(k), []byte(v))
				if !ilog.WrapError(pe) {
					return pe
				}
			}

			return nil
		})
		if !ilog.WrapError(err) {
			return err
		}

		return bucket.Put([]byte("fix-house-stats-offset"), []byte("1"))
	})
}
