package htmlutil

import "golang.org/x/net/html"

type Formatter func(*html.Node) string

func DefaultFormatter(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}

	return ""
}
