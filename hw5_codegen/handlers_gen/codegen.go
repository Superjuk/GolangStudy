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

	out, _ := os.Create(os.Args[2])

	fmt.Fprintln(out, `package `+apiFile.Name.Name)

	// // it work's
	// for _, decl := range apiFile.Decls {
	// 	pos := decl.Pos()
	// 	relPosition := fset.Position(pos)
	// 	log.Println(relPosition.String())
	// }

	/*API realization*/

}
