package parser

import (
	"analizador-sintactico-proj-final/internal/lexer"
	"fmt"
	"strings"
)

type SyntaxError struct {
	Line    int
	Message string
}

type Parser struct {
	tokens []lexer.Token
	pos    int
	errors []SyntaxError
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
		errors: []SyntaxError{},
	}
}

func (p *Parser) Parse() []SyntaxError {
	// Esta función revisa si el programa tiene una forma correcta.
	if len(p.tokens) == 0 {
		return p.errors
	}

	if p.current().Value == "var" {
		p.parseDeclarations()
	}

	p.expect("begin")
	for p.current() != nil && p.current().Value != "end" {
		p.parseStatement()
	}
	p.expect("end")

	return p.errors
}

func (p *Parser) parseDeclarations() {
	// Aquí se revisan las variables declaradas al inicio del programa.
	p.consume("var")
	for p.current() != nil && p.current().Type == lexer.IDENTIFIER {
		p.consumeID()
		p.expect(":")
		p.expectType()
		p.expect(";")
	}
}

func (p *Parser) parseStatement() {
	// Cada instrucción del programa se revisa por separado.
	tok := p.current()
	if tok == nil {
		return
	}

	switch strings.ToLower(tok.Value) {
	case "if":
		p.parseIf()
	case "while":
		p.parseWhile()
	default:
		if tok.Type == lexer.IDENTIFIER {
			p.parseAssignment()
		} else {
			p.errorf("Instrucción o palabra clave inesperada: \"%s\"", tok.Value)
			p.pos++
		}
	}
}

func (p *Parser) parseIf() {
	// Se revisa la estructura de un if y sus ramas.
	p.consume("if")
	p.parseExpression()
	p.consume("then")
	p.parseStatement()
	if p.current() != nil && strings.EqualFold(p.current().Value, "else") {
		p.consume("else")
		p.parseStatement()
	}
}

func (p *Parser) parseWhile() {
	// Se revisa la estructura de un while.
	p.consume("while")
	p.parseExpression()
	p.consume("do")
	p.parseStatement()
}

func (p *Parser) parseAssignment() {
	// Se revisa que una asignación tenga el formato correcto.
	p.consumeID()
	p.expect(":=")
	p.parseExpression()
	p.expect(";")
}

func (p *Parser) parseExpression() {
	// Aquí se arma una expresión con sus partes pequeñas.
	p.parseTerm()
	for p.current() != nil && strings.ContainsRune("+-=<>!", rune(p.current().Value[0])) {
		p.advance()
		p.parseTerm()
	}
}

func (p *Parser) parseTerm() {
	// Cada parte de la expresión se revisa una por una.
	tok := p.current()
	if tok == nil {
		p.errorf("Expresión incompleta")
		return
	}

	switch tok.Type {
	case lexer.IDENTIFIER, lexer.INTEGER, lexer.REAL, lexer.STRING:
		p.advance()
	default:
		if tok.Type == lexer.OPERATOR && tok.Value == "(" {
			p.consume("(")
			p.parseExpression()
			p.expect(")")
		} else {
			p.errorf("Elemento inválido en expresión: \"%s\"", tok.Value)
			p.pos++
		}
	}
}

// --- Métodos auxiliares ---

// Avanza al siguiente token sin validación
func (p *Parser) advance() {
	if p.pos >= len(p.tokens) {
		p.errorf("Se esperaba fin de archivo estructurado")
		return
	}
	p.pos++
}

// Consume el token actual solo si coincide con el valor esperado (insensible a mayúsculas)
func (p *Parser) consume(expected string) {
	tok := p.current()
	if tok == nil {
		p.errorf("Se esperaba '%s' pero se alcanzó el fin de entrada", expected)
		return
	}
	// CORRECCIÓN SA6005: Usar strings.EqualFold en lugar de ToLower+comparación
	if !strings.EqualFold(tok.Value, expected) {
		p.errorf("Se esperaba \"%s\" pero se encontró \"%s\"", expected, tok.Value)
		return
	}
	p.pos++
}

// Consume el token actual solo si es un identificador
func (p *Parser) consumeID() {
	tok := p.current()
	if tok == nil || tok.Type != lexer.IDENTIFIER {
		p.errorf("Se esperaba un identificador")
		return
	}
	p.pos++
}

// Espera que el token actual sea un tipo de dato válido
func (p *Parser) expectType() {
	tok := p.current()
	if tok == nil {
		p.errorf("Se esperaba un tipo de dato")
		return
	}
	val := tok.Value // No necesitamos ToLower aquí porque validamos contra literales en minúsculas
	if val != "integer" && val != "real" && val != "string" && val != "boolean" {
		p.errorf("Tipo de dato no válido: %s", tok.Value)
	}
	p.pos++
}

// Espera que el token actual coincida con el valor esperado (wrapper semántico)
func (p *Parser) expect(expected string) {
	p.consume(expected)
}

// Obtiene el token actual sin avanzar
func (p *Parser) current() *lexer.Token {
	if p.pos >= len(p.tokens) {
		return nil
	}
	return &p.tokens[p.pos]
}

// Registra un error de sintaxis
func (p *Parser) errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	line := 1
	if p.current() != nil {
		line = p.current().Line
	} else if len(p.tokens) > 0 {
		line = p.tokens[len(p.tokens)-1].Line
	}
	p.errors = append(p.errors, SyntaxError{Line: line, Message: msg})
}
