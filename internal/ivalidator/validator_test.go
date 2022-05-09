package ivalidator

import (
	"testing"
)

func TestStruct(t *testing.T) {
	t.Log(Struct(nil))
	t.Log(Struct(nil) == nil)
}

func TestVar(t *testing.T) {
	t.Log(Var("abc", "AccessToken", "min=1"))
	t.Log(Var("abc", "AccessToken", "min=5"))
}
