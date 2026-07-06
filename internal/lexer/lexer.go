package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

// SyntaxError representa un error encontrado durante el análisis léxico
type SyntaxError struct {
	Line    int
	Message string
	Context string
}

type Lexer struct {
	source      string
	tokens      []Token
	errors      []SyntaxError
	lines       []string
	currentLine int
	pos         int
	keywords    map[string]bool
}

func New() *Lexer {
	return &Lexer{
		keywords: map[string]bool{
			"var": true, "begin": true, "end": true, "if": true, "then": true,
			"else": true, "while": true, "do": true, "for": true, "to": true,
			"integer": true, "real": true, "string": true, "boolean": true,
		},
	}
}

func (l *Lexer) Tokenize(code string) ([]Token, []SyntaxError) {
	// Este bloque convierte el código en piezas pequeñas llamadas tokens.
	l.source = code
	l.lines = strings.Split(code, "\n")
	l.tokens = []Token{}
	l.errors = []SyntaxError{}
	l.currentLine = 1
	l.pos = 0

	for l.pos < len(l.source) {
		l.skipWhitespace()
		if l.pos >= len(l.source) {
			break
		}

		ch := l.source[l.pos]

		// Si encuentra un comentario, lo ignora para no confundir al analizador.
		if ch == '{' {
			startLine := l.currentLine
			l.pos++
			for l.pos < len(l.source) && l.source[l.pos] != '}' {
				if l.source[l.pos] == '\n' {
					l.currentLine++
				}
				l.pos++
			}
			if l.pos >= len(l.source) {
				l.errors = append(l.errors, SyntaxError{
					Line:    startLine,
					Message: "Comentario no cerrado",
					Context: l.getLineContext(startLine),
				})
				break
			}
			l.pos++
			l.tokens = append(l.tokens, Token{
				Value: "{}",
				Type:  COMMENT,
				Line:  startLine,
			})
			continue
		}

		// Las cadenas se guardan completas para que no se rompan.
		if ch == '"' || ch == '\'' {
			quote := byte(ch)
			startLine := l.currentLine
			l.pos++
			for l.pos < len(l.source) && l.source[l.pos] != quote {
				if l.source[l.pos] == '\n' {
					l.currentLine++
				}
				l.pos++
			}
			if l.pos >= len(l.source) {
				l.errors = append(l.errors, SyntaxError{
					Line:    startLine,
					Message: "Cadena no cerrada",
					Context: l.getLineContext(startLine),
				})
				break
			}
			l.pos++
			l.tokens = append(l.tokens, Token{
				Value: l.source[l.pos-1 : l.pos],
				Type:  STRING,
				Line:  startLine,
			})
			continue
		}

		// Aquí se reconocen símbolos como := y otros operadores dobles.
		if l.pos+1 < len(l.source) {
			twoChar := l.source[l.pos : l.pos+2]
			if twoChar == ":=" || twoChar == "<>" || twoChar == "<=" || twoChar == ">=" {
				l.tokens = append(l.tokens, Token{
					Value: twoChar,
					Type:  ASSIGN,
					Line:  l.currentLine,
				})
				if twoChar != ":=" {
					last := &l.tokens[len(l.tokens)-1]
					last.Type = OPERATOR
				}
				l.pos += 2
				continue
			}
		}

		// Operadores simples y puntuación
		if strings.ContainsRune("+-*/=<>();:.,[]{}", rune(ch)) {
			tType := OPERATOR
			if !strings.ContainsRune("+-*/<>", rune(ch)) {
				tType = PUNCTUATION
			}
			l.tokens = append(l.tokens, Token{
				Value: string(ch),
				Type:  tType,
				Line:  l.currentLine,
			})
			l.pos++
			continue
		}

		// Si ve un número, lo clasifica como entero o real.
		if unicode.IsDigit(rune(ch)) {
			start := l.pos
			dots := 0
			for l.pos < len(l.source) && (unicode.IsDigit(rune(l.source[l.pos])) || l.source[l.pos] == '.') {
				if l.source[l.pos] == '.' {
					dots++
				}
				l.pos++
			}
			num := l.source[start:l.pos]
			if dots > 1 || strings.HasSuffix(num, ".") {
				l.errors = append(l.errors, SyntaxError{
					Line:    l.currentLine,
					Message: fmt.Sprintf("Número mal formado: \"%s\"", num),
					Context: l.getLineContext(l.currentLine),
				})
				l.tokens = append(l.tokens, Token{Value: num, Type: ERROR, Line: l.currentLine})
			} else {
				tType := INTEGER
				if dots == 1 {
					tType = REAL
				}
				l.tokens = append(l.tokens, Token{Value: num, Type: tType, Line: l.currentLine})
			}
			continue
		}

		// Aquí se reconocen los nombres de variables y las palabras reservadas.
		if unicode.IsLetter(rune(ch)) || ch == '_' {
			start := l.pos
			for l.pos < len(l.source) && (unicode.IsLetter(rune(l.source[l.pos])) || unicode.IsDigit(rune(l.source[l.pos])) || l.source[l.pos] == '_') {
				l.pos++
			}
			word := l.source[start:l.pos]
			if l.keywords[strings.ToLower(word)] {
				l.tokens = append(l.tokens, Token{Value: word, Type: KEYWORD, Line: l.currentLine})
			} else {
				l.tokens = append(l.tokens, Token{Value: word, Type: IDENTIFIER, Line: l.currentLine})
			}
			continue
		}

		// Si aparece algo raro, se marca como error para que el usuario lo corrija.
		l.errors = append(l.errors, SyntaxError{
			Line:    l.currentLine,
			Message: fmt.Sprintf("Carácter inválido: \"%s\"", string(ch)),
			Context: l.getLineContext(l.currentLine),
		})
		l.tokens = append(l.tokens, Token{Value: string(ch), Type: ERROR, Line: l.currentLine})
		l.pos++
	}

	return l.tokens, l.errors
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.source) && strings.ContainsRune(" \t\r\n", rune(l.source[l.pos])) {
		if l.source[l.pos] == '\n' {
			l.currentLine++
		}
		l.pos++
	}
}

func (l *Lexer) getLineContext(line int) string {
	if line < 1 || line > len(l.lines) {
		return ""
	}
	return l.lines[line-1]
}
