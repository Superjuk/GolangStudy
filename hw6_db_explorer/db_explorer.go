package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	//"net/url"
	//"sync"
)

type ApiError struct {
	HTTPStatus int
	Err        error
}

func (ae ApiError) Error() string {
	return ae.Err.Error()
}

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

// обращаю ваше внимание - в этом задании запрещены глобальные переменные
type DbApi struct {
	//mu sync.Mutex
	db *sql.DB
}

func NewDbExplorer(db *sql.DB) (*DbApi, error) {
	db.SetMaxOpenConns(5)

	err := db.Ping()
	if err != nil {
		log.Fatalln("DB not connected")
	}

	log.Println("DB connected")

	return &DbApi{db}, nil
}

/*

* GET, PUT, POST, DELETE - это http-метод, которым был отправлен запрос
 */

func (db *DbApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		db.Read(w, r)
	case http.MethodPut:
		db.Create(w, r)
	case http.MethodPost:
		db.Update(w, r)
	case http.MethodDelete:
		db.Delete(w, r)
	default:
		log.Println("Unknown method")
	}
}

func (db *DbApi) Read(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Read")
	cmd := strings.Split(r.URL.Path, "/")

	var table string
	var id string
	switch len(cmd) {
	/* GET / - возвращает список все таблиц (которые мы можем использовать в дальнейших запросах) */
	case 2:
		table = cmd[1]
		if table == "" {
			fmt.Println("Get all tables...")
			rows, err := db.db.Query("SHOW TABLES")
			if err != nil {
				log.Println("Error on show table's list:", err.Error())
				return
			}
			for rows.Next() {
				var s string
				err = rows.Scan(&s)
				if err != nil {
					log.Println("Error on load rows:", err.Error())
				}
				fmt.Println(s)
			}
			rows.Close()
			/* GET /$table?limit=5&offset=7 - возвращает список из 5 записей (limit) начиная с 7-й (offset) из таблицы $table. limit по-умолчанию 5, offset 0 */
		} else {
			fmt.Println("Table =", table)
		}
	/* GET /$table/$id - возвращает информацию о самой записи или 404 */
	case 3:
		table = cmd[1]
		id = cmd[2]
		fmt.Println("Table =", table, "; ID =", id)
	default:
		log.Println("Wrong command")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, Curl!"))
}

/* PUT /$table - создаёт новую запись, данный по записи в теле запроса (POST-параметры) */
func (db *DbApi) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create")
	cmd := strings.Split(r.URL.Path, "/")

	if len(cmd) != 2 {
		log.Println("Wrong command")
		return
	}

	table := cmd[1]
	fmt.Println("Table =", table)
}

/* POST /$table/$id - обновляет запись, данные приходят в теле запроса (POST-параметры) */
func (db *DbApi) Update(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update")
	cmd := strings.Split(r.URL.Path, "/")

	if len(cmd) != 3 {
		log.Println("Wrong command")
		return
	}

	table := cmd[1]
	id := cmd[2]
	fmt.Println("Table =", table, "; ID =", id)
}

/* DELETE /$table/$id - удаляет запись */
func (db *DbApi) Delete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete")
	cmd := strings.Split(r.URL.Path, "/")

	if len(cmd) != 3 {
		log.Println("Wrong command")
		return
	}

	table := cmd[1]
	id := cmd[2]
	fmt.Println("Table =", table, "; ID =", id)
}
