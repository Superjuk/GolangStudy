package main

import (
	//"io"
	"net/http"
	//"net/http/httptest"
	//"errors"
	"testing"
)

type TestCase struct {
	Users *SearchResponse
	err   error
}

type TestRequest struct {
	Limit   int
	Err     string
	IsError bool
}

func SearchServer(rw http.ResponseWriter, r *http.Request) {
	// limit := r.FormValue("limit")
	// switch limit {
	// case limit < 0:
	// 	w.WriteHeader(http.StatusOK)
	// 	io.WriteString(w, `{"status": 200, "balance": 100500}`)
	// case "100500":
	// 	w.WriteHeader(http.StatusOK)
	// 	io.WriteString(w, `{"status": 400, "err": "bad_balance"}`)
	// default:
	// 	w.WriteHeader(http.StatusInternalServerError)
	// }
}

func TestLimit(t *testing.T) {
	cases := []TestRequest{
		TestRequest{
			Limit:   -1,
			Err:     "limit must be > 0",
			IsError: true,
		},
		TestRequest{
			Limit:   26,
			Err:     "",
			IsError: false,
		},
	}

	for _, item := range cases {
		req := &SearchRequest{}
		req.Limit = item.Limit

		srv := &SearchClient{}
		_, err := srv.FindUsers(*req)

		if item.IsError && err == nil {
			t.Errorf("Expecting error, when limit = %d", item.Limit)
		}
		if req.Limit > -1 && req.Limit != 25 {
			t.Errorf("req.Limit must be equal 25, got %d", req.Limit)
		}
	}
}
