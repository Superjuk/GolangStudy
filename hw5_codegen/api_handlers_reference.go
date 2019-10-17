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
		sendResponse(w, &ApiError{http.StatusNotFound, fmt.Errorf("unknown method")}, nil)
	}
}

//! /user/profile
func (h *MyApi) handlerProfile(w http.ResponseWriter, r *http.Request) {
	var query url.Values
	if r.Method == http.MethodGet {
		query = r.URL.Query()
	}
	if r.Method == http.MethodPost {
		r.ParseForm()
		query = r.PostForm
	}
	// валидирование параметров
	params, errVal := h.validateProfileParams(query)
	if errVal != nil {
		sendResponse(w, errVal, nil)
		return
	}

	ctx := context.Background()
	res, err := h.Profile(ctx, *params)
	if err != nil {
		if ae, ok := err.(ApiError); ok {
			sendResponse(w, &ae, nil)
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
	var query url.Values
	if r.Method == http.MethodGet {
		sendResponse(w, &ApiError{http.StatusNotAcceptable, fmt.Errorf("bad method")}, nil)
		return
	}
	if r.Method == http.MethodPost {
		r.ParseForm()
		query = r.PostForm
	}

	// проверка авторизации
	if r.Header.Get("X-Auth") != "100500" {
		sendResponse(w, &ApiError{http.StatusForbidden, fmt.Errorf("unauthorized")}, nil)
		return
	}

	// валидирование параметров
	params, errVal := h.validateCreateParams(query)
	if errVal != nil {
		sendResponse(w, errVal, nil)
		return
	}

	ctx := context.Background()
	res, err := h.Create(ctx, *params)
	if err != nil {
		if ae, ok := err.(ApiError); ok {
			sendResponse(w, &ae, nil)
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

// /*OtherApi*/
func (h *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/user/create":
		h.handlerCreate(w, r)
	default:
		sendResponse(w, &ApiError{http.StatusNotFound, fmt.Errorf("unknown method")}, nil)
	}
}

//! /user/create
func (h *OtherApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
	// проверка метода
	var query url.Values
	if r.Method == http.MethodGet {
		sendResponse(w, &ApiError{http.StatusNotAcceptable, fmt.Errorf("bad method")}, nil)
		return
	}
	if r.Method == http.MethodPost {
		r.ParseForm()
		query = r.PostForm
	}

	// проверка авторизации
	if r.Header.Get("X-Auth") != "100500" {
		sendResponse(w, &ApiError{http.StatusForbidden, fmt.Errorf("unauthorized")}, nil)
		return
	}

	// валидирование параметров
	params, errVal := h.validateCreateParams(query)
	if errVal != nil {
		sendResponse(w, errVal, nil)
		return
	}

	ctx := context.Background()
	res, err := h.Create(ctx, *params)
	if err != nil {
		if ae, ok := err.(ApiError); ok {
			sendResponse(w, &ae, nil)
		} else {
			sendResponse(w, &ApiError{http.StatusInternalServerError, err}, nil)
		}
		return
	}

	// // прочие обработки
	//! \todo обработать context

	sendResponse(w, nil, res)
}

func (h *OtherApi) validateCreateParams(query url.Values) (*OtherCreateParams, *ApiError) {
	out := &OtherCreateParams{}

	//Username
	out.Username = query.Get("username")
	//required
	if out.Username == "" {
		return nil, &ApiError{
			http.StatusBadRequest,
			fmt.Errorf("username must be not empty")}
	}
	//min length
	if len(out.Username) < 3 {
		return nil, &ApiError{
			http.StatusBadRequest,
			fmt.Errorf("username len must be >= 3")}
	}

	//Name
	out.Name = query.Get("account_name")

	//Class
	out.Class = query.Get("class")
	//default
	if out.Class == "" {
		out.Class = "warrior"
	}
	//enum
	switch out.Class {
	case "warrior", "sorcerer", "rouge":
		break
	default:
		return nil, &ApiError{
			http.StatusBadRequest,
			fmt.Errorf("class must be one of [warrior, sorcerer, rouge]")}
	}

	//Level
	var errLevel error
	out.Level, errLevel = strconv.Atoi(query.Get("level"))
	fmt.Println(out.Level)
	if errLevel != nil {
		return nil, &ApiError{
			http.StatusBadRequest,
			fmt.Errorf("level must be int")}
	}
	//min
	if out.Level < 1 {
		return nil, &ApiError{
			http.StatusBadRequest,
			fmt.Errorf("level must be >= 0")}
	}
	//max
	if out.Level > 50 {
		return nil, &ApiError{
			http.StatusBadRequest,
			fmt.Errorf("level must be <= 50")}
	}

	return out, nil
}
