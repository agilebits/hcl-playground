package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
)

type Demo struct {
	Name   string
	Params Params
}

type Params struct {
	URL string
	TLS TLS
}

type TLS struct {
	Enabled bool
	PEM     []string
}

func processNode(node ast.Node) {
	switch t := node.(type) {
	case *ast.File:
		processNode(t.Node)
	case *ast.ListType:
		for _, node := range t.List {
			processNode(node)
		}
	case *ast.ObjectType:
		processList(t.List)
	case *ast.ObjectList:
		processList(t)
	case *ast.ObjectItem:
		processItem(t)
	case *ast.LiteralType:
		if t.Token.Type == token.STRING {
			value, err := strconv.Unquote(t.Token.Text)
			if err == nil {
				t.Token.Text = strconv.Quote("***" + value + "***")
			}
		}
	default:
		fmt.Println("Unknown type: ", reflect.TypeOf(t))
	}
}

func processList(list *ast.ObjectList) {
	for _, item := range list.Items {
		fmt.Println(item.Val, item.Keys[0], reflect.TypeOf(item.Val))
		processItem(item)
	}
}

func processItem(item *ast.ObjectItem) {
	processNode(item.Val)
}

func main() {
	buf, err := ioutil.ReadFile("./demo.hcl")
	if err != nil {
		panic(err)
	}

	tree, err := hcl.Parse(string(buf))
	if err != nil {
		panic(err)
	}

	processNode(tree)

	var demo Demo
	if err := hcl.DecodeObject(&demo, tree); err != nil {
		panic(err)
	}

	fmt.Printf("Hello, %s with %q.\n", demo.Name, demo.Params.URL)
	for i, pem := range demo.Params.TLS.PEM {
		fmt.Printf("PEM[%d] = %q\n", i+1, pem)
	}
}
