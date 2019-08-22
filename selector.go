package htmlutil

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Selector func(*html.Node) bool

func And(selector Selector, selectors ...Selector) Selector {
	return func(node *html.Node) bool {
		if !selector(node) {
			return false
		}

		for _, selector := range selectors {
			if !selector(node) {
				return false
			}
		}

		return true
	}
}

func (selector Selector) And(selector2 Selector) Selector {
	return func(node *html.Node) bool {
		return selector(node) && selector2(node)
	}
}

func Or(selector Selector, selectors ...Selector) Selector {
	return func(node *html.Node) bool {
		if selector(node) {
			return true
		}

		for _, selector := range selectors {
			if selector(node) {
				return true
			}
		}

		return false
	}
}

func (selector Selector) Or(selector2 Selector) Selector {
	return func(node *html.Node) bool {
		return selector(node) || selector2(node)
	}
}

func Not(selector Selector) Selector {
	return func(node *html.Node) bool {
		return !selector(node)
	}
}

func (selector Selector) Not() Selector {
	return func(node *html.Node) bool {
		return !selector(node)
	}
}

func Tree(selector1 Selector, selector2 ...Selector) Selector {
	selectors := append([]Selector{selector1}, selector2...)

	return func(node *html.Node) bool {
		for i := len(selectors) - 1; i >= 0; i-- {
			if node == nil {
				return false
			}

			if !selectors[i](node) {
				return false
			}

			node = node.Parent
		}

		return true
	}
}

func Includes(selector1 Selector, selector2 ...Selector) Selector {
	selectors := append([]Selector{selector1}, selector2...)

	return func(node *html.Node) bool {
		i := len(selectors) - 1

		if !selectors[i](node) {
			return false
		}
		node = node.Parent
		i--

		for ; i >= 0; i-- {
			for {
				if node == nil {
					return false
				}

				if selectors[i](node) {
					break
				}

				node = node.Parent
			}
		}

		return true
	}
}

func (selector Selector) FirstChild() Selector {
	return func(node *html.Node) bool {
		if node == nil {
			return false
		}

		return selector(node.Parent)
	}
}

func (selector Selector) Before() Selector {
	return func(node *html.Node) bool {
		if node == nil {
			return false
		}

		return selector(node.NextSibling)
	}
}

func (selector Selector) After() Selector {
	return func(node *html.Node) bool {
		if node == nil {
			return false
		}

		return selector(node.PrevSibling)
	}
}

func Elem(elem string) Selector {
	elem = strings.ToLower(elem)

	return func(node *html.Node) bool {
		if node == nil {
			return false
		}

		if node.Type != html.ElementNode {
			return false
		}

		return node.Data == elem
	}
}

func Atom(atom atom.Atom) Selector {
	return func(node *html.Node) bool {
		if node == nil {
			return false
		}

		return node.DataAtom == atom
	}
}

func HasAttr(key string) Selector {
	return func(node *html.Node) bool {
		if node == nil {
			return false
		}

		for _, attr := range node.Attr {
			if attr.Key == key {
				return true
			}
		}

		return false
	}
}

func Attr(key string, val string) Selector {
	return func(node *html.Node) bool {
		if node == nil {
			return false
		}

		for _, attr := range node.Attr {
			if attr.Key == key && attr.Val == val {
				return true
			}
		}

		return false
	}
}

func ID(id string) Selector {
	return Attr("id", id)
}

func Class(class string) Selector {
	return func(node *html.Node) bool {
		if node == nil {
			return false
		}

		for _, attr := range node.Attr {
			if attr.Key != "class" {
				continue
			}

			for _, className := range strings.Fields(attr.Val) {
				if class == className {
					return true
				}
			}
		}

		return false
	}
}
