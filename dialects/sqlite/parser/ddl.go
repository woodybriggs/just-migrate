package sqlite

import (
	"woodybriggs/justmigrate/core/ast"
	"woodybriggs/justmigrate/core/report"
	"woodybriggs/justmigrate/core/tik"
)

func (p *SqliteParser) CreateStatement() ast.Statement {
	p.PushParseContext("create statement")
	defer p.PopParseContext()

	switch p.Peeked().Kind {
	case tik.TokenKind_Keyword_TABLE:
		return p.CreateTableStatement(false)
	case tik.TokenKind_Keyword_VIEW:
		return p.CreateViewStatement(false)
	case tik.TokenKind_Keyword_TRIGGER:
		return p.CreateTriggerStatement(false)
	case tik.TokenKind_Keyword_INDEX:
		isUnique := false
		return p.CreateIndexStatement(isUnique)
	case tik.TokenKind_Keyword_UNIQUE:
		isUnique := true
		return p.CreateIndexStatement(isUnique)
	case tik.TokenKind_Keyword_VIRTUAL:
		return p.CreateVirtualTableStatement()
	case tik.TokenKind_Keyword_TEMPORARY:
		return p.CreateTemporaryStatement()
	default:
		err := report.
			NewReport("parse error").
			WithLabels([]report.Label{
				{
					Source: p.Current().SourceCode,
					Range:  p.Current().SourceRange,
					Note:   "unknown token for create statement",
				},
			})
		p.ReportError(err)
		return nil
	}
}

func (p *SqliteParser) CreateTableStatement(isTemporary bool) *ast.CreateTable {
	p.PushParseContext("create table statement")
	defer p.PopParseContext()

	var temporaryKeyword *ast.Keyword = nil

	createKeyword := ast.Keyword(p.Expect(tik.TokenKind_Keyword_CREATE))

	if isTemporary {
		temporaryKeyword = ast.MakeKeyword(p.Expect(tik.TokenKind_Keyword_TEMPORARY))
	}

	tableKeyword := ast.Keyword(p.Expect(tik.TokenKind_Keyword_TABLE))

	ifnotexists := p.MaybeIfNotExists()

	tableIdent := p.CatalogObjectIdentifier()

	tableDefinition := p.TableDefinition()

	tableOptions := p.TableOptions()

	return ast.MakeCreateTable(
		createKeyword,
		temporaryKeyword,
		tableKeyword,
		ifnotexists,
		tableIdent,
		tableDefinition,
		tableOptions,
	)
}

func (p *SqliteParser) CreateViewStatement(false bool) ast.Statement {
	panic("unimplemented")
}

func (p *SqliteParser) CreateTriggerStatement(false bool) ast.Statement {
	panic("unimplemented")
}

func (p *SqliteParser) CreateIndexStatement(isUnique bool) ast.Statement {
	panic("unimplemented")
}

func (p *SqliteParser) CreateVirtualTableStatement() ast.Statement {
	panic("unimplemented")
}

func (p *SqliteParser) CreateTemporaryStatement() ast.Statement {
	panic("unimplemented")
}
