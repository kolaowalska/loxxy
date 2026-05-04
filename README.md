# loxxy
> a high-fidelity tree-walk interpreter for the lox programming language, written entirely in go.

![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/emiliamacek/da44b047ba1e0c16c4f14224e28e2618/raw/coverage.json)
![Go version](https://img.shields.io/github/go-mod/go-version/kolaowalska/loxxy)
![License](https://img.shields.io/github/license/kolaowalska/loxxy)


### language & interpreter features
the project implements a full suite of modern language features and navigates through the standard pipeline of lexical analysis, parsing, semantic analysis, and interpretation.

#### core implementation
- **tokens and lexing**: thorough scanning of source text into tokens
- **abstract syntax trees (ast)**: representation of program structure for both expressions and statements
- **recursive descent parsing**: a top-down parsing strategy to ensure syntactic correctness
- **prefix and infix expressions**: support for mathematical and logical operations
- **runtime representation of objects**: dynamic handling of lox values in the go environment
- **code interpretation**: utilizing idiomatic go type switches for ast traversal
- **lexical scope**: accurate block-level scoping and variable visibility
- **environment chains**: nested storage for variable tracking across scopes

#### advanced semantics
- **control flow**: support for branching and looping
- **functions and closures**: supporting parameters and the ability to capture the surrounding lexical environment
- **static variable resolution**: a dedicated pass analyzing variable bindings prior to execution for error detection
- **object-oriented programming**: complete implementation of classes, including constructors, fields, and methods
- **inheritance**: support for single inheritance using the `<` operator and method resolution via `super` calls

[//]: # (features:)
[//]: # (- **turing-complete scripting:** support for basic arithmetic, booleans, strings, and standard control flow)
[//]: # (- **first-class functions:** support for closures, anonymous functions, and returning functions)
[//]: # (- **object-oriented programming:** full support for classes, methods, instance fields, `this` binding, and constructors)
[//]: # (- **inheritance:** single inheritance through the `<` operator, along with `super` method resolution)
[//]: # (- **semantic analysis:** a `Resolver` phase that statically analyzes the syntax tree before execution to cinch lexical scoping, catch unused variables, and prevent invalid returns)
[//]: # (- **repl & file execution:** interactively from the command line or through a `.lox` script file)

### structure tour d'horizon
the architecture is structured chronologically by compiler phase to maximize modularity and exploit go's innate package system
1. **scanning** (`src/scanning/`) - lexical analysis
2. **parsing** (`src/parsing/`) - syntax analysis 
3. **representation** (`src/representation/`) - syntax tree node definitions, i.e. expressions and statements
4. **resolving** (`src/resolving/`) - semantic analysis 
5. **evaluation** (`src/evaluation/`) - runtime environment 
6. **testing** (`tests/`, `src/testutils/`) - tests and utilities for the testing pipeline


## quick start
### prerequisites 
- Go version 1.26.1 or higher

### installation
1. clone the repository to your local machine 
2. build the binary
~~~bash
git clone https://github.com/kolaowalska/loxxy.git
cd loxxy
go build -o loxxy main.go
~~~

### usage
loxxy can be used in REPL mode or script execution mode.

- **interactive REPL**:
run 
~~~bash 
./loxxy
~~~
to start a prompt, for example
~~~bash
> var greeting = "feeling... foxxy ;)";
> print greeting;
feeling... foxxy ;)
~~~

- **script execution**:
run 
~~~bash
./loxxy path/to/script
~~~
to execute a `.lox` file.

## testing 
the project is equipped with an extensive testing suite that utilizes a central `testutils` package to validate the scanner, parser, resolver, and interpreter in a unified pipeline.

the following command runs all available tests:
~~~bash
go test ./...
~~~

the tests for specialized language features are as follows:
- **classes and methods**: `tests/classes_test.go`
- **control flow**: `tests/control_flow_test.go`
- **inheritance**: `tests/inheritance_test.go`
- **function closures**: `tests/functions_test.go`

---


