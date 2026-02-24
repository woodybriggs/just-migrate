package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"maps"
	"os"
	"slices"

	"woodybriggs/justmigrate/core/ast"
	"woodybriggs/justmigrate/core/luther"
	"woodybriggs/justmigrate/core/report"
	"woodybriggs/justmigrate/database"
	sqlite "woodybriggs/justmigrate/dialects/sqlite/generator"
	"woodybriggs/justmigrate/diff"
	"woodybriggs/justmigrate/parser"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrInvalidNode = errors.New("invalid ast node")
)

var (
	ErrParserErrors = errors.New("parser has errors")
)

func assert(cond bool, err error) {
	if !cond {
		panic(err)
	}
}

type Database interface {
	Url() string
	ExportDataDefinitions() (string, error)
}

func ShowErrors(errors []report.Report, w io.Writer) {
	errorRenderer := report.Renderer{}
	for _, report := range errors {
		w.Write([]byte(errorRenderer.Render(report)))
	}
}

func ShowWarnings(warnings []report.Report, w io.Writer) {
	renderer := report.Renderer{}
	for _, report := range warnings {
		w.Write([]byte(renderer.Render(report)))
	}
}

func AstFromDatabase(database Database) (luther.SourceCode, []ast.Statement, error) {
	source, err := database.ExportDataDefinitions()
	if err != nil {
		return luther.SourceCode{}, nil, err
	}

	lexer := luther.NewLexer(
		luther.SourceCode{
			FileName: database.Url(),
			Raw:      []rune(source),
		},
	)

	parser := parser.NewParser(lexer)

	nodes := parser.Statements()
	errors := slices.Collect(maps.Values(parser.Errors))
	if len(errors) > 0 {
		ShowErrors(errors, os.Stderr)
		return parser.Lexer.SourceCode, nil, ErrParserErrors
	}

	// warnings := slices.Collect(maps.Values(parser.Warnings))
	// if len(warnings) > 0 {
	// 	ShowWarnings(warnings, os.Stderr)
	// }

	return parser.Lexer.SourceCode, nodes, nil
}

func AstFromFile(file *os.File) (luther.SourceCode, []ast.Statement, error) {
	lexer, err := luther.NewLexerFromFile(file)
	if err != nil {
		return lexer.SourceCode, nil, err
	}

	parser := parser.NewParser(lexer)

	nodes := parser.Statements()
	errors := slices.Collect(maps.Values(parser.Errors))

	if len(errors) > 0 {
		ShowErrors(errors, os.Stderr)
		return lexer.SourceCode, nil, ErrParserErrors
	}

	// warnings := slices.Collect(maps.Values(parser.Warnings))
	// if len(warnings) > 0 {
	// 	ShowWarnings(warnings, os.Stderr)
	// }

	return lexer.SourceCode, nodes, nil
}

func main() {

	var err error

	databaseURL := "/Users/woodybriggs/Projects/ts/currx/database/local.db"

	conn, err := sql.Open("sqlite3", databaseURL)
	if err != nil {
		log.Panicln(err)
	}

	db := &database.Sqlite{DB: conn, FileName: databaseURL}

	fileName := "resources/schema.sql"
	file, err := os.Open(fileName)
	if err != nil {
		os.Exit(1)
	}

	_, dstAst, err := AstFromFile(file)
	if err != nil {
		os.Exit(1)
	}

	_, srcAst, err := AstFromDatabase(db)
	if err != nil {
		os.Exit(1)
	}

	differ := diff.Diff{}

	edits, err := differ.DiffSchema(srcAst, dstAst)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	generate := sqlite.NewSqliteGenerator(edits)

	generate.Generate(os.Stderr)
}
