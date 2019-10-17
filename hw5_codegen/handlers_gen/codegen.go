package main

import (
	"os"
	// "go/ast"
	"go/parser"
	"go/token"
	"log"

	//"io/ioutil"
	"fmt"
)

// type ResponseOk struct {
// 	Error string      `json:"error"`
// 	Data  interface{} `json:"response"`
// }

const (
	responseErr = `type ResponseErr struct {
	Error string ` + "`" + `json:"error"` + "`" + `
}`

	responseOk = `type ResponseOk struct {
	Error string ` + "`" + `json:"error"` + "`" + `
	Data  interface{} ` + "`" + `json:"response"` + "`" + `
}`

	sendResponse = `func sendResponse(w http.ResponseWriter, err *ApiError, response interface{}) {
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
}`
)

func main() {
	fset := token.NewFileSet()

	apiFile, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal("ParseFile error: ", err.Error())
	}

	out, _ := os.Create(os.Args[2])

	fmt.Fprintln(out, `package `+apiFile.Name.Name)
	fmt.Fprintln(out)
	fmt.Fprintln(out, `import "context"`)
	fmt.Fprintln(out, `import "encoding/json"`)
	fmt.Fprintln(out, `import "fmt"`)
	fmt.Fprintln(out, `import "log"`)
	fmt.Fprintln(out, `import "net/http"`)
	fmt.Fprintln(out, `import "net/url"`)
	fmt.Fprintln(out, `import "strconv"`)
	fmt.Fprintln(out)
	fmt.Fprintln(out, responseErr)
	fmt.Fprintln(out)
	fmt.Fprintln(out, responseOk)
	fmt.Fprintln(out)
	fmt.Fprintln(out, sendResponse)

	// it work's
	fmt.Println("Declarations:")
	for _, decl := range apiFile.Decls {
		pos := decl.Pos()
		relPosition := fset.Position(pos)
		log.Println(relPosition.String())
	}

	/*API realization*/

}
