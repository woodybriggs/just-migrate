package main

import (
	"fmt"
	"os"
	"woodybriggs/justmigrate/core/luther"
	"woodybriggs/justmigrate/core/report"
	"woodybriggs/justmigrate/parser"
)

func main() {

	filename := "./resources/test.sql"

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	tokenizer, err := luther.NewLexerFromFile(file)
	if err != nil {
		panic(err)
	}

	parser := parser.NewParser(tokenizer)
	parser.Statements()

	renderer := report.Renderer{}

	for _, diag := range parser.Errors {
		fmt.Println(renderer.Render(diag))
	}
}
