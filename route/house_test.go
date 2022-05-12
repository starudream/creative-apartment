package route

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/starudream/creative-apartment/internal/igin"
	"github.com/starudream/creative-apartment/internal/json"
)

func TestGetHouseData(t *testing.T) {
	req := &GetHouseDataReq{
		Phone:     "13312341234",
		StartDate: "2022-05-01",
		EndDate:   "2022-05-20",
	}

	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(json.MustMarshal(req)))
	w := httptest.NewRecorder()

	e := igin.New()
	e.POST("/", GetHouseData)
	e.ServeHTTP(w, r)

	t.Log(w.Body.String())
}
