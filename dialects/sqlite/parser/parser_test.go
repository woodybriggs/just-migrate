package sqlite

import (
	"fmt"
	"runtime"
	"testing"
	"woodybriggs/justmigrate/core/luther"
)

func makeParser(input string) *SqliteParser {

	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("unable to get caller info")
	}
	funcInfo := runtime.FuncForPC(pc)
	file, _ := funcInfo.FileLine(pc)

	lex := luther.NewLexer(luther.SourceCode{
		FileName: fmt.Sprintf("%s/%s", file, funcInfo.Name()),
		Raw:      []rune(input),
	})

	return NewSqliteParser(lex)
}

func TestCreateTable(t *testing.T) {
	parser := makeParser("CREATE TABLE IF NOT EXISTS users (id integer PRIMARY KEY AUTOINCREMENT)")

	createTable := parser.CreateTableStatement(false)
	fmt.Println(createTable)
}

func TestParseIdentifier(t *testing.T) {
	parser := makeParser("user_id [user_id] `user_id` \"user_id\"")

	for !parser.EndOfFile() {
		ident := parser.Identifier()
		if ident.Text != "user_id" {
			t.FailNow()
		}
	}
}
