package htmlutil

import "golang.org/x/net/html"

type Formatter func(*html.Node) (string, bool)

func DefaultFormatter(node *html.Node) (string, bool) {
	if node.Type == html.TextNode {
		return node.Data, true
	}

	return "", true
}
