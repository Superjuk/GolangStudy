package main

import (
	"context"
	//"fmt"
	//"context"
	//"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

/*out must be slice of strings or struct of strings*/
func fillProfileParams(query url.Values) (out *ProfileParams) {
	out = &ProfileParams{}
	out.Login = query.Get("login")
	return out
}

/*in must be slice of strings or struct of strings*/
func validateProfileParams(in *ProfileParams) (out *ProfileParams, err error) {
	out = in
	err = nil
	//required
	if in.Login == "" {
		err = errors.New("login must not me empty")
	}
	return out, err
}

/*out must be slice of strings or struct of strings*/
func fillCreateParams(query url.Values) (out *CreateParams) {
	out = &CreateParams{}
	out.Login = query.Get("login")
	out.Name = query.Get("full_name")
	out.Status = query.Get("status")
	out.Age, _ = strconv.Atoi(query.Get("age"))
	return out
}

/*in must be slice of strings or struct of strings*/
func validateCreateParams(in *CreateParams) (out *CreateParams, err error) {
	//required
	if in.Login == "" {
		err = errors.New("login must not me empty")
	}
	return out, err
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
	if err != nil {
		/* исправить вывод в соответствии с main_test.go */
		w.Write([]byte(err.Error()))
		return
	}
	// прочие обработки
	/* вывод должен быть в json */
	result := `{"error": "", "response": {
		"id": ` + strconv.FormatUint(res.ID, 10) + `,
		"login": "` + res.Login + `",
		"full_name": "` + res.FullName + `",
		"status": ` + strconv.Itoa(res.Status) + `}}`

	w.Write([]byte(result))
}

func (h *MyApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
	// заполнение структуры params
	raw := fillCreateParams(r.URL.Query())
	// валидирование параметров
	params, err := validateCreateParams(raw)
	if err != nil {
		/* исправить вывод в соответствии с main_test.go */
		w.Write([]byte(err.Error()))
		return
	}
	ctx := context.Background()
	res, err := h.Create(ctx, *params)
	if err != nil {
		/* исправить вывод в соответствии с main_test.go */
		w.Write([]byte(err.Error()))
		return
	}
	// прочие обработки
	/* вывод должен быть в json */
	result := `{"error": "", "response": {
		"id": ` + strconv.FormatUint(res.ID, 10) + `,
		"login": "` + res.Login + `",
		"full_name": "` + res.FullName + `",
		"status": ` + strconv.Itoa(res.Status) + `}}`

	w.Write([]byte(result))
}
