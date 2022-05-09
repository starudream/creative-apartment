package api

import (
	"testing"

	"github.com/starudream/creative-apartment/internal/json"
)

func TestGetLatestAPK(t *testing.T) {
	t.Log(json.MustMarshalString(GetLatestAPK()))
	t.Log(json.MustMarshalString(GetLatestAPK(true)))
}
