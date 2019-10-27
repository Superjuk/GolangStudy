package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"text/scanner"
	"text/template"
)

type TagType struct {
	Required   bool
	Paramname  string
	Enum       []string
	DefaultStr string
	Min        *int
	Max        *int
}

type StructField struct {
	Name string
	Type string
	Tag  TagType
}

//for json unmarshall
type ApigenApi struct {
	Url    string
	Auth   bool
	Method string
}

type ApigenData struct {
	Name   string
	InType string
	Url    string
	Auth   bool
	Method string
}

//------------------------------------
// type Apigen struct {
// 	Type   string
// 	Name   string
// 	InType string
// 	Url    string
// 	Auth   bool
// 	Method string
// }

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
	serveHttpBegin = `func (h {{template "TYPE"}}) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {` + "\n"

	serveHttpCase = `	case "{{.Url}}":
		h.handler{{.Name}}(w, r)` + "\n"

	serveHttpEnd = `	default:
		sendResponse(w, &ApiError{http.StatusNotFound, fmt.Errorf("unknown method")}, nil)
	}` + "\n}\n"

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

	// Парсим tag
	parseTag := func(tag string) (out TagType) {
		// cleaning
		firstClean := strings.Trim(tag, "`\"\"`")
		str := strings.Replace(firstClean, "apivalidator:\"", "", 1)
		str1 := strings.ReplaceAll(str, "=", " ")
		str2 := strings.ReplaceAll(str1, ",", " ")

		// Parsing
		var parser scanner.Scanner
		parser.Init(strings.NewReader(str2))

		for tok := parser.Scan(); tok != scanner.EOF; tok = parser.Scan() {
			switch parser.TokenText() {
			case "required":
				out.Required = true
			case "default":
				tok = parser.Scan()
				if tok != scanner.EOF {
					out.DefaultStr = parser.TokenText()
				}
			case "paramname":
				tok = parser.Scan()
				if tok != scanner.EOF {
					out.Paramname = parser.TokenText()
				}
			case "enum", "|":
				tok = parser.Scan()
				if tok != scanner.EOF {
					out.Enum = append(out.Enum, parser.TokenText())
				}
			case "min":
				tok = parser.Scan()
				if tok != scanner.EOF {
					if min, err := strconv.Atoi(parser.TokenText()); err == nil {
						out.Min = &min
					}
				}
			case "max":
				tok = parser.Scan()
				if tok != scanner.EOF {
					if max, err := strconv.Atoi(parser.TokenText()); err == nil {
						out.Max = &max
					}
				}
			default:
				break
			}
		}

		return
	}

	// Парсим api.go
	//var apigens []Apigen
	apigens := make(map[string]([]ApigenData))
	var apivalidators []Apivalidator

	for _, decl := range apigen.Decls {
		// анализируем функции
		if gen, ok := decl.(*ast.FuncDecl); ok {
			if gen.Doc.Text() != "" {
				var apigenKey string
				if gen.Recv.NumFields() > 0 {
					start := gen.Recv.List[0].Type.Pos() - 1
					end := gen.Recv.List[0].Type.End() - 1
					if end > start {
						apigenKey = src[start:end]
					} else {
						fmt.Println("Type apigen error: start = end")
						continue
					}
				}

				var apigenFunc ApigenData
				apigenFunc.Name = gen.Name.String()

				var apgn ApigenApi
				doc := gen.Doc.Text()
				if strings.HasPrefix(doc, "apigen:api ") {
					apigenStr := strings.TrimPrefix(doc, "apigen:api ")
					err := json.Unmarshal([]byte(apigenStr), &apgn)
					if err != nil {
						fmt.Println("Json apigen error:", err.Error())
						continue
					}
					apigenFunc.Auth = apgn.Auth
					apigenFunc.Method = apgn.Method
					apigenFunc.Url = apgn.Url
				}

				for _, p := range gen.Type.Params.List {
					for _, in := range p.Names {
						if in.Name == "in" {
							start := p.Type.Pos() - 1
							end := p.Type.End() - 1
							if end > start {
								apigenFunc.InType = src[start:end]
							} else {
								fmt.Println("InType apigen error: start = end")
								continue
							}
						}
					}
				}

				apigens[apigenKey] = append(apigens[apigenKey], apigenFunc)
			}
			continue
		}

		// анализируем структуры
		if gen, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range gen.Specs {
				var apival Apivalidator
				specType, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				apival.Type = specType.Name.Name

				structType, ok := specType.Type.(*ast.StructType)
				if !ok {
					continue
				}

				var apivalFields []StructField
				tagsExist := false
				for _, field := range structType.Fields.List {
					var apivalField StructField
					tag := field.Tag
					if tag != nil {
						if !strings.Contains(tag.Value, "apivalidator:") {
							continue
						}
						tagsExist = true
						apivalField.Tag = parseTag(field.Tag.Value)
					}
					for _, name := range field.Names {
						apivalField.Name = name.Name
						break
					}
					if apivalField.Tag.Paramname == "" {
						apivalField.Tag.Paramname = strings.ToLower(apivalField.Name)
					}
					fieldType, ok := field.Type.(*ast.Ident)
					if ok {
						apivalField.Type = fieldType.Name
					}
					apivalFields = append(apivalFields, apivalField)
				}

				if len(apivalFields) > 0 && tagsExist {
					apival.Fields = apivalFields
					apivalidators = append(apivalidators, apival)
				}
			}
		}
	}

	/*API realization*/
	out, _ := os.Create(os.Args[2])

	fmt.Fprintln(out, `package `+apigen.Name.Name)
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
	fmt.Fprintln(out)

	/*
	   // Create a new template and parse the letter into it.
	   t := template.Must(template.New("letter").Parse(letter))

	   // Execute the template for each recipient.
	   for _, r := range recipients {
	       err := t.Execute(os.Stdout, r)
	       if err != nil {
	           log.Println("executing template:", err)
	       }
	   }
	*/
	/*
		`{{define "T1"}}ONE{{end}}
		{{define "T2"}}TWO{{end}}
		{{define "T3"}}{{template "T1"}} {{template "T2"}}{{end}}
		{{template "T3"}}`
	*/
	// generate serveHTTP
	serveHttpCaseTmpl := template.Must(template.New("serveHttpCase").Parse(serveHttpCase))
	for key, _ := range apigens {
		serveHttpBeginTmpl := template.Must(template.New("serveHttpBegin").Parse(`{{define "TYPE"}}` + key + `{{end}}` + serveHttpBegin))
		err := serveHttpBeginTmpl.Execute(out, "")
		if err != nil {
			log.Fatalln("ServeHttp gen err =", err.Error())
		}
		for _, sh := range apigens[key] {
			err = serveHttpCaseTmpl.Execute(out, sh)
			if err != nil {
				log.Fatalln("ServeHttp gen err =", err.Error())
			}
		}

		fmt.Fprintln(out, serveHttpEnd)
	}
}
