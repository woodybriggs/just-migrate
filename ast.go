package main

type AstNode interface {
	Children() []AstNode
}

type AstNodeList []AstNode

func (l AstNodeList) Children() []AstNode {
	return l
}

type AstNode_KeywordList []AstNode_Keyword

func (l AstNode_KeywordList) Children() []AstNode {
	return l.Children()
}

type AstNode_Comment struct {
	Comment Token
}

func (n *AstNode_Comment) Children() []AstNode {
	return nil
}

type AstNode_Statements struct {
	Statements []AstNode

	Comments []AstNode_Comment
}

func (n *AstNode_Statements) Children() []AstNode {
	return n.Statements
}

type AstNode_Statement struct {
	Explain *AstNode_Keyword
	Query   *AstNode_Keyword
	Plan    *AstNode_Keyword

	Statement AstNode

	Comments []AstNode_Comment
}

func (n *AstNode_Statement) Children() []AstNode {
	return []AstNode{n.Statement}
}

type AstNode_CreateTableStmt struct {
	Create      *AstNode_Keyword
	Temporary   *AstNode_Keyword
	Table       *AstNode_Keyword
	IfNotExists AstNode_KeywordList

	TableIdentifier *AstNode_TableIdentifer
	TableDefinition *AstNode_TableDefinition

	TableOptions AstNode

	Comments []AstNode_Comment
}

func (n *AstNode_CreateTableStmt) Children() []AstNode {
	return []AstNode{n.TableIdentifier, n.TableDefinition, n.TableOptions}
}

type AstNode_Keyword struct {
	Keyword Token

	Comments []AstNode_Comment
}

func (n *AstNode_Keyword) Children() []AstNode {
	return nil
}

type AstNode_TableIdentifer struct {
	SchemaName *AstNode_Identifier
	TableName  *AstNode_Identifier

	Comments []AstNode_Comment
}

func (n *AstNode_TableIdentifer) Children() []AstNode {
	return []AstNode{n.SchemaName, n.TableName}
}

type AstNode_TableDefinition struct {
	// AsSelect *AstNode_SelectStatement
	ColumnDefinitions *AstNode_ColumnDefinitions
	TableConstraints  *AstNode_TableConstraints

	Comments []AstNode_Comment
}

func (n *AstNode_TableDefinition) Children() []AstNode {
	return []AstNode{n.ColumnDefinitions, n.TableConstraints}
}

type AstNode_ColumnDefinitions struct {
	Definitions []AstNode

	Comments []AstNode_Comment
}

func (n *AstNode_ColumnDefinitions) Children() []AstNode {
	return n.Definitions
}

func (n *AstNode_ColumnDefinitions) ColumnNamed(name string) *AstNode_ColumnDefinition {
	for _, columndef := range n.Definitions {
		switch col := columndef.(type) {
		case *AstNode_ColumnDefinition:
			if col.ColumnName.Identifier.Text == name {
				return col
			}
		default:
			continue
		}
	}

	return nil
}

type AstNode_TableConstraints struct {
	Constraints []AstNode

	Comments []AstNode_Comment
}

func (n *AstNode_TableConstraints) Children() []AstNode {
	return n.Constraints
}

type AstNode_TableConstraint struct {
	Name       *AstNode_Identifier
	Constraint AstNode

	Comments []AstNode_Comment
}

func (n *AstNode_TableConstraint) Children() []AstNode {
	return []AstNode{n.Constraint}
}

type AstNode_ColumnConstraints struct {
	Constraints []AstNode

	Comments []AstNode_Comment
}

func (n *AstNode_ColumnConstraints) Children() []AstNode {
	return n.Constraints
}

type AstNode_ColumnConstraint struct {
	Name       *AstNode_Identifier
	Constraint AstNode

	Comments []AstNode_Comment
}

func (n *AstNode_ColumnConstraint) Children() []AstNode {
	return []AstNode{n.Name, n.Constraint}
}

type AstNode_ColumnDefinition struct {
	ColumnName        *AstNode_Identifier
	TypeName          *AstNode_TypeName
	ColumnConstraints *AstNode_ColumnConstraints

	Comments []AstNode_Comment
}

func (n *AstNode_ColumnDefinition) Children() []AstNode {
	return []AstNode{n.TypeName, n.ColumnConstraints}
}

type AstNode_TypeName struct {
	TypeName Token

	Comments []AstNode_Comment
}

func (n *AstNode_TypeName) Children() []AstNode {
	return nil
}

type AstNode_Identifier struct {
	Identifier Token

	Comments []AstNode_Comment
}

func (n *AstNode_Identifier) Children() []AstNode {
	return nil
}

type AstNode_TableConstraint_PrimaryKey struct {
	PrimaryKeyword *AstNode_Keyword
	KeyKeyword     *AstNode_Keyword
	IndexedColumns AstNodeList
	ConflictClause *AstNode_ConflictClause

	Comments []AstNode_Comment
}

func (n *AstNode_TableConstraint_PrimaryKey) Children() []AstNode {
	return []AstNode{
		n.IndexedColumns,
		n.ConflictClause,
	}
}

type AstNode_ColumnConstraint_PrimaryKey struct {
	PrimaryKeyword       *AstNode_Keyword
	KeyKeyword           *AstNode_Keyword
	OrderKeyword         *AstNode_Keyword
	ConflictClause       *AstNode_ConflictClause
	AutoIncrementKeyword *AstNode_Keyword

	Comments []AstNode_Comment
}

func (n *AstNode_ColumnConstraint_PrimaryKey) Children() []AstNode {
	return []AstNode{
		n.ConflictClause,
	}
}

type AstNode_ConflictClause struct {
	OnConflict *AstNode_Keyword

	Comments []AstNode_Comment
}

func (n *AstNode_ConflictClause) Children() []AstNode {
	return nil
}

type AstNode_Constraint_NotNull struct {
	Not  *AstNode_Keyword
	Null *AstNode_Keyword

	Comments []AstNode_Comment
}

func (n *AstNode_Constraint_NotNull) Children() []AstNode {
	return []AstNode{n.Not, n.Null}
}
