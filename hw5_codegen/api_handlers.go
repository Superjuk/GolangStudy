package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type ResponseErr struct {
	Error string `json:"error"`
}

type ResponseOk struct {
	Error string      `json:"error"`
	Data  interface{} `json:"response"`
}

func sendResponse(w http.ResponseWriter, err *ApiError, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	header := http.StatusOK

	send := func(resp interface{}) {
		jsonStr, _ := json.Marshal(resp)
		fmt.Printf("%s\n", jsonStr)

		w.WriteHeader(header)
		w.Write(jsonStr)
	}

	if err != nil {
		resp := ResponseErr{err.Err.Error()}
		header = err.HTTPStatus
		send(resp)
	} else if response != nil {
		resp := ResponseOk{}
		resp.Data = response
		send(resp)
	} else {
		log.Fatalln("Err and response equal nil. This is must not be.")
	}
}

/*MyApi*/
func (h *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/user/profile":
		h.handlerProfile(w, r)
	case "/user/create":
		h.handlerCreate(w, r)
	default:
		//		h.handlerUnknown(w, r)
	}
}

//! /user/profile
func (h *MyApi) handlerProfile(w http.ResponseWriter, r *http.Request) {
	// проверка метода
	if r.Method == "GET" {
		sendResponse(w, &ApiError{http.StatusNotAcceptable, fmt.Errorf("bad method")}, nil)
		return
	}

	// валидирование параметров
	params, errVal := h.validateProfileParams(r.URL.Query())
	if errVal != nil {
		sendResponse(w, errVal, nil)
		return
	}

	ctx := context.Background()
	res, err := h.Profile(ctx, *params)
	if err != nil {
		if ae, ok := err.(*ApiError); ok {
			sendResponse(w, ae, nil)
		} else {
			sendResponse(w, &ApiError{http.StatusInternalServerError, err}, nil)
		}
		return
	}

	// // прочие обработки
	//! \todo обработать context

	sendResponse(w, nil, res)
}

func (h *MyApi) validateProfileParams(query url.Values) (*ProfileParams, *ApiError) {
	out := &ProfileParams{}

	//Login
	out.Login = query.Get("login")
	//required
	if out.Login == "" {
		return nil, &ApiError{
			http.StatusBadRequest,
			fmt.Errorf("login must be not empty")}
	}

	return out, nil
}

//! /user/create
func (h *MyApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
	// проверка метода
	if r.Method == "GET" {
		sendResponse(w, &ApiError{http.StatusNotAcceptable, fmt.Errorf("bad method")}, nil)
		return
	}

	// валидирование параметров
	params, errVal := h.validateCreateParams(r.URL.Query())
	if errVal != nil {
		sendResponse(w, errVal, nil)
		return
	}

	ctx := context.Background()
	res, err := h.Create(ctx, *params)
	if err != nil {
		if ae, ok := err.(*ApiError); ok {
			sendResponse(w, ae, nil)
		} else {
			sendResponse(w, &ApiError{http.StatusInternalServerError, err}, nil)
		}
		return
	}

	// // прочие обработки
	//! \todo обработать context

	sendResponse(w, nil, res)
}

func (h *MyApi) validateCreateParams(query url.Values) (*CreateParams, *ApiError) {
	out := &CreateParams{}

	//Login
	out.Login = query.Get("login")
	//required
	if out.Login == "" {
		return nil, &ApiError{
			http.StatusBadRequest,
			fmt.Errorf("login must be not empty")}
	}
	//min length
	if len(out.Login) < 10 {
		return nil, &ApiError{
			http.StatusBadRequest,
			fmt.Errorf("login len must be >= 10")}
	}

	//Name
	out.Name = query.Get("full_name")

	//Status
	out.Status = query.Get("status")
	//default
	if out.Status == "" {
		out.Status = "user"
	}
	//enum
	switch out.Status {
	case "user", "moderator", "admin":
		break
	default:
		return nil, &ApiError{
			http.StatusBadRequest,
			fmt.Errorf("status must be one of [user, moderator, admin]")}
	}

	//Age
	var errAge error
	out.Age, errAge = strconv.Atoi(query.Get("age"))
	fmt.Println(out.Age)
	if errAge != nil {
		return nil, &ApiError{
			http.StatusBadRequest,
			fmt.Errorf("age must be int")}
	}
	//min
	if out.Age < 0 {
		return nil, &ApiError{
			http.StatusBadRequest,
			fmt.Errorf("age must be >= 0")}
	}
	//max
	if out.Age > 128 {
		return nil, &ApiError{
			http.StatusBadRequest,
			fmt.Errorf("age must be <= 128")}
	}

	return out, nil
}

// //! unknonw URL path
// func (h *MyApi) handlerUnknown(w http.ResponseWriter, r *http.Request) {

// }

// /*OtherApi*/
// func (o *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	switch r.URL.Path {
// 	case "/user/create":
// 		o.handlerCreate(w, r)
// 	default:
// 		ae := ApiError{http.StatusNotFound, errors.New("Not found")}
// 		ae.Error()
// 	}
// }

// //! /user/create
// func (o *OtherApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
// 	// заполнение структуры params
// 	raw := o.fillCreateParams(r.URL.Query())
// 	// валидирование параметров
// 	params, errVal := o.validateCreateParams(raw)
// 	if &errVal != nil {
// 		sendResponse(w, nil, &errVal)
// 		return
// 	}
// 	ctx := context.Background()
// 	res, err := o.Create(ctx, *params)
// 	if err != nil {
// 		sendResponse(w, nil, &ApiError{http.StatusNotFound, err})
// 		return
// 	}
// 	// прочие обработки
// 	result := `{"error": "", "response": {
// 		"id": ` + strconv.FormatUint(res.ID, 10) + `,
// 		"login": "` + res.Login + `",
// 		"full_name": "` + res.FullName + `",
// 		"level": ` + strconv.Itoa(res.Level) + `}}`

// 	sendResponse(w, &result, nil)
// }

// func (o *OtherApi) fillCreateParams(query url.Values) (out *OtherCreateParams) {
// 	out = &OtherCreateParams{}
// 	out.Username = query.Get("username")
// 	out.Name = query.Get("account_name")
// 	out.Class = query.Get("class")
// 	out.Level, _ = strconv.Atoi(query.Get("level"))
// 	return out
// }

// func (o *OtherApi) validateCreateParams(in *OtherCreateParams) (out *OtherCreateParams, err ApiError) {
// 	out = in
// 	//! required
// 	if out.Username == "" {
// 		err = ApiError{http.StatusBadRequest, errors.New("username must not be empty")}
// 		return out, err
// 	}
// 	//! min = 3
// 	if len(out.Username) < 3 {
// 		err = ApiError{http.StatusBadRequest, errors.New("username len must be >= 3")}
// 		return out, err
// 	}
// 	//! default=warrior
// 	if out.Class == "" {
// 		out.Class = "warrior"
// 	}
// 	//! enum=warrior|sorcerer|rouge
// 	switch out.Class {
// 	case "warrior", "sorcerer", "rouge":
// 		break
// 	default:
// 		err = ApiError{http.StatusBadRequest, errors.New("class must be one of [warrior, sorcerer, rouge]")}
// 		return out, err
// 	}
// 	//! min = 1
// 	if out.Level < 1 {
// 		err = ApiError{http.StatusBadRequest, errors.New("level must be >= 1")}
// 		return out, err
// 	}
// 	//! max = 50
// 	if out.Level > 50 {
// 		err = ApiError{http.StatusBadRequest, errors.New("level must be <= 50")}
// 		return out, err
// 	}

// 	err = ApiError{}
// 	return out, err
// }
