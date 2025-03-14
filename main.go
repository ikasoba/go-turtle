package main

import (
	"fmt"
	"os"

	"github.com/ikasoba/go-turtle/parser"
	"github.com/k0kubun/pp"
)

func main() {
	filename := os.Args[1]

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	p := parser.New()

	n, se, err := p.Parse(f)

	if err != nil {
		panic(err)
	}

	if len(se) > 0 {
		fmt.Println(se)
	}

	pp.Println(n)
}
