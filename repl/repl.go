package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/amiosamu/tiny-interpreter/eval"
	"github.com/amiosamu/tiny-interpreter/lexer"
	"github.com/amiosamu/tiny-interpreter/object"
	"github.com/amiosamu/tiny-interpreter/parser"
)

const prompt = ">> "

// Start starts Tiny REPL.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	for {
		fmt.Print(prompt)
		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		// Process macros
		eval.DefineMacros(program, macroEnv)
		expanded := eval.ExpandMacros(program, macroEnv)

		// Evaluate AST
		evaluated := eval.Eval(expanded, env)
		if evaluated == nil {
			continue
		}

		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, msg)
		io.WriteString(out, "\n")
	}
}
