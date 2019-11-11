package main

import (
	"database/sql"
	"fmt"
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

func (db *DbApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Println("It works!!!")
	default:
		command := strings.Split(r.URL.Path, "/")
		fmt.Println(command)
		if len(command) <= 3 {
			fmt.Println(r.URL.Query())
		}
	}
}
