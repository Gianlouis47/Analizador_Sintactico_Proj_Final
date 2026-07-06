package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"analizador-sintactico-proj-final/internal/lexer"
	"analizador-sintactico-proj-final/internal/parser"
	"analizador-sintactico-proj-final/internal/semantic"
)

func center(s string, width int) string {
	if len(s) >= width {
		return s
	}
	padding := width - len(s)
	left := padding / 2
	right := padding - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

func main() {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println(center("ANALIZADOR LÉXICO-SINTÁCTICO-SEMÁNTICO (PASCAL-LIKE)", 70))
	fmt.Println(strings.Repeat("=", 70))
	fmt.Print("\nINGRESE SU CÓDIGO (Escriba 'FIN' en una línea vacía para compilar)\n")
	fmt.Print(strings.Repeat("-", 70) + "\n")

	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	lineNum := 1
	for scanner.Scan() {
		line := scanner.Text()
		if line == "FIN" {
			break
		}
		lines = append(lines, line)
		lineNum++
	}
	code := strings.Join(lines, "\n")

	l := lexer.New()
	tokens, lexErrors := l.Tokenize(code)

	p := parser.New(tokens)
	syntaxErrors := p.Parse()

	s := semantic.New()
	symbols, semanticErrors := s.Analyze(tokens)

	printResults(lexErrors, syntaxErrors, semanticErrors, symbols)
}

func printResults(lexErrors []lexer.SyntaxError, syntaxErrors []parser.SyntaxError, semanticErrors []semantic.SemanticError, symbols map[string]semantic.Symbol) {
	fmt.Print("\n" + strings.Repeat("=", 70))
	fmt.Print("\n" + center("RESULTADO GLOBAL DE LAS FASES", 70))
	fmt.Print("\n" + strings.Repeat("=", 70))

	fmt.Printf("\n\n[Fase Léxica]: ")
	if len(lexErrors) > 0 {
		fmt.Print("ERRORES")
		for _, err := range lexErrors {
			fmt.Printf("\n  Línea %d: %s -> \"%s\"", err.Line, err.Message, err.Context)
		}
	} else {
		fmt.Print("OK")
	}

	fmt.Printf("\n\n[Fase Sintáctica]: ")
	if len(syntaxErrors) > 0 {
		fmt.Print("ERRORES")
		for _, err := range syntaxErrors {
			fmt.Printf("\n  Línea %d: %s", err.Line, err.Message)
		}
	} else {
		fmt.Print("OK")
	}

	fmt.Printf("\n\n[Fase Semántica]: ")
	if len(semanticErrors) > 0 {
		fmt.Print("ERRORES")
		for _, err := range semanticErrors {
			fmt.Printf("\n  Línea %d: %s", err.Line, err.Message)
		}
	} else {
		fmt.Print("OK")
	}

	if len(symbols) > 0 {
		fmt.Print("\n\n" + strings.Repeat("-", 45))
		fmt.Print("\n  {'Identificador':<15} {'Tipo':<12} {'Linea':<8}")
		fmt.Print("\n   " + strings.Repeat("-", 38))
		for _, sym := range symbols {
			fmt.Printf("\n  %-15s %-12s %-8d", sym.Name, sym.Type, sym.Line)
		}
	}
	fmt.Print("\n" + strings.Repeat("=", 70) + "\n")
}
