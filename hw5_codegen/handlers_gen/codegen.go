package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
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

	//ServeHTTP
	serveHttpBegin = `func (h {{.structType}}) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {`

	serveHttpCase = `	case "{{.apigenUrl}}":
		h.handler{{.apigenMethod}}(w, r)`

	serveHttpEnd = `	default:
		sendResponse(w, &ApiError{http.StatusNotFound, fmt.Errorf("unknown method")}, nil)
	}`

	//handler
	handlerBegin = `func (h {{.structType}}) handler{{.apigenMethod}}(w http.ResponseWriter, r *http.Request) {
	var query url.Values`

	handlerMethodGetTrue = `if r.Method == http.MethodGet {
		query = r.URL.Query()
	}`

	handlerMethodGetFalse = `if r.Method == http.MethodGet {
		sendResponse(w, &ApiError{http.StatusNotAcceptable, fmt.Errorf("bad method")}, nil)
		return
	}`

	handlerMethodPostTrue = `if r.Method == http.MethodPost {
		r.ParseForm()
		query = r.PostForm
	}`

	handlerMethodPostFalse = `if r.Method == http.MethodPost {
		sendResponse(w, &ApiError{http.StatusNotAcceptable, fmt.Errorf("bad method")}, nil)
		return
	}`

	handlerAuthTrue = `if r.Header.Get("X-Auth") != "100500" {
		sendResponse(w, &ApiError{http.StatusForbidden, fmt.Errorf("unauthorized")}, nil)
		return
	}`

	handlerEnd = `// валидирование параметров
	params, errVal := h.validate{{.apivalidatorStructType}}(query)
	if errVal != nil {
		sendResponse(w, errVal, nil)
		return
	}

	ctx := context.Background()
	res, err := h.{{.apigenMethod}}(ctx, *params)
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

	sendResponse(w, nil, res)`

	//validate
	validateStart = `func (h {{.structType}}) validate{{.apivalidatorStructType}}(query url.Values) (*{{.apivalidatorStructType}}, *ApiError) {
	out := &{{.apivalidatorStructType}}{}`

	validateFieldTypeString = `out.{{.apivalidatorFieldName}} = query.Get("{{.apivalidatorParamname}}")`
)

func main() {
	fset := token.NewFileSet()

	apigen, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal("ParseFile error: ", err.Error())
	}

	srcByte, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Print(err)
	}
	// convert bytes to string
	src := string(srcByte)

	// out, _ := os.Create(os.Args[2])

	// fmt.Fprintln(out, `package `+apiFile.Name.Name)
	// fmt.Fprintln(out)
	// fmt.Fprintln(out, `import "context"`)
	// fmt.Fprintln(out, `import "encoding/json"`)
	// fmt.Fprintln(out, `import "fmt"`)
	// fmt.Fprintln(out, `import "log"`)
	// fmt.Fprintln(out, `import "net/http"`)
	// fmt.Fprintln(out, `import "net/url"`)
	// fmt.Fprintln(out, `import "strconv"`)
	// fmt.Fprintln(out)
	// fmt.Fprintln(out, responseErr)
	// fmt.Fprintln(out)
	// fmt.Fprintln(out, responseOk)
	// fmt.Fprintln(out)
	// fmt.Fprintln(out, sendResponse)
	// fmt.Fprintln(out)

	// it work's
	fmt.Println("Declarations:")
	for _, decl := range apigen.Decls {
		// анализируем функции
		if gen, ok := decl.(*ast.FuncDecl); ok {
			if gen.Doc.Text() != "" {
				fmt.Println("Name:", gen.Name.String())
				fmt.Println("Doc:", gen.Doc.Text())
				if gen.Recv.NumFields() > 0 {
					start := gen.Recv.List[0].Type.Pos() - 1
					end := gen.Recv.List[0].Type.End() - 1
					fmt.Println("Recv:", src[start:end])
				}
				fmt.Println("@func@")
			}
		}

		// анализируем структуры
		if gen, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range gen.Specs {
				specType, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				fmt.Println("specType:", specType.Name.Name)
				structType, ok := specType.Type.(*ast.StructType)
				if !ok {
					continue
				}
				for _, field := range structType.Fields.List {
					for _, name := range field.Names {
						fmt.Println("Name:", name)
					}
					if field.Tag != nil {
						fmt.Println("Description:", field.Tag.Value)
					}

					fieldType, ok := field.Type.(*ast.Ident)
					if ok {
						fmt.Println("Type:", fieldType)
					}
				}
				fmt.Println("@struct@")
			}

			fmt.Println("--")
		}
	}

	// all comments

	// filtering

	/*API realization*/

}
