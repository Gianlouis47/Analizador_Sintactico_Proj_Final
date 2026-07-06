Analizador Lexico, Sintactico y Semantico (Pascal-Like)
Este proyecto consiste en un compilador educativo en su fase inicial (Front-End) que realiza el analisis lexico, sintactico y semantico de un lenguaje de programacion simplificado inspirado en Pascal. El objetivo del programa es validar que el codigo fuente ingresado por el usuario cumpla estrictamente con las reglas gramaticales y de tipado establecidas.

Caracteristicas del Analizador
1. Analisis Lexico
Escanea el codigo caracter por caracter en un flujo global plano.
Soporta comentarios multilinea delimitados por { ... } y cadenas de texto delimitadas por comillas simples o dobles.
Identifica y clasifica tokens en categorias como: Palabras Clave (var, begin, end, if, while), Identificadores, Operadores, Enteros, Reales y Puntuacion.
Realiza validaciones estrictas sobre representaciones numericas (evitando numeros mal formados con multiples puntos decimales).
2. Analisis Sintactico
Implementa un Analizador de Descenso Recurrente (Recursive Descent Parser).
Valida la estructura jerarquica del codigo, exigiendo la correcta apertura y cierre de bloques mediante las estructuras begin y end.
Soporta el reconocimiento de estructuras de control condicionales (if-then-else), bucles (while-do) y asignaciones basicas (:=).
3. Analisis Semantico
Administra una Tabla de Simbolos global basada en las declaraciones dentro del bloque var.
Detecta el uso de variables no declaradas en cualquier seccion del bloque de ejecucion.
Realiza un chequeo de tipos estricto para evitar asignaciones invalidas (por ejemplo, asignar valores reales o cadenas de texto a variables de tipo entero).
Requisitos de Ejecucion
Python 3.6 o superior instalado en el sistema.
No requiere dependencias externas o librerias adicionales de terceros (utiliza exclusivamente módulos nativos y estructuras de datos estandard).
Estructura del Codigo Fuente Ingresado
El analizador espera un programa estructurado bajo el siguiente esquema:

var
  nombre_variable : tipo_dato;
begin
  nombre_variable := valor;
end
