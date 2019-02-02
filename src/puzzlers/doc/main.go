package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"log"
	"os"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/types/typeutil"
	// TODO: these will use std go/types after Feb 2016
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: doc <package> <object>")
	}

	pkgpath, name := os.Args[1], os.Args[2]

	fmt.Printf("pkgpath %s,name %s \n", pkgpath, name)

	conf := loader.Config{ParserMode: parser.ParseComments}
	conf.Import(pkgpath)
	lprog, err := conf.Load()
	if err != nil {
		log.Fatal(err)
	}

	pkg := lprog.Package(pkgpath).Pkg
	obj := pkg.Scope().Lookup(name)

	if obj == nil {
		log.Fatalf("%s.%s not found", pkg.Path(), name)
	}

	fmt.Println(obj)

	for _, sel := range typeutil.IntuitiveMethodSet(obj.Type(), nil) {
		fmt.Printf("%s: %s\n", lprog.Fset.Position(sel.Obj().Pos()), sel)
	}

	// Find the path from the root of the AST to the object's position.
	// Walk up to the enclosing ast.Decl for the doc comment.
	_, path, _ := lprog.PathEnclosingInterval(obj.Pos(), obj.Pos())
	for _, n := range path {
		switch n := n.(type) {
		case *ast.GenDecl:
			fmt.Println("\n", n.Doc.Text())
			return
		case *ast.FuncDecl:
			fmt.Println("\n", n.Doc.Text())
			return
		}
	}
}
