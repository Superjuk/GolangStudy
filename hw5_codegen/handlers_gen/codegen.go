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

/*
***original***
//apivalidator
type CreateParams struct {
	Login  string `apivalidator:"required,min=10"`
	Name   string `apivalidator:"paramname=full_name"`
	Status string `apivalidator:"enum=user|moderator|admin,default=user"`
	Age    int    `apivalidator:"min=0,max=128"`
}

//json
type User struct {
	ID       uint64 `json:"id"`
	Login    string `json:"login"`
	FullName string `json:"full_name"`
	Status   int    `json:"status"`
}

// apigen:api {"url": "/user/create", "auth": true, "method": "POST"}
func (srv *MyApi) Create(ctx context.Context, in CreateParams) (*NewUser, error)
------------------------------------------------------------------------------------
***modified***
//apivalidator
type {{GenDecl.specType} | {FuncDecl.InParam}} struct {
	{{GenDecl.structType.Name}}  {{GenDecl.structType.Type}} {{GenDecl.structType.Tag(`apivalidator:`)}}
}

//json
type {{GenDecl.specType}} struct {
	{{GenDecl.structType.Name}}  {{GenDecl.structType.Type}} {{GenDecl.structType.Tag(`json:`)}}
}

//apigen
// {{FuncDecl.Doc}}
func (srv {{FuncDecl.Recv(src[:])}}) {{FuncDecl.Name}}(ctx context.Context, in {{FuncDecl.InParam}}) (*NewUser, error)
*/

type StructField struct {
	Name string
	Type string
	Tag  string
}

type ApigenApi struct {
	url    string
	auth   bool
	method string
}

//------------------------------------
type Apigen struct {
	Type   string
	Name   string
	InType string
	Api    ApigenApi
}

type Apivalidator struct {
	Type   string
	Fields []StructField
}

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

	//var apigens []Apigen
	//var apivalidators []Apivalidator

	fmt.Println("Declarations:")
	for _, decl := range apigen.Decls {
		// анализируем функции
		if gen, ok := decl.(*ast.FuncDecl); ok {
			if gen.Doc.Text() != "" {
				fmt.Println("FuncDecl.Name:", gen.Name.String())
				fmt.Println("FuncDecl.Doc:", gen.Doc.Text())
				if gen.Recv.NumFields() > 0 {
					start := gen.Recv.List[0].Type.Pos() - 1
					end := gen.Recv.List[0].Type.End() - 1
					fmt.Println("FuncDecl.Recv(src[:]):", src[start:end])
				}
				for _, p := range gen.Type.Params.List {
					for _, in := range p.Names {
						if in.Name == "in" {
							start := p.Type.Pos() - 1
							end := p.Type.End() - 1
							fmt.Println("FuncDecl.InParam:", src[start:end])
						}
					}
				}
				//fmt.Println("FuncDecl.Out:", src[st:en])
				fmt.Println("@func@")
				fmt.Println("--")
			}
			continue
		}

		// анализируем структуры
		if gen, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range gen.Specs {
				specType, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				fmt.Println("GenDecl.specType:", specType.Name.Name)
				structType, ok := specType.Type.(*ast.StructType)
				if !ok {
					continue
				}
				for _, field := range structType.Fields.List {
					for _, name := range field.Names {
						fmt.Println("GenDecl.structType.Name:", name)
					}
					if field.Tag != nil {
						fmt.Println("GenDecl.structType.Tag:", field.Tag.Value)
					}

					fieldType, ok := field.Type.(*ast.Ident)
					if ok {
						fmt.Println("GenDecl.structType.Type:", fieldType)
					}
				}
				fmt.Println("@struct@")
				fmt.Println("--")
			}
		}
	}

	// all comments

	// filtering

	/*API realization*/

}
