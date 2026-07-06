package semantic

import (
	"analizador-sintactico-proj-final/internal/lexer"
	"fmt"
	"strings"
)

type SemanticError struct {
	Line    int
	Message string
}

type Symbol struct {
	Name string
	Type string // "integer", "real", "string", "boolean"
	Line int
}

type SemanticAnalyzer struct {
	symbols map[string]Symbol
	errors  []SemanticError
}

func New() *SemanticAnalyzer {
	return &SemanticAnalyzer{
		symbols: make(map[string]Symbol),
		errors:  []SemanticError{},
	}
}

func (s *SemanticAnalyzer) Analyze(tokens []lexer.Token) (map[string]Symbol, []SemanticError) {
	s.symbols = make(map[string]Symbol)
	s.errors = []SemanticError{}

	// Filtrar comentarios y pasar a fase 1: construir tabla de símbolos
	cleanTokens := s.filterComments(tokens)
	s.buildSymbolTable(cleanTokens)

	// Fase 2: validar asignaciones y expresiones
	s.checkAssignments(cleanTokens)

	return s.symbols, s.errors
}

func (s *SemanticAnalyzer) filterComments(tokens []lexer.Token) []lexer.Token {
	result := []lexer.Token{}
	for _, t := range tokens {
		if t.Type != lexer.COMMENT {
			result = append(result, t)
		}
	}
	return result
}

func (s *SemanticAnalyzer) buildSymbolTable(tokens []lexer.Token) {
	i := 0
	for i < len(tokens) {
		if i < len(tokens) && strings.EqualFold(tokens[i].Value, "var") {
			i++ // Pasar 'var'
			for i < len(tokens) && tokens[i].Type == lexer.IDENTIFIER {
				varTok := tokens[i]
				i++ // Nombre de variable
				if i < len(tokens) && tokens[i].Value == ":" {
					i++ // Dos puntos
					if i < len(tokens) {
						typeTok := tokens[i]
						i++ // Tipo
						s.symbols[varTok.Value] = Symbol{
							Name: varTok.Value,
							Type: strings.ToLower(typeTok.Value),
							Line: varTok.Line,
						}
					}
				}
				if i < len(tokens) && tokens[i].Value == ";" {
					i++ // Punto y coma
				}
			}
		} else {
			i++
		}
	}
}

func (s *SemanticAnalyzer) checkAssignments(tokens []lexer.Token) {
	i := 0
	for i < len(tokens) {
		if tokens[i].Type == lexer.IDENTIFIER {
			varName := tokens[i].Value
			// Verificar si está declarada
			if _, exists := s.symbols[varName]; !exists {
				s.errors = append(s.errors, SemanticError{
					Line:    tokens[i].Line,
					Message: fmt.Sprintf("Variable \"%s\" no declarada", varName),
				})
			}

			// Verificar si es una asignación
			if i+1 < len(tokens) && tokens[i+1].Type == lexer.ASSIGN {
				targetType := ""
				if sym, ok := s.symbols[varName]; ok {
					targetType = sym.Type
				}
				if targetType == "" { // No declarada
					i += 2 // Saltar := y continuar
					continue
				}

				// Recoger tokens del lado derecho hasta ;
				j := i + 2
				rhsTokens := []lexer.Token{}
				for j < len(tokens) && tokens[j].Value != ";" {
					rhsTokens = append(rhsTokens, tokens[j])
					j++
				}

				// Validar tipos en la expresión
				s.checkExpressionTypes(targetType, rhsTokens, tokens[i].Line)
				i = j // Posicionar en ;
			} else {
				i++
			}
		} else {
			i++
		}
	}
}

func (s *SemanticAnalyzer) checkExpressionTypes(targetType string, expr []lexer.Token, line int) {
	for _, tok := range expr {
		var currentType string
		switch tok.Type {
		case lexer.IDENTIFIER:
			if sym, ok := s.symbols[tok.Value]; ok {
				currentType = sym.Type
			} else {
				continue // Ya se reportó como no declarada
			}
		case lexer.INTEGER:
			currentType = "integer"
		case lexer.REAL:
			currentType = "real"
		case lexer.STRING:
			currentType = "string"
		default:
			continue // Operadores, paréntesis, etc. no afectan tipo directo
		}

		// Reglas de compatibilidad (igual que tu Python)
		switch targetType {
		case "integer":
			if currentType == "string" || currentType == "real" {
				s.errors = append(s.errors, SemanticError{
					Line:    line,
					Message: fmt.Sprintf("Incompatibilidad de tipos: No se puede asignar %s a tipo entero", currentType),
				})
			}
		case "real":
			if currentType == "string" {
				s.errors = append(s.errors, SemanticError{
					Line:    line,
					Message: "Incompatibilidad de tipos: No se puede operar string con números reales",
				})
			}
		case "string":
			if currentType == "integer" || currentType == "real" {
				s.errors = append(s.errors, SemanticError{
					Line:    line,
					Message: fmt.Sprintf("Incompatibilidad de tipos: No se puede asignar tipos numéricos a cadenas (%s)", currentType),
				})
			}
		}
	}
}
