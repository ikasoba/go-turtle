package parser

import (
	"strconv"

	"github.com/ikasoba/go-turtle/ast"
	"github.com/ikasoba/go-turtle/grammar"
)

type TurtleSyntaxTreeBuilder struct {
	tree *ast.WorkingAstNode
}

func NewTurtleSyntaxTreeBuilder() *TurtleSyntaxTreeBuilder {
	return &TurtleSyntaxTreeBuilder{}
}

func (b *TurtleSyntaxTreeBuilder) Shift(kindName string, text string, row, col int) grammar.SyntaxTreeNode {
	res := &ast.WorkingAstNode{
		Type:     grammar.NodeTypeTerminal,
		KindName: kindName,
		Text:     text,
	}

	switch kindName {
	case "a":
		t := ast.PrefixedName{
			Name: "rdf:type",
		}

		res.Node = ast.AstNode(&t)
	}

	return res
}

func (b *TurtleSyntaxTreeBuilder) ShiftError(kindName string) grammar.SyntaxTreeNode {
	return &ast.WorkingAstNode{
		Type:     grammar.NodeTypeError,
		KindName: kindName,
	}
}

func (b *TurtleSyntaxTreeBuilder) Reduce(kindName string, children []grammar.SyntaxTreeNode) grammar.SyntaxTreeNode {
	switch kindName {
	case "iri", "object", "predicate", "verb", "directive", "literal", "statement", "string":
		return children[0]
	}

	cNodes := make([]grammar.SyntaxTreeNode, len(children))

	for i, c := range children {
		cNodes[i] = c.(*ast.WorkingAstNode)
	}

	res := &ast.WorkingAstNode{
		Type:     grammar.NodeTypeNonTerminal,
		KindName: kindName,
		Children: cNodes,
	}

	switch kindName {
	case "iriref":
		if len(children) == 3 {
			n := children[1].(*ast.WorkingAstNode)

			u := ""

			children := n.Children
			for _, x := range children {
				n := x.(*ast.WorkingAstNode)

				u += n.Text
			}

			iri := ast.IRI{
				IRI: u,
			}

			res.Node = ast.AstNode(&iri)
		} else {
			iri := ast.IRI{
				IRI: "",
			}

			res.Node = ast.AstNode(&iri)
		}

	case "prefixed_name", "pname_ns":
		name := ast.PrefixedName{
			Name: children[0].(*ast.WorkingAstNode).Text,
		}

		res.Node = ast.AstNode(&name)

	case "predicate_object_list_atom":
		v := children[0].(*ast.WorkingAstNode)
		o := children[1].(*ast.WorkingAstNode)

		objects := make([]ast.AstNode, len(o.Children))

		for i, c := range o.Children {
			objects[i] = c.(*ast.WorkingAstNode).Node
		}

		predicate := ast.Predicate{
			Verb:    v.Node,
			Objects: objects,
		}

		res.Node = ast.AstNode(&predicate)

	case "turtle_doc":
		cNodes := make([]ast.AstNode, len(children))

		for i, c := range children {
			cNodes[i] = c.(*ast.WorkingAstNode).Node
		}

		doc := ast.Document{
			Children: cNodes,
		}

		res.Node = ast.AstNode(&doc)

	case "triples":
		n := children[0].(*ast.WorkingAstNode)

		if n.KindName == "subject" {
			triple := ast.Triple{
				Subject:    n.Children[0].(*ast.WorkingAstNode).Node,
				Predicates: []ast.Predicate{},
			}

			for _, x := range children[1].(*ast.WorkingAstNode).Children {
				n := x.(*ast.WorkingAstNode)
				p, ok := n.Node.(*ast.Predicate)
				if !ok {
					continue
				}

				triple.Predicates = append(triple.Predicates, *p)
			}

			res.Node = ast.AstNode(&triple)
		} else {
			res.Node = children[0].(*ast.WorkingAstNode).Node
		}

	case "prefix_id":
		n := children[0].(*ast.WorkingAstNode)

		prefix := ast.PrefixDirective{
			Name: *n.Node.(*ast.PrefixedName),
			IRI:  *children[1].(*ast.WorkingAstNode).Node.(*ast.IRI),
		}

		res.Node = ast.AstNode(&prefix)

	case "string_literal_quote":
		t := ""

		for _, x := range children[1].(*ast.WorkingAstNode).Children {
			n := x.(*ast.WorkingAstNode).Children[0].(*ast.WorkingAstNode)

			switch n.KindName {
			case "string_literal_quote_echar":
				switch n.Text[1:] {
				case "t":
					t += "\t"

				case "b":
					t += "\b"

				case "r":
					t += "\r"

				case "f":
					t += "\f"

				default:
					t += n.Text[1:]
				}

			case "string_literal_quote_uchar":
				h := n.Text[2:]

				i, _ := strconv.ParseInt(h, 16, 64)

				t += string(rune(i))

			default:
				t += n.Text
			}
		}

		s := ast.StringLiteral{
			Value: t,
		}

		res.Children = []grammar.SyntaxTreeNode{}

		res.Node = ast.AstNode(&s)

	case "rdf_literal":
		n := children[0].(*ast.WorkingAstNode).Node.(*ast.StringLiteral)

		if len(children) == 3 {
			n.IRI = children[2].(*ast.WorkingAstNode).Node
		} else if len(children) == 2 {
			t := children[1].(*ast.WorkingAstNode).Text
			n.LangDir = &t
		}

		res.Node = n

	case "blank_node_property_list":
		l := ast.BlankNodePropertyList{
			Properties: []ast.Predicate{},
		}

		for _, x := range children {
			n, ok := x.(*ast.WorkingAstNode).Node.(*ast.Predicate)
			if !ok {
				continue
			}

			l.Properties = append(l.Properties, *n)
		}

		res.Node = &l

	case "collection":
		c := ast.Collection{
			Objects: []ast.AstNode{},
		}

		for _, x := range children[1].(*ast.WorkingAstNode).Children {
			c.Objects = append(c.Objects, x.(*ast.WorkingAstNode).Node)
		}

		res.Node = &c
	}

	return res
}

func (b *TurtleSyntaxTreeBuilder) Accept(f grammar.SyntaxTreeNode) {
	w := f.(*ast.WorkingAstNode)

	b.tree = w
}

func (b *TurtleSyntaxTreeBuilder) Tree() ast.AstNode {
	return b.tree.Node
}

func (b *TurtleSyntaxTreeBuilder) RawTree() *ast.WorkingAstNode {
	return b.tree
}
