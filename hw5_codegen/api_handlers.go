package main

import (
	//"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	//"strconv"
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

	resp := ResponseOk{}
	header := http.StatusOK

	if err != nil {
		resp.Error = err.Err.Error()
		header = err.HTTPStatus
	} else if response != nil {
		resp.Data = response
	} else {
		log.Fatalln("Err and response equal nil. This is must not be.")
	}

	jsonStr, _ := json.Marshal(resp)
	fmt.Printf("%s\n", jsonStr)

	w.WriteHeader(header)
	// if err != nil {
	// 	w.WriteHeader(err.HTTPStatus)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }

	// if json != nil {
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte(*json))
	// }
}

/*MyApi*/
func (h *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/user/profile":
		h.handlerProfile(w, r)
	case "/user/create":
		//		h.handlerCreate(w, r)
	default:
		//		h.handlerUnknown(w, r)
	}
}

//! /user/profile
func (h *MyApi) handlerProfile(w http.ResponseWriter, r *http.Request) {
	// проверка метода
	if r.Method == "GET" {
		sendResponse(w /*&ApiError{http.StatusNotAcceptable, fmt.Errorf("bad method")}, &NewUser{43}*/, nil, nil)
	}
	// // заполнение структуры params
	// raw := h.fillProfileParams(r.URL.Query())
	// // валидирование параметров
	// params, errVal := h.validateProfileParams(raw)
	// if errVal != nil {
	// 	sendResponse(w, nil, &errVal)
	// 	return
	// }
	// ctx := context.Background()
	// res, err := h.Profile(ctx, *params)
	// if err != nil {
	// 	sendResponse(w, nil, &ApiError{http.StatusNotFound, err})
	// 	return
	// }
	// // прочие обработки
	// result := `{"error": "", "response": {
	// 	"id": ` + strconv.FormatUint(res.ID, 10) + `,
	// 	"login": "` + res.Login + `",
	// 	"full_name": "` + res.FullName + `",
	// 	"status": ` + strconv.Itoa(res.Status) + `}}`

	// sendResponse(w, &result, nil)
}

func (h *MyApi) fillProfileParams(query url.Values) (out *ProfileParams) {
	out = &ProfileParams{}
	out.Login = query.Get("login")
	return out
}

func (h *MyApi) validateProfileParams(in *ProfileParams) (out *ProfileParams, err *ApiError) {
	out = in
	//required
	if out.Login == "" {
		err = &ApiError{http.StatusBadRequest, fmt.Errorf("login must not be empty")}
		return out, err
	}

	return out, nil
}

// //! /user/create
// func (h *MyApi) handlerCreate(w http.ResponseWriter, r *http.Request) {
// 	emptyErr := ApiError{}
// 	// заполнение структуры params
// 	raw := h.fillCreateParams(r.URL.Query())
// 	// валидирование параметров
// 	params, errVal := h.validateCreateParams(raw)
// 	if errVal != emptyErr {
// 		sendResponse(w, nil, &errVal)
// 		return
// 	}
// 	ctx := context.Background()
// 	res, err := h.Create(ctx, *params)
// 	if err != nil {
// 		//tempErr :=
// 		sendResponse(w, nil, &ApiError{http.StatusNotFound, err})
// 		return
// 	}
// 	// прочие обработки
// 	// \todo json marshaler and unmarshaler
// 	/* вывод должен быть в json */
// 	result := `{"error": "", "response": {
// 		"id": ` + strconv.FormatUint(res.ID, 10) + `}}`

// 	sendResponse(w, &result, nil)
// }

// func (h *MyApi) fillCreateParams(query url.Values) (out *CreateParams) {
// 	out = &CreateParams{}
// 	out.Login = query.Get("login")
// 	out.Name = query.Get("full_name")
// 	out.Status = query.Get("status")
// 	out.Age, _ /*errage*/ = strconv.Atoi(query.Get("age"))
// 	// if err_age != nil {
// 	// 	err = ApiError{http.StatusBadRequest, errors.New("age must be int")}
// 	// 	return out
// 	// }
// 	return out
// }

// func (h *MyApi) validateCreateParams(in *CreateParams) (out *CreateParams, err ApiError) {
// 	out = in
// 	//! required
// 	if out.Login == "" {
// 		err = ApiError{http.StatusBadRequest, errors.New("login must not be empty")}
// 		return out, err
// 	}
// 	//! min = 10
// 	if len(out.Login) < 10 {
// 		err = ApiError{http.StatusBadRequest, errors.New("login len must be >= 10")}
// 		return out, err
// 	}
// 	//! default=user
// 	if out.Status == "" {
// 		out.Status = "user"
// 	}
// 	//! enum=user|moderator|admin
// 	switch out.Status {
// 	case "user", "moderator", "admin":
// 		break
// 	default:
// 		err = ApiError{http.StatusBadRequest, errors.New("status must be one of [user, moderator, admin]")}
// 		return out, err
// 	}
// 	//! min = 0
// 	if out.Age < 0 {
// 		err = ApiError{http.StatusBadRequest, errors.New("age must be >= 0")}
// 		return out, err
// 	}
// 	//! max = 128
// 	if out.Age > 128 {
// 		err = ApiError{http.StatusBadRequest, errors.New("age must be <= 128")}
// 		return out, err
// 	}

// 	err = ApiError{}
// 	return out, err
// }

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
