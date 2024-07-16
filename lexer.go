package main

import (
	"io"
	"os"
	"strings"
	"unicode"
)

type Location struct {
	FileName string
	Line     int
	Col      int
}

type Lexer struct {
	FileName   string
	Text       []rune
	CurrentPos int
}

func NewLexer(text string) *Lexer {
	return &Lexer{
		Text:       []rune(text),
		CurrentPos: 0,
	}
}

func NewLexerFromFile(file *os.File) (*Lexer, error) {
	text, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	lex := NewLexer(string(text))
	lex.FileName = file.Name()

	return lex, nil
}

func (l Lexer) current() rune {
	return l.Text[l.CurrentPos]
}

func (l Lexer) peek() (rune, error) {
	l.CurrentPos += 1
	if l.Eof() {
		return '\uFFFF', io.EOF
	}
	return l.current(), nil
}

func (l *Lexer) consume() {
	l.CurrentPos += 1
}

func (l *Lexer) backup() {
	l.CurrentPos -= 1
}

type ComparableFn func(r rune) bool

func (l *Lexer) consumeWhile(fn ComparableFn) {
	for !l.Eof() && fn(l.current()) {
		l.consume()
	}
}

func (l Lexer) Eof() bool {
	return l.CurrentPos == len(l.Text)
}

func (l Lexer) tokenFromCurrent(kind TokenKind) Token {
	tok := Token{
		Kind: kind,
		Text: string(l.current()),
		Pos: TextRange{
			Start: l.CurrentPos,
			End:   l.CurrentPos + 1,
		},
	}
	return tok
}

func (l Lexer) tokenFromRange(kind TokenKind, start int, end int) Token {
	tok := Token{
		Kind: kind,
		Text: string(l.Text[start:end]),
		Pos: TextRange{
			Start: start,
			End:   end,
		},
	}
	return tok
}

func (l *Lexer) SkipWhitespace() {
	for !l.Eof() && unicode.IsSpace(l.current()) {
		l.consume()
	}
}

func (l Lexer) PeekToken(count int) Token {
	res := Token{}
	for i := 0; i < count; i++ {
		res = l.ConsumeToken()
	}
	return res
}

func (l *Lexer) ConsumeToken() Token {
	tok := Token{}
	if l.Eof() {
		return tok
	}

	l.SkipWhitespace()

	start := l.CurrentPos
	current := l.current()
	switch current {
	case '/':
		l.consume()
		if l.current() == '*' {
			l.consume()
			l.consumeWhile(func(r rune) bool {
				return r != '*'
			})
			l.consume()
			if l.current() == '/' {
				l.consume()
				end := l.CurrentPos
				tok = l.tokenFromRange(TokenKind_Comment, start, end)
			}
		}
		return tok
	case '0':
		peeked, _ := l.peek()
		if unicode.ToLower(peeked) == 'x' {
			tok = l.HexNumeric()
			return tok
		}
		fallthrough
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		tok = l.Numeric()
		return tok
	case '.':
		peeked, _ := l.peek()
		if unicode.IsDigit(peeked) {
			tok = l.Numeric()
			return tok
		}
		tok = l.tokenFromCurrent(TokenKind_Period)
		l.consume()
		return tok
	case '-':
		peeked, _ := l.peek()
		if peeked == '-' {
			l.consumeWhile(func(r rune) bool { return r != '\n' })
			end := l.CurrentPos
			tok = l.tokenFromRange(TokenKind_Comment, start, end)
			return tok
		} else if peeked == '.' {
			l.consume()
			peeked, _ = l.peek()
			if unicode.IsDigit(peeked) {
				l.consume()
				numeric := l.Numeric()
				tok = l.tokenFromRange(TokenKind_NumericLiteral, start, numeric.Pos.End)
				return tok
			} else {
				l.backup()
			}
		} else if unicode.IsDigit(peeked) {
			l.consume()
			numeric := l.Numeric()
			tok = l.tokenFromRange(TokenKind_NumericLiteral, start, numeric.Pos.End)
			return tok
		}
		tok = l.tokenFromCurrent(TokenKind_Minus)
		l.consume()
		return tok
	case '`':
		l.consume()
		start = l.CurrentPos
		l.consumeWhile(func(r rune) bool { return r != '`' })
		end := l.CurrentPos
		l.consume()
		tok = l.tokenFromRange(TokenKind_Identifier, start, end)
		return tok
	case '+':
		tok = l.tokenFromCurrent(TokenKind_Plus)
		l.consume()
		return tok
	case '(':
		tok = l.tokenFromCurrent(TokenKind_LParen)
		l.consume()
		return tok
	case ')':
		tok = l.tokenFromCurrent(TokenKind_RParen)
		l.consume()
		return tok
	case ',':
		tok = l.tokenFromCurrent(TokenKind_Comma)
		l.consume()
		return tok
	case ';':
		tok = l.tokenFromCurrent(TokenKind_SemiColon)
		l.consume()
		return tok
	}

	for !l.Eof() && unicode.IsLetter(l.current()) {
		l.consume()
	}
	end := l.CurrentPos

	tok = l.tokenFromRange(TokenKind_Identifier, start, end)

	keywordmatch := strings.ToLower(tok.Text)
	if kind, ok := keywordIndex.GetValue(keywordmatch); ok {
		tok.Kind = kind
		return tok
	}

	return tok
}

func (l *Lexer) HexNumeric() Token {
	current := l.current()
	peeked, err := l.peek()
	if err != nil {
		return Token{}
	}

	if current != '0' || unicode.ToLower(peeked) != 'x' {
		return Token{}
	}

	start := l.CurrentPos
	l.consume()
	l.consume()

	for !l.Eof() && unicode.Is(unicode.ASCII_Hex_Digit, l.current()) {
		l.consume()
	}
	end := l.CurrentPos

	return l.tokenFromRange(TokenKind_HexNumericLiteral, start, end)
}

func (l *Lexer) Numeric() Token {
	tok := Token{}
	start := l.CurrentPos
	for !l.Eof() && unicode.IsDigit(l.current()) {
		l.consume()
	}
	if l.current() == '.' {
		l.consume()
		tok.Kind = TokenKind_NumericLiteral
		for !l.Eof() && unicode.IsDigit(l.current()) {
			l.consume()
		}
	}
	end := l.CurrentPos
	return l.tokenFromRange(TokenKind_NumericLiteral, start, end)
}

func (l Lexer) GetLocation(token Token) *Location {
	newlinepos := 0
	line := 1

	for i := 0; i < token.Pos.Start; i++ {
		if l.Text[i] == '\n' {
			line++
			newlinepos = i
		}
	}

	col := token.Pos.Start - newlinepos
	if newlinepos < 1 {
		col++
	}

	return &Location{
		FileName: l.FileName,
		Line:     line,
		Col:      col,
	}
}
