package main

import (
	"reflect"
	"strings"
)

type SqlBuilder struct {
	strings.Builder
}

type SqlGeneratorConfig struct {
	EscapeIdentifiers     bool
	EscapeIdentifiersWith rune
}

var DefaultSqlGeneratorConfig = &SqlGeneratorConfig{
	EscapeIdentifiers:     true,
	EscapeIdentifiersWith: '`',
}

type SqlGenerator struct {
	config *SqlGeneratorConfig
	ast    AstNode
	out    SqlBuilder
}

func (b *SqlBuilder) WriteSpace() {
	b.WriteRune(' ')
}

func (b *SqlBuilder) WriteKeyword(keyword string) {
	b.WriteString(strings.ToUpper(keyword))
	b.WriteSpace()
}

func NewSqlGenerator(root AstNode, out strings.Builder, config *SqlGeneratorConfig) *SqlGenerator {
	if config == nil {
		config = DefaultSqlGeneratorConfig
	}
	return &SqlGenerator{
		config: config,
		ast:    root,
		out:    SqlBuilder{out},
	}
}

func (gen *SqlGenerator) Generate() string {
	gen.generateSql(gen.ast)
	return gen.out.String()
}

func (gen *SqlGenerator) generateSql(node AstNode) {

	if reflect.ValueOf(node).Kind() == reflect.Ptr && reflect.ValueOf(node).IsNil() {
		return
	}

	switch node := node.(type) {
	case AstNodeList:
		for _, n := range node {
			gen.generateSql(n)
		}
	case *AstNode_Statements:
		for _, statement := range node.Statements {
			gen.generateSql(statement)
			gen.out.Write([]byte{'\n', '\n'})
		}
	case *AstNode_Statement:
		gen.generateSql(node.Explain)
		gen.generateSql(node.Query)
		gen.generateSql(node.Plan)
		gen.generateSql(node.Statement)
		gen.out.WriteByte(';')
	case *AstNode_CreateTableStmt:
		gen.generateSql(node.Create)
		if node.Temporary != nil {
			gen.generateSql(node.Temporary)
		}
		gen.generateSql(node.Table)
		if node.IfNotExists != nil {
			gen.generateSql(node.IfNotExists)
		}
		gen.generateSql(node.TableIdentifier)
		gen.generateSql(node.TableDefinition)
	case *AstNode_TableDefinition:
		gen.out.WriteByte('(')
		gen.generateSql(node.ColumnDefinitions)
		if node.TableConstraints != nil {
			gen.out.WriteByte(',')
			gen.generateSql(node.TableConstraints)
		}
		gen.out.WriteByte('\n')
		gen.out.WriteByte(')')
	case *AstNode_ColumnDefinitions:
		for i, definition := range node.Definitions {
			gen.out.WriteByte('\n')
			gen.out.WriteByte('\t')
			gen.generateSql(definition)
			if i < len(node.Definitions)-1 {
				gen.out.WriteByte(',')
			}
		}
	case *AstNode_ColumnDefinition:
		gen.generateSql(node.ColumnName)
		gen.out.WriteSpace()
		gen.generateSql(node.TypeName)
		gen.out.WriteSpace()
		gen.generateSql(node.ColumnConstraints)
	case *AstNode_ColumnConstraints:
		for _, constraint := range node.Constraints {
			gen.generateSql(constraint)
		}
	case *AstNode_ColumnConstraint:
		if node.Name != nil {
			gen.out.WriteKeyword(Keyword_CONSTRAINT)
			gen.generateSql(node.Name)
		}
		gen.generateSql(node.Constraint)
	case *AstNode_Constraint_NotNull:
		gen.generateSql(node.Not)
		gen.generateSql(node.Null)
	case *AstNode_TableConstraints:
		for i, constraint := range node.Constraints {
			gen.out.WriteByte('\n')
			gen.out.WriteByte('\t')
			gen.generateSql(constraint)
			if i < len(node.Constraints)-1 {
				gen.out.WriteByte(',')
			}
		}
	case *AstNode_TableConstraint:
		if node.Name != nil {
			gen.out.WriteKeyword(Keyword_CONSTRAINT)
			gen.generateSql(node.Name)
		}
		gen.generateSql(node.Constraint)
	case *AstNode_TableConstraint_PrimaryKey:
		gen.out.WriteKeyword(Keyword_PRIMARY)
		gen.out.WriteKeyword(Keyword_KEY)
		gen.out.WriteByte('(')
		for i, column := range node.IndexedColumns {
			gen.generateSql(column)
			if i < len(node.IndexedColumns)-1 {
				gen.out.WriteByte(',')
			}
		}
		gen.out.WriteByte(')')
		if node.ConflictClause != nil {
			gen.generateSql(node.ConflictClause)
		}
	case *AstNode_ConflictClause:
		gen.out.WriteKeyword(Keyword_ON)
		gen.out.WriteKeyword(Keyword_CONFLICT)
		gen.generateSql(node.OnConflict)
	case *AstNode_TypeName:
		gen.out.WriteString(node.TypeName.Text)
	case *AstNode_TableIdentifer:
		if node.SchemaName != nil {
			gen.generateSql(node.SchemaName)
			gen.out.WriteByte('.')
		}
		gen.generateSql(node.TableName)
	case *AstNode_Identifier:
		if gen.config.EscapeIdentifiers {
			gen.out.WriteRune(gen.config.EscapeIdentifiersWith)
			gen.out.WriteString(node.Identifier.Text)
			gen.out.WriteRune(gen.config.EscapeIdentifiersWith)
		} else {
			gen.out.WriteString(node.Identifier.Text)
		}
	case AstNode_KeywordList:
		for _, n := range node {
			gen.generateSql(AstNode(&n))
		}
	case *AstNode_Keyword:
		text, _ := keywordIndex.GetKey(node.Keyword.Kind)
		gen.out.WriteKeyword(text)
	default:
		return
	}
}
