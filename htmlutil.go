package htmlutil

import (
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func Parse(r io.Reader) (Selection, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	return []*html.Node{node}, nil
}

func ParseFragment(r io.Reader, context *html.Node) (Selection, error) {
	nodes, err := html.ParseFragment(r, context)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func ContextElem(elem string) *html.Node {
	return ContextAtom(atom.Lookup([]byte(strings.ToLower(elem))))
}

func ContextAtom(atom atom.Atom) *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		Data:     atom.String(),
		DataAtom: atom,
	}
}
