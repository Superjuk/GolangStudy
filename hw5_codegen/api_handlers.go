package main

import (
	//"context"
	"errors"
	"net/http"
)

//type MyApi struct{}

func (h *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/user/profile":
		h.handlerProfile(w, r)
	case "/user/create":
		h.handlerCreate(w, r)
	default:
		ae := ApiError{404, errors.New("Not found")}
		ae.Error()
	}
}

func (h *MyApi) handlerProfile(w http.ResponseWriter, r *http.Request) {
	// заполнение структуры params
	// валидирование параметров
	//res, err := h.DoSomeJob(ctx, params)
	// прочие обработки
}

func (h *MyApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
	// заполнение структуры params
	// валидирование параметров
	//res, err := h.DoSomeJob(ctx, params)
	// прочие обработки
}
