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

/*MyApi*/
/*out must be slice of strings or struct of strings*/
func (h *MyApi) fillProfileParams(query url.Values) (out *ProfileParams) {
	out = &ProfileParams{}
	out.Login = query.Get("login")
	return out
}

/*in must be slice of strings or struct of strings*/
func (h *MyApi) validateProfileParams(in *ProfileParams) (out *ProfileParams, err error) {
	out = in
	err = nil
	//required
	if in.Login == "" {
		err = errors.New("login must not be empty")
	}
	return out, err
}

/*out must be slice of strings or struct of strings*/
func (h *MyApi) fillCreateParams(query url.Values) (out *CreateParams) {
	out = &CreateParams{}
	out.Login = query.Get("login")
	out.Name = query.Get("full_name")
	out.Status = query.Get("status")
	out.Age, _ = strconv.Atoi(query.Get("age"))
	return out
}

/*in must be slice of strings or struct of strings*/
func (h *MyApi) validateCreateParams(in *CreateParams) (out *CreateParams, err error) {
	out = in
	err = nil
	//! required
	if out.Login == "" {
		err = errors.New("login must not be empty")
	}
	//! min = 10
	if len(out.Login) < 10 {
		err = errors.New("login len must be >= 10")
	}
	//! default=user
	if out.Status == "" {
		out.Status = "user"
	}
	//! enum=user|moderator|admin
	switch out.Status {
	case "user", "moderator", "admin":
		break
	default:
		err = errors.New("status must be one of [user, moderator, admin]")
	}
	//! min = 0
	if out.Age < 0 {
		err = errors.New("age must be >= 0")
	}
	//! max = 128
	if out.Age > 128 {
		err = errors.New("age must be <= 128")
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
	raw := h.fillProfileParams(r.URL.Query())
	// валидирование параметров
	params, err := h.validateProfileParams(raw)
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
	raw := h.fillCreateParams(r.URL.Query())
	// валидирование параметров
	params, err := h.validateCreateParams(raw)
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
		"id": ` + strconv.FormatUint(res.ID, 10) + `}}`

	w.Write([]byte(result))
}

/*OtherApi*/
func (o *OtherApi) fillCreateParams(query url.Values) (out *OtherCreateParams) {
	out = &OtherCreateParams{}
	out.Username = query.Get("username")
	out.Name = query.Get("account_name")
	out.Class = query.Get("class")
	out.Level, _ = strconv.Atoi(query.Get("level"))
	return out
}

func (o *OtherApi) validateCreateParams(in *OtherCreateParams) (out *OtherCreateParams, err error) {
	out = in
	err = nil
	//! required
	if out.Username == "" {
		err = errors.New("username must not be empty")
	}
	//! min = 3
	if len(out.Username) < 3 {
		err = errors.New("username len must be >= 3")
	}
	//! default=warrior
	if out.Class == "" {
		out.Class = "warrior"
	}
	//! enum=warrior|sorcerer|rouge
	switch out.Class {
	case "warrior", "sorcerer", "rouge":
		break
	default:
		err = errors.New("class must be one of [warrior, sorcerer, rouge]")
	}
	//! min = 1
	if out.Level < 1 {
		err = errors.New("level must be >= 1")
	}
	//! max = 50
	if out.Level > 50 {
		err = errors.New("level must be <= 50")
	}
	return out, err
}

func (o *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/user/create":
		o.handlerCreate(w, r)
	default:
		ae := ApiError{404, errors.New("Not found")}
		ae.Error()
	}
}

func (o *OtherApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
	// заполнение структуры params
	raw := o.fillCreateParams(r.URL.Query())
	// валидирование параметров
	params, err := o.validateCreateParams(raw)
	if err != nil {
		/* исправить вывод в соответствии с main_test.go */
		w.Write([]byte(err.Error()))
		return
	}
	ctx := context.Background()
	res, err := o.Create(ctx, *params)
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
		"level": ` + strconv.Itoa(res.Level) + `}}`

	w.Write([]byte(result))
}
