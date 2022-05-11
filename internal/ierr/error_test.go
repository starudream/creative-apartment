package ierr

import (
	"testing"

	"github.com/starudream/creative-apartment/internal/json"
)

func TestError(t *testing.T) {
	e := New()
	t.Logf("%#v", e)
	t.Logf("%#v", e.WithCode(200))
	t.Logf("%#v", e.WithMessage("ok"))
	t.Logf("%#v", e.WithMetadata(MD{"k1": "v1"}))
	t.Logf("%#v", e.WithMetadata(MD{"k2": "v2"}))

	t.Logf("%#v", New(200))
	t.Logf("%#v", New(200, "ok"))
	t.Logf("%#v", New(999, "bool: %t, string: %s", true, "ok"))

	t.Logf("%s", New(999, "bool: %t", true).WithMetadata(MD{"foo": "bar"}))

	t.Logf("%s", json.MustMarshal(New(999, "bool: %t", true)))
	t.Logf("%s", json.MustMarshal(New(999, "bool: %t", true).WithMetadata(MD{"foo": "bar"})))
}
