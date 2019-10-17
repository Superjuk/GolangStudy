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

func main() {
	fset := token.NewFileSet()

	apiFile, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal("ParseFile error: ", err.Error())
	}

	//out, _ := os.Create(os.Args[2])

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

	// it work's
	fmt.Println("Declarations:")
	for _, decl := range apiFile.Decls {
		pos := decl.Pos()
		relPosition := fset.Position(pos)
		log.Println(relPosition.String())
	}

	/*API realization*/

}
