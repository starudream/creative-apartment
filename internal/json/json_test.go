package json

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMustUnmarshalTo(t *testing.T) {
	type X struct {
		A string
		B int
		C float64
		D bool
		E time.Time
		F []byte
	}

	var x = &X{
		A: "foo",
		B: 5,
		C: 3.14,
		D: true,
		E: time.Now().Truncate(time.Second),
		F: []byte("bar"),
	}

	t.Logf("%#v", x)

	bs := MustMarshal(x)

	t.Log(string(bs))

	xx := MustUnmarshalTo[*X](bs)

	t.Logf("%#v", xx)

	require.Equal(t, x.A, xx.A)
	require.Equal(t, x.B, xx.B)
	require.Equal(t, x.C, xx.C)
	require.Equal(t, x.D, xx.D)
	require.Equal(t, x.E, xx.E)
	require.Equal(t, x.F, xx.F)
}
