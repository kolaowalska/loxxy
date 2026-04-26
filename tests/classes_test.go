package tests

import (
	"bytes"
	"testing"

	"github.com/kolaowalska/loxxy/src/evaluation"
	parser "github.com/kolaowalska/loxxy/src/parsing"
	resolving "github.com/kolaowalska/loxxy/src/resolving"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
	"github.com/kolaowalska/loxxy/src/testutils"
)

func TestClasses(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expected      any
		expectedError bool
	}{
		{
			name: "CLASSES - print class name",
			source: `
				class DevonshireCream {
					serveOn() {
						return "Scones";
					}
				}

				print DevonshireCream;
			`,
			expected:      "DevonshireCream\n",
			expectedError: false,
		},
		{
			name: "CLASSES - print class instance",
			source: `
			class Bagel {}
			var bagel = Bagel();
			print bagel;
			`,
			expected:      "Bagel instance\n",
			expectedError: false,
		},
		{
			name: "CLASSES - function call using the same syntax as a method call",
			source: `
			class Box {}

			fun notMethod(argument) {
			  print "called function with " + argument;
			}

			var box = Box();
			box.function = notMethod;
			box.function("argument");
			`,
			expected:      "called function with argument\n",
			expectedError: false,
		},
		{
			name: "CLASSES - bound method retains 'this' context",
			source: `
			class Person {
			  sayName() {
				print this.name;
			  }
			}

			var jane = Person();
			jane.name = "Jane";

			var method = jane.sayName;
			method(); // ?
						`,
			expected:      "Jane\n",
			expectedError: false,
		},
		{
			name: "CLASSES - bound method retains 'this' context 2.0",
			source: `
			class Person {
			  sayName() {
				print this.name;
			  }
			}

			var jane = Person();
			jane.name = "Jane";

			var bill = Person();
			bill.name = "Bill";

			bill.sayName = jane.sayName;
			bill.sayName(); // ?
			`,
			expected:      "Jane\n",
			expectedError: false,
		},
		{
			name: "CLASSES - hoist method lookup into a variable",
			source: `
			class Omelette {
			  filledWith(ingredient) {
				print "Omelette filled with " + ingredient;
			  }
			}

			var omelette = Omelette();
			// hoist the lookup part into a variable
			var eggs = omelette.filledWith; 

			// call it later
			eggs("cheese");
			`,
			expected:      "Omelette filled with cheese\n",
			expectedError: false,
		},
		{
			name: "CLASSES - pass bound method as a callback",
			source: `
			class Greeter {
			  sayHi(first, last) {
				print "Hi, " + first + " " + last + "! I am " + this.name;
			  }
			}

			// a standard function that takes a callback and executes it
			fun executeCallback(callback) {
			  callback("Dear", "Reader");
			}

			var greeter = Greeter();
			greeter.name = "Loxxy";

			// pass the bound method directly, without manually wrapping it in a function
			executeCallback(greeter.sayHi);
			`,
			expected:      "Hi, Dear Reader! I am Loxxy\n",
			expectedError: false,
		},
		{
			name: "CLASSES - execute class method",
			source: `
			class Bacon {
			  eat() {
				print "Crunch crunch crunch!";
			  }
			}

			Bacon().eat();
			`,
			expected:      "Crunch crunch crunch!\n",
			expectedError: false,
		},
		{
			name: "CLASSES - this inside method",
			source: `
			class Egotist {
			  speak() {
				print this;
			  }
			}

			var method = Egotist().speak;
			method();
		`,
			expected:      "Egotist instance\n",
			expectedError: false,
		},
		{
			name: "CLASSES - bound method",
			source: `
			class Cake {
			  taste() {
				var adjective = "delicious";
				print "The " + this.flavor + " cake is " + adjective + "!";
			  }
			}

			var cake = Cake();
			cake.flavor = "German chocolate";
			cake.taste();
			`,
			expected:      "The German chocolate cake is delicious!\n",
			expectedError: false,
		},
		{
			name: "CLASSES - ",
			source: `
			class Thing {
			  getCallback() {
				fun localFunction() {
				  print this;
				}

				return localFunction;
			  }
			}

			var callback = Thing().getCallback();
			callback();
			`,
			expected:      "Thing instance\n",
			expectedError: false,
		},
		{
			name: "CLASSES - this out of context",
			source: `
			print this;
			`,
			expected:      "nil\n",
			expectedError: true, // can't use 'this' outside of a class
		},
		// INIT
		{
			name: "CLASSES - init() always returns this",
			source: `
			class Foo {
			  init() {
				print this;
			  }
			}

			var foo = Foo();
			print foo.init();
			`,
			expected:      "Foo instance\nFoo instance\nFoo instance\n",
			expectedError: false,
		},
		{
			name: "CLASSES - static error, cannot return a value from an initializer",
			source: `
			class Foo {
			  init() {
				return "something else";
			  }
			}

			var foo = Foo();
			print foo.init();
			`,
			expected:      "nil\n",
			expectedError: true,
		},
		{
			name: "CLASSES - empty return statement in initializer",
			source: `
			class Foo {
			  init() {
				return;
			  }
			}

			var foo = Foo();
			print foo.init();
			`,
			expected:      "Foo instance\n",
			expectedError: false,
		},
		{
			name: "CLASSES - variable initialization inside init()",
			source: `
			class Foo {
			  init() {
				  var a = 1;
			  }
			}

			var foo = Foo();
			print foo.init();
			`,
			expected:      "Foo instance\n",
			expectedError: false,
		},
		{
			name: "CLASSES - variable initialization inside init()",
			source: `
			class Foo {
			  init() {
				  var a = 1;
				  return a;
			  }
			}

			var foo = Foo();
			print foo.init();
			`,
			expected:      "nil\n",
			expectedError: true,
		},
		{
			name: "CLASSES - undefined property init() inside class",
			source: `
			class Bagel {}
			var b = Bagel();
			print b.init();
			`,
			expected:      "nil\n",
			expectedError: true,
		},
		{
			name: "CLASSES - simple execution",
			source: `
			class Bagel {}
			Bagel();
			`,
			expected:      "",
			expectedError: false,
		},
		{
			name: "CLASSES - cannot declare fields in class body",
			source: `
				class Foo {
					var a = 2;
				}
				var foo = Foo();
				print foo.a;
			`,
			expected:      "",
			expectedError: true, // expect method name
		},
		{
			name: "CLASSES - overriding fields inside the class",
			source: `
			class Foo {
				init() {
					this.a = 1;
				}
				myFun(){
					this.a = 2;
				}
				getA(){
					return this.a;
				}
			}

			var foo = Foo();
			foo.init();
			foo.myFun();
			print foo.getA();
			`,
			expected:      "2\n",
			expectedError: false,
		},
		{
			name: "CLASSES - init() with parameters",
			source: `
				class Point {
					init(x, y) {
						this.x = x;
						this.y = y;
					}
				}
				var p = Point(1, 2);
				print p.x;
				print p.y;
			`,
			expected:      "1\n2\n",
			expectedError: false,
		},
		{
			name: "CLASSES - constructor arity mismatch",
			source: `
				class Point {
					init(x, y) {
						this.x = x;
						this.y = y;
					}
				}
				var p = Point(1); // missing 'y'
			`,
			expected:      "",
			expectedError: true, // runtimeError: expected 2 arguments but got 1
		},
		{
			name: "CLASSES - default constructor arity mismatch",
			source: `
				class Bagel {}
				var b = Bagel("cream cheese"); 
			`,
			expected:      "",
			expectedError: true, // runtimeError: Expected 0 arguments but got 1
		},
		{
			name: "CLASSES - fields shadow methods",
			source: `
				class Foo {
					bar() {
						print "method";
					}
				}
				var foo = Foo();
				foo.bar = "field"; 
				print foo.bar;
			`,
			expected:      "field\n",
			expectedError: false,
		},
		{
			name: "CLASSES - calling a non-callable field",
			source: `
				class Foo {}
				var foo = Foo();
				foo.bar = "string";
				foo.bar(); 
			`,
			expected:      "",
			expectedError: true, // runtimeError: can only call functions and classes
		},
		{
			name: "CLASSES - 'this' inside a normal function",
			source: `
				fun notMethod() {
					print this;
				}
			`,
			expected:      "",
			expectedError: true, // resolver Error: Can't use 'this' outside of a class
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reporter := &testutils.TestReporter{}
			s := scanner.NewScanner(test.source, reporter)
			tokens := s.ScanTokens()

			p := parser.NewParser(tokens, reporter)
			statements, err := p.Parse()

			if err != nil {
				if test.expectedError {
					return
				}
				t.Fatalf("Parser returned nil for source: %s\nError: %v", test.source, err)
			}

			var out bytes.Buffer
			i := evaluation.NewInterpreter() // i := evaluation.NewInterpreter(&out)
			i.Stdout = &out

			resolver := resolving.NewResolver(i, reporter)
			_ = resolver.ResolveStatements(statements)

			if testutils.CheckError(t, test.expectedError, nil, reporter.HadError, "RESOLVING") {
				return
			}

			err = i.Interpret(statements)
			if testutils.CheckError(t, test.expectedError, err, reporter.HadError, "INTERPRETING") {
				return
			}
			if test.expectedError {
				t.Fatalf("expected an error for source: %s, but execution succeeded.", test.source)
			}
		})
	}
}
