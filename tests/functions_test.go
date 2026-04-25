package tests

import (
	"bytes"
	"testing"

	"github.com/kolaowalska/loxxy/src/evaluation"
	parser "github.com/kolaowalska/loxxy/src/parsing"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

func TestFunctions(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expected      any
		expectedError bool
	}{
		{
			name:          "CALL - native clock function",
			source:        "var t = clock(); print t > 0;",
			expected:      "true\n",
			expectedError: false,
		},
		{
			name:          "CALL - calling a non-function",
			source:        "var a = \"string\"; a();",
			expected:      "",
			expectedError: true,
		},
		{
			name:          "CALL - arity mismatch (too many)",
			source:        "fun foo(a) {} foo(1, 2);",
			expected:      "",
			expectedError: true,
		},
		{
			name:          "CALL - arity mismatch (too few)",
			source:        "fun foo(a, b) {} foo(1);",
			expected:      "",
			expectedError: true,
		},
		{
			name:          "BUILDER - simple function with no return",
			source:        "fun sayHi(first, last) { print \"Hi, \" + first + \" \" + last + \"!\"; } sayHi(\"Dear\", \"Reader\");",
			expected:      "Hi, Dear Reader!\n",
			expectedError: false,
		},
		{
			name:          "BUILDER - parameters are locally scoped",
			source:        "fun foo(a) { print a; } foo(1); print a;",
			expected:      "1\n",
			expectedError: true,
		},
		{
			name:          "RETURN - standard return value",
			source:        "fun add(a, b) { return a + b; } print add(10, 20);",
			expected:      "30\n",
			expectedError: false,
		},
		{
			name:          "RETURN - early return skips rest of function",
			source:        "fun early() { return \"done\"; print \"never happens\"; } print early();",
			expected:      "done\n",
			expectedError: false,
		},
		{
			name:          "RETURN - empty return evaluates to nil",
			source:        "fun empty() { return; } print empty();",
			expected:      "nil\n",
			expectedError: false,
		},
		{
			name:          "RETURN - nested inside control flow",
			source:        "fun isEven(n) { if (n == 2) return true; return false; } print isEven(2); print isEven(3);",
			expected:      "true\nfalse\n",
			expectedError: false,
		},
		{
			name:          "INTEGRATION - recursive fibonacci",
			source:        "fun fib(n) { if (n <= 1) return n; return fib(n - 2) + fib(n - 1); } print fib(7);",
			expected:      "13\n",
			expectedError: false,
		},
		{
			name:          "body must be a block",
			source:        "fun f() 123;",
			expected:      "",
			expectedError: true,
		},
		{
			name:          "empty body",
			source:        "fun f() {}",
			expected:      "",
			expectedError: false,
		},
		{
			name:          "local mutual recursion",
			source:        "{ fun isEven(n){ if (n==0) return true; return isOdd(n-1);} fun isOdd(n){ if (n==0) return false; return isEven(n-1};} isEven(4);}",
			expected:      "",
			expectedError: true,
		},
		{
			name:          "missing comma in parameters",
			source:        "fun foo(a, b c, d, e, f) {}",
			expected:      "",
			expectedError: true,
		},
		{
			name:          "mutual recursion",
			source:        "fun isEven(n){ if(n==0) return true; return isOdd(n-1); } fun isOdd(n) { if (n==0) return false; return isEven(n-1);} print isEven(4); print isOdd(3);",
			expected:      "true\ntrue\n",
			expectedError: false,
		},
		{
			name:          "print",
			source:        "fun foo() {} print foo; print clock;",
			expected:      "<fn foo>\n<native fn>\n",
			expectedError: false,
		},
		{
			name:          "nested call with arguments",
			source:        "fun returnArg(arg){ return arg;} fun returnFunCallWithArg(func, arg){ return returnArg(func)(arg); } fun printArg(arg){ print arg;} returnFunCallWithArg(printArg, \"hello world\");",
			expected:      "hello world\n",
			expectedError: false,
		},
		{
			name:          "too many arguments",
			source:        "fun foo() {}\n{\n  var a = 1;\n  foo(\n     a, // 1\n     a, // 2\n     a, // 3\n     a, // 4\n     a, // 5\n     a, // 6\n     a, // 7\n     a, // 8\n     a, // 9\n     a, // 10\n     a, // 11\n     a, // 12\n     a, // 13\n     a, // 14\n     a, // 15\n     a, // 16\n     a, // 17\n     a, // 18\n     a, // 19\n     a, // 20\n     a, // 21\n     a, // 22\n     a, // 23\n     a, // 24\n     a, // 25\n     a, // 26\n     a, // 27\n     a, // 28\n     a, // 29\n     a, // 30\n     a, // 31\n     a, // 32\n     a, // 33\n     a, // 34\n     a, // 35\n     a, // 36\n     a, // 37\n     a, // 38\n     a, // 39\n     a, // 40\n     a, // 41\n     a, // 42\n     a, // 43\n     a, // 44\n     a, // 45\n     a, // 46\n     a, // 47\n     a, // 48\n     a, // 49\n     a, // 50\n     a, // 51\n     a, // 52\n     a, // 53\n     a, // 54\n     a, // 55\n     a, // 56\n     a, // 57\n     a, // 58\n     a, // 59\n     a, // 60\n     a, // 61\n     a, // 62\n     a, // 63\n     a, // 64\n     a, // 65\n     a, // 66\n     a, // 67\n     a, // 68\n     a, // 69\n     a, // 70\n     a, // 71\n     a, // 72\n     a, // 73\n     a, // 74\n     a, // 75\n     a, // 76\n     a, // 77\n     a, // 78\n     a, // 79\n     a, // 80\n     a, // 81\n     a, // 82\n     a, // 83\n     a, // 84\n     a, // 85\n     a, // 86\n     a, // 87\n     a, // 88\n     a, // 89\n     a, // 90\n     a, // 91\n     a, // 92\n     a, // 93\n     a, // 94\n     a, // 95\n     a, // 96\n     a, // 97\n     a, // 98\n     a, // 99\n     a, // 100\n     a, // 101\n     a, // 102\n     a, // 103\n     a, // 104\n     a, // 105\n     a, // 106\n     a, // 107\n     a, // 108\n     a, // 109\n     a, // 110\n     a, // 111\n     a, // 112\n     a, // 113\n     a, // 114\n     a, // 115\n     a, // 116\n     a, // 117\n     a, // 118\n     a, // 119\n     a, // 120\n     a, // 121\n     a, // 122\n     a, // 123\n     a, // 124\n     a, // 125\n     a, // 126\n     a, // 127\n     a, // 128\n     a, // 129\n     a, // 130\n     a, // 131\n     a, // 132\n     a, // 133\n     a, // 134\n     a, // 135\n     a, // 136\n     a, // 137\n     a, // 138\n     a, // 139\n     a, // 140\n     a, // 141\n     a, // 142\n     a, // 143\n     a, // 144\n     a, // 145\n     a, // 146\n     a, // 147\n     a, // 148\n     a, // 149\n     a, // 150\n     a, // 151\n     a, // 152\n     a, // 153\n     a, // 154\n     a, // 155\n     a, // 156\n     a, // 157\n     a, // 158\n     a, // 159\n     a, // 160\n     a, // 161\n     a, // 162\n     a, // 163\n     a, // 164\n     a, // 165\n     a, // 166\n     a, // 167\n     a, // 168\n     a, // 169\n     a, // 170\n     a, // 171\n     a, // 172\n     a, // 173\n     a, // 174\n     a, // 175\n     a, // 176\n     a, // 177\n     a, // 178\n     a, // 179\n     a, // 180\n     a, // 181\n     a, // 182\n     a, // 183\n     a, // 184\n     a, // 185\n     a, // 186\n     a, // 187\n     a, // 188\n     a, // 189\n     a, // 190\n     a, // 191\n     a, // 192\n     a, // 193\n     a, // 194\n     a, // 195\n     a, // 196\n     a, // 197\n     a, // 198\n     a, // 199\n     a, // 200\n     a, // 201\n     a, // 202\n     a, // 203\n     a, // 204\n     a, // 205\n     a, // 206\n     a, // 207\n     a, // 208\n     a, // 209\n     a, // 210\n     a, // 211\n     a, // 212\n     a, // 213\n     a, // 214\n     a, // 215\n     a, // 216\n     a, // 217\n     a, // 218\n     a, // 219\n     a, // 220\n     a, // 221\n     a, // 222\n     a, // 223\n     a, // 224\n     a, // 225\n     a, // 226\n     a, // 227\n     a, // 228\n     a, // 229\n     a, // 230\n     a, // 231\n     a, // 232\n     a, // 233\n     a, // 234\n     a, // 235\n     a, // 236\n     a, // 237\n     a, // 238\n     a, // 239\n     a, // 240\n     a, // 241\n     a, // 242\n     a, // 243\n     a, // 244\n     a, // 245\n     a, // 246\n     a, // 247\n     a, // 248\n     a, // 249\n     a, // 250\n     a, // 251\n     a, // 252\n     a, // 253\n     a, // 254\n     a, // 255\n a);}",
			expected:      "",
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reporter := TestReporter{}
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
			err = i.Interpret(statements)

			if err != nil {
				if test.expectedError {
					return
				}
				t.Fatalf("Interpreter returned an error for source: %s\nError: %v", test.source, err)
			}
			if test.expectedError {
				t.Fatalf("Expected an error for source: %s, but execution succeeded.", test.source)
			}
			if out.String() != test.expected {
				t.Errorf("For source:\n%s\n\nExpected:\n%v\n\nGot:\n%v", test.source, test.expected, out.String())
			}
		})
	}
}
