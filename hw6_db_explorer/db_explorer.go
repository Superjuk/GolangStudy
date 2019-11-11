package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

// обращаю ваше внимание - в этом задании запрещены глобальные переменные
type DbApi struct {
	mu sync.Mutex
}

func NewDbExplorer(db *sql.DB) (*DbApi, error) {
	return &DbApi{}, nil
}

/*
* GET / - возвращает список все таблиц (которые мы можем использовать в дальнейших запросах)
* GET /$table?limit=5&offset=7 - возвращает список из 5 записей (limit) начиная с 7-й (offset) из таблицы $table. limit по-умолчанию 5, offset 0
* GET /$table/$id - возвращает информацию о самой записи или 404
* PUT /$table - создаёт новую запись, данный по записи в теле запроса (POST-параметры)
* POST /$table/$id - обновляет запись, данные приходят в теле запроса (POST-параметры)
* DELETE /$table/$id - удаляет запись
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
		// command := strings.Split(r.URL.Path, "/")
		// fmt.Println(command)
		// if len(command) <= 3 {
		// 	fmt.Println(r.URL.Query())
		// }
	}
}

func (db *DbApi) Read(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Read")
	cmd := strings.Split(r.URL.Path, "/")

	var table string
	var id string
	switch len(cmd) {
	case 2:
		table = cmd[1]
		if table == "" {
			fmt.Println("Get all tables")
		} else {
			fmt.Println("Table =", table)
		}
	case 3:
		table = cmd[1]
		id = cmd[2]
		fmt.Println("Table =", table, "; ID =", id)
	default:
		log.Println("Wrong command")
	}
}

func (db *DbApi) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create")
	cmd := strings.Split(r.URL.Path, "/")

	var table string
	switch len(cmd) {
	case 2:
		table = cmd[1]
		fmt.Println("Table =", table)
	default:
		log.Println("Wrong command")
	}
}

func (db *DbApi) Update(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update")
	cmd := strings.Split(r.URL.Path, "/")

	var table string
	var id string
	switch len(cmd) {
	case 3:
		table = cmd[1]
		id = cmd[2]
		fmt.Println("Table =", table, "; ID =", id)
	default:
		log.Println("Wrong command")
	}
}

func (db *DbApi) Delete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete")
	cmd := strings.Split(r.URL.Path, "/")

	var table string
	var id string
	switch len(cmd) {
	case 3:
		table = cmd[1]
		id = cmd[2]
		fmt.Println("Table =", table, "; ID =", id)
	default:
		log.Println("Wrong command")
	}
}
