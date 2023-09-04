package main

import (
	"flag"
	"fmt"
	"interpreter/glox"
	"os"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		return
	}
	filename := flag.Arg(0)
	// todo: REPL?
	if err := runFile(filename); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runFile(filename string) error {
	fBytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return run(fBytes)
}

func run(source []byte) error {
	scanner := glox.NewScanner(source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		return err
	}
	// glox.NewParser(tokens).PrintAST()
	parser := glox.NewParser(tokens)
	env := glox.NewEnvironment()
	parser.Execute(env)
	return nil
}
