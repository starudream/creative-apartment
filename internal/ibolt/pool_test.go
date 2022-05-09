package ibolt

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/starudream/creative-apartment/internal/json"
)

func Test(t *testing.T) {
	Init("test.bolt")
	defer Close()

	e1 := X(func(db *DB) error {
		t.Log(db.String())
		t.Log(json.MustMarshalString(db.Stats()))
		t.Log(json.MustMarshalString(db.Info()))
		return nil
	})
	require.NoError(t, e1)

	e2 := Update(func(tx *Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("daily"))
		if err != nil {
			return err
		}
		return bucket.Put([]byte("foo"), []byte("bar"))
	})
	require.NoError(t, e2)

	e3 := View(func(tx *Tx) error {
		return tx.ForEach(func(name []byte, b *Bucket) error {
			t.Log(string(name))
			t.Log(json.MustMarshalString(b.Stats()))
			return nil
		})
	})
	require.NoError(t, e3)
}
