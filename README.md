# Analizador Léxico, Sintáctico y Semántico (Pascal-Like)

Este es un compilador educativo en su fase inicial (Front-End) que realiza el análisis léxico, sintáctico y semántico de un lenguaje de programación simplificado inspirado en Pascal. Su objetivo es validar que el código fuente ingresado por el usuario cumpla estrictamente con las reglas gramaticales y de tipado establecidas.

## Características del Analizador

### 1. Análisis Léxico
- Escanea el código carácter por carácter en un flujo global plano.
- Soporta comentarios multilínea delimitados por `{ ... }` y cadenas de texto delimitadas por comillas simples o dobles.
- Identifica y clasifica tokens en categorías como: Palabras Clave (`var`, `begin`, `end`, `if`, `while`), Identificadores, Operadores, Enteros, Reales y Puntuación.
- Realiza validaciones estrictas sobre representaciones numéricas (evitando números mal formados con múltiples puntos decimales).

### 2. Análisis Sintáctico
- Implementa un Analizador de Descenso Recurrente (Recursive Descent Parser).
- Valida la estructura jerárquica del código, exigiendo la correcta apertura y cierre de bloques mediante las estructuras `begin` y `end`.
- Soporta el reconocimiento de estructuras de control condicionales (`if-then-else`), bucles (`while-do`) y asignaciones básicas (`:=`).

### 3. Análisis Semántico
- Administra una Tabla de Símbolos global basada en las declaraciones dentro del bloque `var`.
- Detecta el uso de variables no declaradas en cualquier sección del bloque de ejecución.
- Realiza un chequeo de tipos estricto para evitar asignaciones inválidas (por ejemplo, asignar valores reales o cadenas de texto a variables de tipo entero).

## Requisitos de Ejecución
- Go 1.19 o superior instalado en el sistema.
- No requiere dependencias externas o bibliotecas de terceros (utiliza exclusivamente módulos nativos y estructuras de datos estándar).

## Estructura del Código Fuente Ingresado
El analizador espera un programa estructurado bajo el siguiente esquema:

```pascal
var
  nombre_variable : tipo_dato;
begin
  nombre_variable := valor;
end
