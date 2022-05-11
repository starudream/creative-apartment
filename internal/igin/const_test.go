package igin

import (
	"net/http"
	"testing"
)

func TestConsts(t *testing.T) {
	headers := []string{"content-type", "x-request-id"}
	for i := 0; i < len(headers); i++ {
		t.Log(http.CanonicalHeaderKey(headers[i]))
	}
}
