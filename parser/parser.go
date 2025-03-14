package parser

import (
	"io"

	"github.com/ikasoba/go-turtle/ast"
	"github.com/ikasoba/go-turtle/grammar"
)

type TurtleParser struct{}

func New() *TurtleParser {
	return &TurtleParser{}
}

func (_ *TurtleParser) Parse(r io.Reader) (ast.AstNode, []*grammar.SyntaxError, error) {
	toks, err := grammar.NewTokenStream(r)
	if err != nil {
		return nil, nil, err
	}

	gram := grammar.NewGrammar()
	tb := NewTurtleSyntaxTreeBuilder()
	p, err := grammar.NewParser(toks, gram, grammar.SemanticAction(
		grammar.NewASTActionSet(gram, tb),
	))
	if err != nil {
		return nil, nil, err
	}

	err = p.Parse()
	if err != nil {
		return nil, nil, err
	}

	synErrs := p.SyntaxErrors()

	return tb.Tree(), synErrs, nil
}
