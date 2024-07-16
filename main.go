package main

import (
	"fmt"
	"os"
	"strings"
)

func walk(node AstNode, depth int) {
	fmt.Printf("%s%T %+v\n", strings.Repeat("  ", depth), node, node)

	if node == nil {
		return
	} else {
		for _, child := range node.Children() {
			walk(child, depth+1)
		}
	}
}

func main() {

	filename := "./resources/test.sql"

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	lex, err := NewLexerFromFile(file)
	if err != nil {
		panic(err)
	}

	parser := Parser{
		lexer: lex,
	}

	stmts := parser.Statements()

	if !parser.HasErrors() {
		out := strings.Builder{}
		gen := NewSqlGenerator(stmts, out, nil)
		fmt.Print(gen.Generate())
	}
}
