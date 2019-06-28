package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type TestRequestLimit struct {
	Limit   int
	Err     string
	IsError bool
}

type TestRequestOffset struct {
	Limit   int
	Offset  int
	Err     string
	IsError bool
}

type TestRequest struct {
	Request SearchRequest
	Err     string
	IsError bool
}

const (
	securityToken = "SearchServer"
	wrongToken    = "error"

	wrongField        = "wrong"
	wrongUnknownField = "unknown"
)

var (
	searchFields = []string{"Id", "Age", "Name"}
)

func SearchServer(rw http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("AccessToken")
	if token != securityToken {
		rw.WriteHeader(http.StatusUnauthorized)
	}

	field := r.FormValue("order_field")
	if field != "" {
		isContain := false

		for _, item := range searchFields {
			if item == field {
				isContain = true
			}
		}

		if !isContain {
			rw.WriteHeader(http.StatusBadRequest)
			if field == wrongField {
				io.WriteString(rw, `{"error": "ErrorBadOrderField"}`)
			} else {
				io.WriteString(rw, `{"error": "AnotherErrorBadOrderField"}`)
			}
			return
		}
	} else {
		field = "Name"
	}

	query := r.FormValue("query")

	rw.WriteHeader(http.StatusOK)
	switch query {
	case "Rose":
		io.WriteString(rw, `[{"id": 9, "name": "Rose Carney", "age": 36, "about": "many words", "gender": "female"}, {"id": 13, "name": "Rose Carney", "age": 36, "about": "so on", "gender": "female"}]`)
	case "Wolf":
		io.WriteString(rw, `[{"id": 2, "name": "Boyd Wolf", "age": 22, "about": "few words", "gender": "male"}, {"id": 3, "name": "Boyd Wolf", "age": 22, "about": "silence", "gender": "male"}]`)
	}
}

func TimeoutServer(rw http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
}

func FatalErrorServer(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusInternalServerError)
}

func WrongErrorJsonServer(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
	io.WriteString(rw, `{"error": ErrorBadOrderField}`)
}

func WrongResultJsonServer(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	io.WriteString(rw, `{"error": ErrorBadOrderField}`)
}

// Tests
func TestLimit(t *testing.T) {
	cases := []TestRequestLimit{
		TestRequestLimit{
			Limit:   -1,
			Err:     "limit must be > 0",
			IsError: true,
		},
	}

	for _, item := range cases {
		req := SearchRequest{}
		req.Limit = item.Limit

		srv := &SearchClient{}
		_, err := srv.FindUsers(req)

		if err == nil {
			t.Errorf("Expecting [%v], got nil", item.Err)
		}

		if item.IsError && !strings.Contains(err.Error(), item.Err) {
			t.Errorf("Expecting [%v], got [%v]", item.Err, err)
		}
	}
}

func TestOffset(t *testing.T) {
	cases := []TestRequestOffset{
		TestRequestOffset{
			Limit:   1,
			Offset:  -1,
			Err:     "offset must be > 0",
			IsError: true,
		},
	}

	for _, item := range cases {
		req := SearchRequest{}
		req.Limit = item.Limit
		req.Offset = item.Offset

		srv := &SearchClient{}
		_, err := srv.FindUsers(req)

		if err == nil {
			t.Errorf("Expecting [%v], got nil", item.Err)
		}

		if item.IsError && !strings.Contains(err.Error(), item.Err) {
			t.Errorf("Expecting [%v], got [%v]", item.Err, err)
		}
	}
}

func TestUnknownError(t *testing.T) {
	cases := []TestRequest{
		TestRequest{
			Request: SearchRequest{
				Limit: 26,
			},
			Err:     "unknown error",
			IsError: true,
		},
	}

	for _, item := range cases {
		srv := &SearchClient{}
		_, err := srv.FindUsers(item.Request)

		if err == nil {
			t.Errorf("Expecting [%v], got nil", item.Err)
		}

		if item.IsError && !strings.Contains(err.Error(), item.Err) {
			t.Errorf("Expecting [%v], got [%v]", item.Err, err)
		}
	}
}

func TestTimeout(t *testing.T) {
	cases := []TestRequest{
		TestRequest{
			Request: SearchRequest{
				Limit:      10,
				Offset:     0,
				Query:      "Wolf",
				OrderField: "Id",
				OrderBy:    OrderByAsIs,
			},
			Err:     "timeout for",
			IsError: true,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(TimeoutServer))

	for _, item := range cases {
		srv := &SearchClient{
			AccessToken: wrongToken,
			URL:         ts.URL,
		}
		_, err := srv.FindUsers(item.Request)

		if err == nil {
			t.Errorf("Expecting [%v], got nil", item.Err)
		}

		if item.IsError && !strings.Contains(err.Error(), item.Err) {
			t.Errorf("Expecting [%v], got [%v]", item.Err, err)
		}
	}
}

func TestBadAccess(t *testing.T) {
	cases := []TestRequest{
		TestRequest{
			Request: SearchRequest{
				Limit:      10,
				Offset:     0,
				Query:      "Wolf",
				OrderField: "Id",
				OrderBy:    OrderByAsIs,
			},
			Err:     "Bad AccessToken",
			IsError: true,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	for _, item := range cases {
		srv := &SearchClient{
			AccessToken: wrongToken,
			URL:         ts.URL,
		}
		_, err := srv.FindUsers(item.Request)

		if err == nil {
			t.Errorf("Expecting [%v], got nil", item.Err)
		}

		if item.IsError && !strings.Contains(err.Error(), item.Err) {
			t.Errorf("Expecting [%v], got [%v]", item.Err, err)
		}
	}
}

func TestFatalError(t *testing.T) {
	cases := []TestRequest{
		TestRequest{
			Request: SearchRequest{
				Limit:      10,
				Offset:     0,
				Query:      "Wolf",
				OrderField: "Id",
				OrderBy:    OrderByAsIs,
			},
			Err:     "SearchServer fatal error",
			IsError: true,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(FatalErrorServer))

	for _, item := range cases {
		srv := &SearchClient{
			AccessToken: wrongToken,
			URL:         ts.URL,
		}
		_, err := srv.FindUsers(item.Request)

		if err == nil {
			t.Errorf("Expecting [%v], got nil", item.Err)
		}

		if item.IsError && !strings.Contains(err.Error(), item.Err) {
			t.Errorf("Expecting [%v], got [%v]", item.Err, err)
		}
	}
}

func TestBadRequest(t *testing.T) {
	cases := []TestRequest{
		TestRequest{
			Request: SearchRequest{
				Limit:      10,
				Offset:     0,
				Query:      "Wolf",
				OrderField: wrongField,
				OrderBy:    OrderByAsIs,
			},
			Err:     "OrderFeld " + wrongField + " invalid",
			IsError: true,
		},
		TestRequest{
			Request: SearchRequest{
				Limit:      10,
				Offset:     0,
				Query:      "Wolf",
				OrderField: wrongUnknownField,
				OrderBy:    OrderByAsIs,
			},
			Err:     "unknown bad request error",
			IsError: true,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	for _, item := range cases {
		srv := &SearchClient{
			AccessToken: securityToken,
			URL:         ts.URL,
		}
		_, err := srv.FindUsers(item.Request)

		if err == nil {
			t.Errorf("Expecting [%v], got nil", item.Err)
		}

		if item.IsError && !strings.Contains(err.Error(), item.Err) {
			t.Errorf("Expecting [%v], got [%v]", item.Err, err)
		}
	}
}

func TestErrorJson(t *testing.T) {
	cases := []TestRequest{
		TestRequest{
			Request: SearchRequest{
				Limit:      10,
				Offset:     0,
				Query:      "Wolf",
				OrderField: wrongField,
				OrderBy:    OrderByAsIs,
			},
			Err:     "cant unpack error json",
			IsError: true,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(WrongErrorJsonServer))

	for _, item := range cases {
		srv := &SearchClient{
			AccessToken: securityToken,
			URL:         ts.URL,
		}
		_, err := srv.FindUsers(item.Request)

		if err == nil {
			t.Errorf("Expecting [%v], got nil", item.Err)
		}

		if item.IsError && !strings.Contains(err.Error(), item.Err) {
			t.Errorf("Expecting [%v], got [%v]", item.Err, err)
		}
	}
}

func TestResultJson(t *testing.T) {
	cases := []TestRequest{
		TestRequest{
			Request: SearchRequest{
				Limit:      10,
				Offset:     0,
				Query:      "Wolf",
				OrderField: "Id",
				OrderBy:    OrderByAsIs,
			},
			Err:     "cant unpack result json",
			IsError: true,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(WrongResultJsonServer))

	for _, item := range cases {
		srv := &SearchClient{
			AccessToken: securityToken,
			URL:         ts.URL,
		}
		_, err := srv.FindUsers(item.Request)

		if err == nil {
			t.Errorf("Expecting [%v], got nil", item.Err)
		}

		if item.IsError && !strings.Contains(err.Error(), item.Err) {
			t.Errorf("Expecting [%v], got [%v]", item.Err, err)
		}
	}
}

func TestResult(t *testing.T) {
	cases := []TestRequest{
		TestRequest{
			Request: SearchRequest{
				Limit:      2,
				Offset:     0,
				Query:      "Wolf",
				OrderField: "Id",
				OrderBy:    OrderByAsIs,
			},
			Err:     "",
			IsError: false,
		},
		TestRequest{
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "Rose",
				OrderField: "Id",
				OrderBy:    OrderByAsIs,
			},
			Err:     "",
			IsError: false,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	for _, item := range cases {
		srv := &SearchClient{
			AccessToken: securityToken,
			URL:         ts.URL,
		}
		_, err := srv.FindUsers(item.Request)

		if err != nil {
			t.Errorf("Expecting nil, got [%v]", err.Error())
		}
	}
}

/*
	go test -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html

*/
