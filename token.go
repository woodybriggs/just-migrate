package main

import "fmt"

type TokenKind int

var tokenKindString = map[TokenKind]string{
	TokenKind_EOF:                 "eof",
	TokenKind_LParen:              "l-paren",
	TokenKind_RParen:              "r-paren",
	TokenKind_Comma:               "comma",
	TokenKind_Period:              "period",
	TokenKind_SemiColon:           "semi-colon",
	TokenKind_Plus:                "plus",
	TokenKind_Minus:               "minus",
	TokenKind_Backtic:             "backtic",
	TokenKind_Identifier:          "identifier",
	TokenKind_NumericLiteral:      "numeric-literal",
	TokenKind_HexNumericLiteral:   "hexnumeric-literal",
	TokenKind_Comment:             "comment",
	TokenKind_Keyword_CREATE:      "create",
	TokenKind_Keyword_TEMPORARY:   "temporary",
	TokenKind_Keyword_TABLE:       "table",
	TokenKind_Keyword_AS:          "as",
	TokenKind_Keyword_IF:          "if",
	TokenKind_Keyword_NOT:         "not",
	TokenKind_Keyword_EXISTS:      "exists",
	TokenKind_Keyword_NULL:        "null",
	TokenKind_Keyword_PRIMARY:     "primary",
	TokenKind_Keyword_FOREIGN:     "foreign",
	TokenKind_Keyword_KEY:         "key",
	TokenKind_Keyword_PRIMARYKEY:  "primary key",
	TokenKind_Keyword_IFNOTEXISTS: "if not exists",
	TokenKind_Keyword_NOTNULL:     "not null",
	TokenKind_Keyword_ASC:         "asc",
	TokenKind_Keyword_DESC:        "desc",
	TokenKind_Keyword_ON:          "on",
	TokenKind_Keyword_CONFLICT:    "conflict",
	TokenKind_Keyword_ROLLBACK:    "rollback",
	TokenKind_Keyword_ABORT:       "abort",
	TokenKind_Keyword_FAIL:        "fail",
	TokenKind_Keyword_IGNORE:      "ignore",
	TokenKind_Keyword_REPLACE:     "replace",
	TokenKind_Keyword_EXPLAIN:     "explain",
	TokenKind_Keyword_QUERY:       "query",
	TokenKind_Keyword_PLAN:        "plan",
}

func (k TokenKind) String() string {
	if val, ok := tokenKindString[k]; ok {
		return val
	}
	return tokenKindString[-1]
}

const (
	TokenKindOffset_ASCII    TokenKind = 0
	TokenKindOffset_Atoms    TokenKind = 257
	TokenKindOffset_Keywords TokenKind = 400
)

const (
	TokenKind_Error     TokenKind = -1
	TokenKind_EOF       TokenKind = TokenKindOffset_ASCII
	TokenKind_LParen    TokenKind = '('
	TokenKind_RParen    TokenKind = ')'
	TokenKind_Comma     TokenKind = ','
	TokenKind_Period    TokenKind = '.'
	TokenKind_SemiColon TokenKind = ';'
	TokenKind_Plus      TokenKind = '+'
	TokenKind_Minus     TokenKind = '-'
	TokenKind_Backtic   TokenKind = '`'
)

const (
	TokenKind_Identifier TokenKind = iota + 1 + TokenKindOffset_Atoms
	TokenKind_NumericLiteral
	TokenKind_HexNumericLiteral
	TokenKind_Comment
)

const (
	TokenKind_Keyword_CREATE TokenKind = iota + 1 + TokenKindOffset_Keywords
	TokenKind_Keyword_EXPLAIN
	TokenKind_Keyword_QUERY
	TokenKind_Keyword_PLAN
	TokenKind_Keyword_TEMPORARY
	TokenKind_Keyword_VIRTUAL
	TokenKind_Keyword_TABLE
	TokenKind_Keyword_AS
	TokenKind_Keyword_IF
	TokenKind_Keyword_NOT
	TokenKind_Keyword_EXISTS
	TokenKind_Keyword_NULL

	TokenKind_Keyword_IFNOTEXISTS
	TokenKind_Keyword_NOTNULL

	TokenKind_Keyword_CONSTRAINT
	TokenKind_Keyword_PRIMARY
	TokenKind_Keyword_FOREIGN
	TokenKind_Keyword_KEY
	TokenKind_Keyword_PRIMARYKEY
	TokenKind_Keyword_UNIQUE
	TokenKind_Keyword_CHECK
	TokenKind_Keyword_DEFAULT
	TokenKind_Keyword_COLLATE
	TokenKind_Keyword_REFERENCES
	TokenKind_Keyword_GENERATED

	TokenKind_Keyword_ASC
	TokenKind_Keyword_DESC
	TokenKind_Keyword_ON
	TokenKind_Keyword_CONFLICT

	TokenKind_Keyword_ROLLBACK
	TokenKind_Keyword_ABORT
	TokenKind_Keyword_FAIL
	TokenKind_Keyword_IGNORE
	TokenKind_Keyword_REPLACE
)

const (
	Keyword_CREATE     string = "create"
	Keyword_TEMP       string = "temp"
	Keyword_TEMPORARY  string = "temporary"
	Keyword_TABLE      string = "table"
	Keyword_AS         string = "as"
	Keyword_IF         string = "if"
	Keyword_NOT        string = "not"
	Keyword_EXISTS     string = "exists"
	Keyword_NULL       string = "null"
	Keyword_CONSTRAINT string = "constraint"
	Keyword_PRIMARY    string = "primary"
	Keyword_FOREIGN    string = "foreign"
	Keyword_KEY        string = "key"
	Keyword_UNIQUE     string = "unique"
	Keyword_CHECK      string = "check"
	Keyword_DEFAULT    string = "default"
	Keyword_COLLATE    string = "collate"
	Keyword_REFERENCES string = "references"
	Keyword_GENERATED  string = "generated"
	Keyword_ASC        string = "asc"
	Keyword_DESC       string = "desc"
	Keyword_ON         string = "on"
	Keyword_CONFLICT   string = "conflict"
	Keyword_ROLLBACK   string = "rollback"
	Keyword_ABORT      string = "abort"
	Keyword_FAIL       string = "fail"
	Keyword_IGNORE     string = "ignore"
	Keyword_REPLACE    string = "replace"
	Keyword_EXPLAIN    string = "explain"
	Keyword_QUERY      string = "query"
	Keyword_PLAN       string = "plan"
)

var keywordIndex = NewIndex[string, TokenKind]().
	Add(Keyword_CREATE, TokenKind_Keyword_CREATE).
	Add(Keyword_TEMP, TokenKind_Keyword_TEMPORARY).
	Add(Keyword_TEMPORARY, TokenKind_Keyword_TEMPORARY).
	Add(Keyword_TABLE, TokenKind_Keyword_TABLE).
	Add(Keyword_AS, TokenKind_Keyword_AS).
	Add(Keyword_IF, TokenKind_Keyword_IF).
	Add(Keyword_NOT, TokenKind_Keyword_NOT).
	Add(Keyword_EXISTS, TokenKind_Keyword_EXISTS).
	Add(Keyword_NULL, TokenKind_Keyword_NULL).
	Add(Keyword_CONSTRAINT, TokenKind_Keyword_CONSTRAINT).
	Add(Keyword_PRIMARY, TokenKind_Keyword_PRIMARY).
	Add(Keyword_FOREIGN, TokenKind_Keyword_FOREIGN).
	Add(Keyword_KEY, TokenKind_Keyword_KEY).
	Add(Keyword_UNIQUE, TokenKind_Keyword_UNIQUE).
	Add(Keyword_CHECK, TokenKind_Keyword_CHECK).
	Add(Keyword_DEFAULT, TokenKind_Keyword_DEFAULT).
	Add(Keyword_COLLATE, TokenKind_Keyword_COLLATE).
	Add(Keyword_REFERENCES, TokenKind_Keyword_REFERENCES).
	Add(Keyword_GENERATED, TokenKind_Keyword_GENERATED).
	Add(Keyword_ASC, TokenKind_Keyword_ASC).
	Add(Keyword_DESC, TokenKind_Keyword_DESC).
	Add(Keyword_ON, TokenKind_Keyword_ON).
	Add(Keyword_CONFLICT, TokenKind_Keyword_CONFLICT).
	Add(Keyword_ROLLBACK, TokenKind_Keyword_ROLLBACK).
	Add(Keyword_ABORT, TokenKind_Keyword_ABORT).
	Add(Keyword_FAIL, TokenKind_Keyword_FAIL).
	Add(Keyword_IGNORE, TokenKind_Keyword_IGNORE).
	Add(Keyword_REPLACE, TokenKind_Keyword_REPLACE)

var constaintKeywords = map[TokenKind]bool{
	TokenKind_Keyword_CONSTRAINT: true,
	TokenKind_Keyword_PRIMARY:    true,
	TokenKind_Keyword_FOREIGN:    true,
	TokenKind_Keyword_UNIQUE:     true,
	TokenKind_Keyword_CHECK:      true,
	TokenKind_Keyword_DEFAULT:    true,
	TokenKind_Keyword_COLLATE:    true,
	TokenKind_Keyword_REFERENCES: true,
	TokenKind_Keyword_GENERATED:  true,
}

type TextRange struct {
	Start int
	End   int
}

type Token struct {
	Text string
	Kind TokenKind

	Pos     TextRange
	FileLoc *Location
}

func (t Token) String() string {
	if t.FileLoc != nil {
		return fmt.Sprintf("%s:%d:%d { kind: %s, text: '%s' }", t.FileLoc.FileName, t.FileLoc.Line, t.FileLoc.Col, t.Kind.String(), t.Text)
	}
	return fmt.Sprintf("{ kind: %s, text: '%s' }", t.Kind.String(), t.Text)
}
