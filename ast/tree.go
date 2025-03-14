package ast

import (
	"github.com/ikasoba/go-turtle/grammar"
)

type WorkingAstNode struct {
	Type     grammar.NodeType
	KindName string
	Children []grammar.SyntaxTreeNode
	Text     string
	Node     AstNode
}

func (n *WorkingAstNode) ChildCount() int {
	return len(n.Children)
}

func (n *WorkingAstNode) ExpandChildren() []grammar.SyntaxTreeNode {
	return n.Children
}

type NodeKind int

const (
	KindDocument NodeKind = iota
	KindPrefixDirective
	KindBaseDirective
	KindSparqlPrefixDirective
	KindSparqlBaseDirective
	KindTriple
	KindPredicate
	KindIRI
	KindPrefixedName
	KindStringLiteral
	KindBlankNodePropertyList
	KindCollection
)

type AstNode interface {
	Kind() NodeKind
}

type Document struct {
	Children []AstNode
}

func (_ *Document) Kind() NodeKind {
	return KindDocument
}

type PrefixDirective struct {
	Name PrefixedName
	IRI  IRI
}

func (_ *PrefixDirective) Kind() NodeKind {
	return KindPrefixDirective
}

type BaseDirective struct {
	IRI IRI
}

func (_ *BaseDirective) Kind() NodeKind {
	return KindBaseDirective
}

type SparqlPrefixDirective struct {
	Name string
	IRI  IRI
}

func (_ *SparqlPrefixDirective) Kind() NodeKind {
	return KindSparqlPrefixDirective
}

type SparqlBaseDirective struct {
	IRI IRI
}

func (_ *SparqlBaseDirective) Kind() NodeKind {
	return KindSparqlBaseDirective
}

type Triple struct {
	Subject    AstNode
	Predicates []Predicate
}

func (_ *Triple) Kind() NodeKind {
	return KindTriple
}

type Predicate struct {
	Verb    AstNode
	Objects []AstNode
}

func (_ *Predicate) Kind() NodeKind {
	return KindPredicate
}

type IRI struct {
	IRI string
}

func (_ *IRI) Kind() NodeKind {
	return KindIRI
}

type PrefixedName struct {
	Name string
}

func (_ *PrefixedName) Kind() NodeKind {
	return KindPrefixedName
}

type StringLiteral struct {
	Value   string
	LangDir *string
	IRI     AstNode
}

func (_ *StringLiteral) Kind() NodeKind {
	return KindStringLiteral
}

type BlankNodePropertyList struct {
	Properties []Predicate
}

func (_ *BlankNodePropertyList) Kind() NodeKind {
	return KindBlankNodePropertyList
}

type Collection struct {
	Objects []AstNode
}

func (_ *Collection) Kind() NodeKind {
	return KindCollection
}
