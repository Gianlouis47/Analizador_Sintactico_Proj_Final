package lexer

// TokenType representa la categoría de un token (un número entero)
type TokenType int

// Definición de todas las categorías de tokens posibles en nuestro lenguaje
const (
	IDENTIFIER  TokenType = iota // 0: Nombres de variables (ej: edad, x, suma)
	ASSIGN                       // 1: El símbolo de asignación (:=)
	OPERATOR                     // 2: Operadores matemáticos o lógicos (+, -, *, /, <, >, =)
	INTEGER                      // 3: Números enteros (ej: 10, 500)
	REAL                         // 4: Números con decimales (ej: 3.14, 10.5)
	STRING                       // 5: Texto entre comillas (ej: "Hola")
	KEYWORD                      // 6: Palabras reservadas (var, begin, end, if, while...)
	PUNCTUATION                  // 7: Símbolos de puntuación (;, :, (, ))
	COMMENT                      // 8: Comentarios entre { ... }
	ERROR                        // 9: Cualquier carácter que no sea reconocido
)

// Token es la estructura que representa una "palabra" del código fuente
type Token struct {
	Value string    // El texto real que escribió el usuario (ej: "begin")
	Type  TokenType // La categoría a la que pertenece (ej: KEYWORD)
	Line  int       // La línea exacta donde se encuentra el token
}
