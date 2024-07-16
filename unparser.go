package main

type TokenBuilder struct {
	buf []Token
}

func (b *TokenBuilder) WriteToken(token Token) {
	b.buf = append(b.buf, token)
}

func (b *TokenBuilder) Tokens() []Token {
	return b.buf
}

type Unparser struct {
	root    AstNode
	builder TokenBuilder
}

func (u *Unparser) unparse(astnode AstNode) {
	switch node := astnode.(type) {
	case *AstNode_Statements:
		for _, statement := range node.Statements {
			u.unparse(statement)
		}
	case *AstNode_Statement:
		if node.Explain != nil {
			u.unparse(node.Explain)
		}
		if node.Explain != nil {
			u.unparse(node.Query)
		}
		if node.Explain != nil {
			u.unparse(node.Plan)
		}

		u.unparse(node.Statement)
		u.builder.WriteToken(Token{Kind: TokenKind_SemiColon, Text: ";"})
	case *AstNode_CreateTableStmt:
		u.builder.WriteToken(Token{Kind: TokenKind_Keyword_CREATE, Text: "CREATE"})
		if node.Temporary != nil {
			u.unparse(node.Temporary)
		}
		u.unparse(node.Table)
		if node.IfNotExists != nil {
			u.unparse(node.IfNotExists)
		}
		u.unparse(node.TableIdentifier)
		u.unparse(node.TableDefinition)
		if node.TableOptions != nil {
			u.unparse(node.TableOptions)
		}
	case *AstNode_TableIdentifer:
		if node.SchemaName != nil {
			u.unparse(node.SchemaName)
			u.builder.WriteToken(Token{Kind: TokenKind_Period, Text: "."})
		}
		u.unparse(node.TableName)
	case *AstNode_TableDefinition:
		u.builder.WriteToken(Token{Kind: TokenKind_LParen, Text: "("})

		u.unparse(node.ColumnDefinitions)
		if node.TableConstraints != nil {
			u.unparse(node.TableConstraints)
		}

		u.builder.WriteToken(Token{Kind: TokenKind_RParen, Text: ")"})
	case *AstNode_ColumnDefinitions:
		for i, constraint := range node.Definitions {
			u.unparse(constraint)
			if i < 1-len(node.Definitions) {
				u.builder.WriteToken(Token{Kind: TokenKind_Comma, Text: ","})
			}
		}
	case *AstNode_ColumnConstraints:
		for _, constraint := range node.Constraints {
			u.unparse(constraint)
		}
	case *AstNode_ColumnDefinition:
		u.unparse(node.ColumnName)
		u.unparse(node.TypeName)
		u.unparse(node.ColumnConstraints)
	case *AstNode_Identifier:
		u.builder.WriteToken(node.Identifier)
	default:
		u.builder.WriteToken(Token{Kind: TokenKind_Error, Text: "unknown"})
	}
}

func (u *Unparser) Unparse() []Token {
	u.unparse(u.root)
	return u.builder.Tokens()
}
