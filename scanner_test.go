package main

/*
- file does not exist
- file without reading permission
- empty file
- location information test???? - to chyba nie potrzebne
- `{` - LEFT_PARENTHESES
- `}` - RIGHT_PARENTHESES
- `(` - LEFT_BRACE
- `)` - RIGHT_BRACE
- `.` - DOT
- `,` - COMMA
- `-` - MINUS
- `+` - PLUS
- `;` - SEMICOLON
- `*` - STAR
- `=` - EQUAL
- `==` - EQUAL_EQUAL
- `<` - LESS
- `<=` - LESS_EQUAL
- `>` - GREATER
- `>=` - GREATER_EQUAL
- `!=` - BANG_EQUAL
- `/` - SLASH
- `EOF` - END OF FILE
- testy apropos newlines i whitespaces - chyba `' '`, `\r`, `\t`, `\n`
- testy stringów: `""`, `"ala ma kota"`, `"kot"`
- numerki: `123`, `123.456`, `.456` => dot & `456`, `123.` => `123` & `.`
- identifiers - `or` => `OR`, `orchid` => IDENTIFIER, `if_var` => IDENTIFIER, `while123` => IDENTIFIER, `var ala = 10` => {`VAR`, IDENTIFIER, `EQUAL`, `NUMBER`}
- keywords: `and`, `class`, `else`, `false`, `for`, `fun`, `if`, `nil`, `or`, `print`, `return`, `super`, `this`, `true`, `var`, `while`
- comments???
- podrozdział 4.5.1 - characters not used by lox, e.g. `@`, `#`, `^` (to chyba defaultowo wyrzuca error)

*/
