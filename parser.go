package main

import (
	"errors"
	"fmt"
	"os"
)

type ParseError struct {
	err   error
	token Token
	lexer *Lexer
}

func (p *ParseError) Join(err error) *ParseError {
	p.err = errors.Join(err, p.err)
	return p
}

func (p *ParseError) Error() string {
	loc := p.lexer.GetLocation(p.token)
	return fmt.Sprintf("%s:%d:%d\n%s", loc.FileName, loc.Line, loc.Col, p.err)
}

type Parser struct {
	lexer       *Lexer
	parseerrors []ParseError
}

func (p *Parser) NewParseError(err error, token Token) *ParseError {
	return &ParseError{
		err:   err,
		token: token,
		lexer: p.lexer,
	}
}

func (p *Parser) ReportError(err ParseError, token Token) {
	p.parseerrors = append(p.parseerrors, err)
	location := p.lexer.GetLocation(token)
	fmt.Printf(
		"%s:%d:%d\n\t%s\n",
		location.FileName,
		location.Line,
		location.Col,
		err.Error(),
	)
	if len(p.parseerrors) >= 10 {
		fmt.Println("too many errors")
		os.Exit(1)
	}
}

func (p *Parser) HasErrors() bool {
	return len(p.parseerrors) > 0
}

func (p *Parser) AcceptToken(kind TokenKind) (Token, bool, []Token) {
	var comments []Token = nil
	peeked := p.lexer.PeekToken(1)
	if peeked.Kind == TokenKind_Comment {
		comments = append(comments, peeked)
	}
	if peeked.Kind != kind {
		return Token{}, false, comments
	}
	return p.lexer.ConsumeToken(), true, comments
}

func (p *Parser) ConsumeToken() Token {
	return p.lexer.ConsumeToken()
}

func (p *Parser) MakeComments(comments []Token) []AstNode_Comment {
	if comments == nil {
		return nil
	}
	if len(comments) == 0 {
		return nil
	}
	result := make([]AstNode_Comment, len(comments))
	for i, comment := range comments {
		result[i] = AstNode_Comment{
			Comment: comment,
		}
	}
	return result
}

func (p *Parser) Statements() *AstNode_Statements {
	stmts := &AstNode_Statements{
		Statements: make([]AstNode, 0),
	}

	for !p.lexer.Eof() {
		stmt, err := p.Statement()
		if err != nil {
			fmt.Println(err)
			return nil
		}
		stmts.Statements = append(stmts.Statements, stmt)
	}

	return stmts
}

func (p *Parser) Statement() (AstNode, error) {
	var stmt *AstNode_Statement = &AstNode_Statement{}

	explain, err := p.Keyword(TokenKind_Keyword_EXPLAIN)
	if err == nil {
		stmt.Explain = explain

		query, err := p.Keyword(TokenKind_Keyword_QUERY)
		if err == nil {
			stmt.Query = query

			plan, err := p.Keyword(TokenKind_Keyword_PLAN)
			if err != nil {
				return nil, p.NewParseError(errors.New("expected 'PLAN' after 'QUERY' in Statement"), query.Keyword)
			}
			stmt.Plan = plan
		}
	}

	token := p.lexer.PeekToken(1)
	switch token.Kind {
	case TokenKind_Keyword_CREATE:
		create, err := p.Keyword(TokenKind_Keyword_CREATE)
		if err != nil {
			return nil, err
		}
		temporary, _ := p.Keyword(TokenKind_Keyword_TEMPORARY)
		catalogObj := p.ConsumeToken()
		switch catalogObj.Kind {
		case TokenKind_Keyword_TABLE:
			table := &AstNode_Keyword{Keyword: catalogObj}
			statement, err := p.CreateTableStatement(create, table, temporary)
			if err != nil {
				return nil, err
			}
			stmt.Statement = statement
		default:
			return nil, p.NewParseError(errors.New("expected catalog object in create statement"), catalogObj)
		}
	default:
		return nil, p.NewParseError(errors.New("unsupported"), token)
	}

	semi, ok, comments := p.AcceptToken(TokenKind_SemiColon)
	if !ok {
		return nil, p.NewParseError(errors.New("expected semicolon at end of statement"), semi)
	}
	if comments != nil {
		stmt.Comments = append(stmt.Comments, p.MakeComments(comments)...)
	}

	return stmt, nil
}

func (p *Parser) CreateTableStatement(create *AstNode_Keyword, temporary *AstNode_Keyword, table *AstNode_Keyword) (*AstNode_CreateTableStmt, error) {

	createTableStmt := &AstNode_CreateTableStmt{
		Create:    create,
		Table:     table,
		Temporary: temporary,
	}

	ifnotexits, err := p.IfNotExists()
	if err != nil {
		return nil, err.(*ParseError).Join(errors.New("expected if not exists"))
	}

	tableIdentifier, err := p.TableIdentifier()
	if err != nil {
		return nil, errors.Join(err, errors.New("expected table identfier"))
	}
	tableDefinition, err := p.TableDefinition()
	if err != nil {
		return nil, errors.Join(err, errors.New("expected table definition"))
	}
	// createTableStmt.TableOptions = p.TableOptions()

	createTableStmt.IfNotExists = ifnotexits
	createTableStmt.TableIdentifier = tableIdentifier
	createTableStmt.TableDefinition = tableDefinition

	return createTableStmt, nil
}

func TokensToStrings(tokens []TokenKind) []string {
	result := make([]string, len(tokens))
	for i, token := range tokens {
		result[i] = token.String()
	}
	return result
}

func (p *Parser) IfNotExists() (AstNode_KeywordList, error) {
	if_, _ := p.Keyword(TokenKind_Keyword_IF)
	if if_ == nil {
		return nil, nil
	}

	not, err := p.Keyword(TokenKind_Keyword_NOT)
	if err != nil {
		return nil, err
	}

	exists, err := p.Keyword(TokenKind_Keyword_EXISTS)
	if err != nil {
		return nil, err
	}

	return AstNode_KeywordList{*if_, *not, *exists}, nil
}

func (p *Parser) Keyword(keyword TokenKind) (*AstNode_Keyword, *ParseError) {
	kword, ok, comments := p.AcceptToken(keyword)
	if !ok {
		keywordstr, _ := keywordIndex.GetKey(keyword)
		return nil, p.NewParseError(fmt.Errorf("expected keyword %s", keywordstr), kword)
	}
	return &AstNode_Keyword{
		Keyword:  kword,
		Comments: p.MakeComments(comments),
	}, nil

}

func (p *Parser) TableIdentifier() (*AstNode_TableIdentifer, error) {

	schemaOrTable, err := p.Identifier()
	if err != nil {
		return nil, errors.Join(err)
	}

	_, ok, periodcomments := p.AcceptToken(TokenKind_Period)
	if !ok {
		return &AstNode_TableIdentifer{
			TableName: schemaOrTable,
			Comments:  p.MakeComments(periodcomments),
		}, nil
	}

	table, err := p.Identifier()
	if err != nil {
		return nil, errors.Join(err, errors.New("expected table name after schema name"))
	}

	return &AstNode_TableIdentifer{
		SchemaName: schemaOrTable,
		TableName:  table,
	}, nil
}

func (p *Parser) TableDefinition() (*AstNode_TableDefinition, error) {
	lparen, ok, _ /*lparencomments*/ := p.AcceptToken(TokenKind_LParen)
	if !ok {
		return nil, p.NewParseError(errors.New("expected '(' at start of column definitions"), lparen)
	}

	columnDefs, err := p.ColumnDefinitions()
	if err != nil {
		return nil, err
	}
	tableConstraints, err := p.TableConstraints(columnDefs)
	if err != nil {
		return nil, err
	}

	rparen, ok, _ /*rparencomments*/ := p.AcceptToken(TokenKind_RParen)
	if !ok {
		return nil, p.NewParseError(errors.New("expected ')' at end of table definition"), rparen)
	}

	return &AstNode_TableDefinition{
		ColumnDefinitions: columnDefs,
		TableConstraints:  tableConstraints,
	}, nil
}

func IsColumnDefinitionsTerminal(token Token) bool {
	if token.Kind == TokenKind_RParen {
		return true
	}
	_, ok := constaintKeywords[token.Kind]
	return ok
}

func (p *Parser) ColumnDefinitions() (*AstNode_ColumnDefinitions, error) {

	node := &AstNode_ColumnDefinitions{
		Definitions: []AstNode{},
	}

	token := p.lexer.PeekToken(1)

	for {
		if IsColumnDefinitionsTerminal(token) {
			break
		}

		columnDef, err := p.ColumnDefinition()
		if err != nil {
			return nil, err
		}
		node.Definitions = append(node.Definitions, columnDef)

		token = p.lexer.PeekToken(1)
		if token.Kind == TokenKind_Comma {
			comma := p.ConsumeToken()
			token = p.lexer.PeekToken(1)
			if token.Kind == TokenKind_RParen {
				return nil, p.NewParseError(errors.New("unexpected trailing comma"), comma)
			}
		}
	}

	return node, nil
}

func (p *Parser) TableConstraints(columnDefinitions *AstNode_ColumnDefinitions) (*AstNode_TableConstraints, error) {

	var constraints AstNodeList = AstNodeList{}
	token := p.lexer.PeekToken(1)

	for {
		if token.Kind == TokenKind_RParen {
			break
		}

		tableConstraint, err := p.TableConstraint(columnDefinitions)
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, tableConstraint)

		token = p.lexer.PeekToken(1)
		if token.Kind == TokenKind_Comma {
			comma := p.ConsumeToken()
			token = p.lexer.PeekToken(1)
			if token.Kind == TokenKind_RParen {
				return nil, p.NewParseError(errors.New("unexpected trailing comma"), comma)
			}
		}
	}

	if len(constraints) > 0 {
		return &AstNode_TableConstraints{
			Constraints: constraints,
		}, nil
	}
	return nil, nil
}

func (p *Parser) PeekedIsColumnDefTerminal() bool {
	next := p.lexer.PeekToken(1)
	switch next.Kind {
	case TokenKind_Comma:
		return true
	case TokenKind_RParen:
		return true
	default:
		return false
	}
}

func (p *Parser) ColumnDefinition() (AstNode, error) {

	columnIdentifier, err := p.Identifier()
	if err != nil {
		return nil, err
	}
	columnType, err := p.TypeName()
	if err != nil {
		return nil, err
	}
	columnConstraints, err := p.ColumnConstraints()
	if err != nil {
		return nil, err
	}

	return &AstNode_ColumnDefinition{
		ColumnName:        columnIdentifier,
		TypeName:          columnType,
		ColumnConstraints: columnConstraints,
	}, nil
}

func (p *Parser) TableConstraint(columnDefs *AstNode_ColumnDefinitions) (AstNode, error) {

	var name *AstNode_Identifier = nil

	_, err := p.Keyword(TokenKind_Keyword_CONSTRAINT)
	if err == nil {
		ident, err := p.Identifier()
		if err != nil {
			return nil, err
		}
		name = ident
	}

	token := p.lexer.PeekToken(1)

	switch token.Kind {
	case TokenKind_Keyword_PRIMARY:
		primarykey, err := p.TableConstraint_PrimaryKey(columnDefs)
		if err != nil {
			return nil, err
		}
		return &AstNode_TableConstraint{
			Name:       name,
			Constraint: primarykey,
		}, nil
	case TokenKind_Keyword_UNIQUE:
		return nil, p.NewParseError(errors.New("not implemented"), token)
	case TokenKind_Keyword_CHECK:
		return nil, p.NewParseError(errors.New("not implemented"), token)
	case TokenKind_Keyword_FOREIGN:
		return nil, p.NewParseError(errors.New("not implemented"), token)
	default:
		return nil, p.NewParseError(errors.New("expected table constraint"), token)
	}
}

func (p *Parser) TableConstraint_PrimaryKey(columnDefs *AstNode_ColumnDefinitions) (AstNode, error) {

	primary, err := p.Keyword(TokenKind_Keyword_PRIMARY)
	if err != nil {
		return nil, err
	}

	key, err := p.Keyword(TokenKind_Keyword_KEY)
	if err != nil {
		return nil, err
	}

	if lparen, ok, _ /*comments*/ := p.AcceptToken(TokenKind_LParen); !ok {
		return nil, p.NewParseError(errors.New("expected '(' after primary key table constraint"), lparen)
	}

	tableConstraint := &AstNode_TableConstraint_PrimaryKey{
		PrimaryKeyword: primary,
		KeyKeyword:     key,
	}

	parsingIndexCols := true
	for parsingIndexCols {
		indexedcol, err := p.IndexedColumn(columnDefs)
		if err != nil {
			return nil, err
		}
		tableConstraint.IndexedColumns = append(tableConstraint.IndexedColumns, indexedcol)

		if _, ok, _ /*comments*/ := p.AcceptToken(TokenKind_Comma); ok {
			continue
		}

		next := p.lexer.PeekToken(1)
		if next.Kind == TokenKind_RParen {
			parsingIndexCols = false
		}
	}

	if rparen, ok, _ /*comments*/ := p.AcceptToken(TokenKind_RParen); !ok {
		return nil, p.NewParseError(errors.New("expected ')' at end of indexed columns of primary key constraint on table"), rparen)
	}

	conflictclause, _ := p.ConflictClause()
	tableConstraint.ConflictClause = conflictclause

	return tableConstraint, nil
}

func (p *Parser) IndexedColumn(tabledef *AstNode_ColumnDefinitions) (AstNode, error) {
	token := p.lexer.PeekToken(1)

	switch token.Kind {
	case TokenKind_Identifier:
		{
			ident, _, _ /*comments*/ := p.AcceptToken(TokenKind_Identifier)
			if columndef := tabledef.ColumnNamed(ident.Text); columndef != nil {
				return columndef.ColumnName, nil
			}
			return nil, p.NewParseError(fmt.Errorf("could not find column named %s for use in indexed column", token), token)
		}
	default:
		{
			return nil, p.NewParseError(errors.New("expected column name in indexed column"), token)
		}
	}
}

func (p *Parser) ConflictClause() (*AstNode_ConflictClause, error) {

	_, err := p.Keyword(TokenKind_Keyword_ON)
	if err != nil {
		return nil, err
	}

	_, err = p.Keyword(TokenKind_Keyword_CONFLICT)
	if err != nil {
		return nil, err
	}

	conflictcommand := p.lexer.PeekToken(1)
	switch conflictcommand.Kind {
	case TokenKind_Keyword_ROLLBACK:
		rollback, err := p.Keyword(TokenKind_Keyword_ROLLBACK)
		if err != nil {
			return nil, err
		}
		return &AstNode_ConflictClause{
			OnConflict: rollback,
		}, nil
	case TokenKind_Keyword_ABORT:
		abort, err := p.Keyword(TokenKind_Keyword_ABORT)
		if err != nil {
			return nil, err
		}
		return &AstNode_ConflictClause{
			OnConflict: abort,
		}, nil
	case TokenKind_Keyword_FAIL:
		fail, err := p.Keyword(TokenKind_Keyword_FAIL)
		if err != nil {
			return nil, err
		}
		return &AstNode_ConflictClause{
			OnConflict: fail,
		}, nil
	case TokenKind_Keyword_IGNORE:
		ignore, err := p.Keyword(TokenKind_Keyword_IGNORE)
		if err != nil {
			return nil, err
		}
		return &AstNode_ConflictClause{
			OnConflict: ignore,
		}, nil
	case TokenKind_Keyword_REPLACE:
		replace, err := p.Keyword(TokenKind_Keyword_REPLACE)
		if err != nil {
			return nil, err
		}
		return &AstNode_ConflictClause{
			OnConflict: replace,
		}, nil
	default:
		return nil, p.NewParseError(errors.New("expected conflict clause verb"), conflictcommand)
	}
}

func (p *Parser) TypeName() (*AstNode_TypeName, error) {
	ident, ok, comments := p.AcceptToken(TokenKind_Identifier)
	if !ok {
		return nil, p.NewParseError(errors.New("expected type name"), ident)
	}

	return &AstNode_TypeName{
		TypeName: ident,
		Comments: p.MakeComments(comments),
	}, nil
}

func (p *Parser) ColumnConstraints() (*AstNode_ColumnConstraints, error) {

	result := &AstNode_ColumnConstraints{
		Constraints: []AstNode{},
	}

	for !p.PeekedIsColumnDefTerminal() {
		columnConstraint, err := p.ColumnConstraint()
		if err != nil {
			return nil, err
		}
		result.Constraints = append(result.Constraints, columnConstraint)
	}

	return result, nil
}

func (p *Parser) ColumnConstraint() (*AstNode_ColumnConstraint, error) {

	var err error
	var name *AstNode_Identifier = nil

	constraint, _ := p.Keyword(TokenKind_Keyword_CONSTRAINT)
	if constraint != nil {
		name, err = p.Identifier()
		if err != nil {
			return nil, err
		}
	}

	token := p.lexer.PeekToken(1)

	switch token.Kind {
	case TokenKind_Keyword_PRIMARY:

		primary, err := p.Keyword(TokenKind_Keyword_PRIMARY)
		if err != nil {
			return nil, p.NewParseError(errors.New("expected 'primary' keyword"), primary.Keyword)
		}

		key, err := p.Keyword(TokenKind_Keyword_KEY)
		if err != nil {
			return nil, err.Join(errors.New("expected 'KEY' keyword after 'PRIMARY'"))
		}

		asc, _ := p.Keyword(TokenKind_Keyword_ASC)
		desc, _ := p.Keyword(TokenKind_Keyword_DESC)

		conflictclause, _ := p.ConflictClause()

		primarykeyconstraint := &AstNode_ColumnConstraint_PrimaryKey{
			PrimaryKeyword: primary,
			KeyKeyword:     key,
			ConflictClause: conflictclause,
		}

		if asc != nil {
			primarykeyconstraint.OrderKeyword = asc
		}
		if desc != nil {
			primarykeyconstraint.OrderKeyword = desc
		}

		return &AstNode_ColumnConstraint{
			Name:       name,
			Constraint: primarykeyconstraint,
		}, nil
	case TokenKind_Keyword_NOT:
		not, err := p.Keyword(TokenKind_Keyword_NOT)
		if err != nil {
			return nil, err
		}

		null, err := p.Keyword(TokenKind_Keyword_NULL)
		if err != nil {
			return nil, err
		}
		return &AstNode_ColumnConstraint{
			Constraint: &AstNode_Constraint_NotNull{
				Not:  not,
				Null: null,
			},
		}, nil
	case TokenKind_Keyword_DEFAULT:
		// def := p.Keyword(TokenKind_Keyword_DEFAULT)
		fallthrough
	default:
		{
			return nil, p.NewParseError(errors.New("expected identifier"), token)
		}
	}
}

func (p *Parser) Identifier() (*AstNode_Identifier, error) {
	ident, ok, comments := p.AcceptToken(TokenKind_Identifier)
	if !ok {
		return nil, p.NewParseError(errors.New("expected identifier"), ident)
	}

	return &AstNode_Identifier{
		Identifier: ident,
		Comments:   p.MakeComments(comments),
	}, nil
}
