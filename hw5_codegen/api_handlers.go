package main

import (
	"context"
	//"fmt"
	//"context"
	//"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

func fillProfileParams(query url.Values) (out *ProfileParams) {
	out = &ProfileParams{}
	out.Login = query.Get("login")
	return out
}

func validateProfileParams(in *ProfileParams) (out *ProfileParams, err error) {
	out = in
	err = nil
	//required
	if in.Login == "" {
		err = errors.New("login must not me empty")
	}
	return out, err
}

func validateCreateParams() {

}

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
	raw := fillProfileParams(r.URL.Query())
	// валидирование параметров
	params, err := validateProfileParams(raw)
	if err != nil {
		/* исправить вывод в соответствии с main_test.go */
		w.Write([]byte(err.Error()))
		return
	}
	ctx := context.Background()
	res, err := h.Profile(ctx, *params)
	// прочие обработки
	/* вывод должен быть в json */
	result := `{"error": "",
				"response": {
					"id":        res.ID,
					"login":     res.Login,
					"full_name": res.FullName,
					"status":    res.Status,
		},}`
}

func (h *MyApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
	// заполнение структуры params
	// валидирование параметров
	//res, err := h.Create(ctx, params)
	// прочие обработки
}
