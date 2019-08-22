package htmlutil

import (
	"path"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Formatter func(*html.Node) (string, bool)

func DefaultFormatter(node *html.Node) (string, bool) {
	if node.Type == html.TextNode {
		return node.Data, true
	}

	return "", true
}

func ConfluenceStorageFormatter(node *html.Node) (string, bool) {
	switch node.Type {
	case html.TextNode:
		return node.Data, true

	case html.ElementNode:
		switch node.Data {
		case "time":
			for _, attr := range node.Attr {
				if attr.Key == "datetime" {
					return attr.Val, false
				}
			}

			return "", false

		case "ac:task-list":
			return "", true

		case "ac:task":
			return "", true

		case "ac:task-status":
			return "", true

		case "ac:task-body":
			return "", true

		case "ac:placeholder":
			return "", false

		case "ac:link":
			return "", true

		case "ac:plain-text-link-body":
			return "", true

		case "ac:link-body":
			return "", true

		case "ac:image":
			// 必要であれば対応
			// for _, attr := range node.Attr {
			// 	switch attr.Key {
			// 	case "ac:align":
			// 	case "ac:border":
			// 	case "ac:class":
			// 	case "ac:title":
			// 	case "ac:style":
			// 	case "ac:thumbnail":
			// 	case "ac:alt":
			// 	case "ac:height":
			// 	case "ac:width":
			// 	case "ac:vspace":
			// 	case "ac:hspace":
			// 	}
			// }

			return "", true

		case "ac:layout":
			return "", true

		case "ac:layout-section":
			// 必要であれば対応
			// for _, attr := range node.Attr {
			// 	switch attr.Key {
			// 	case "ac:typ":
			// 		switch attr.Val {
			// 		case "single":
			// 		case "two_equal":
			// 		case "two_left_sidebar":
			// 		case "two_right_sidebar":
			// 		case "three_equal":
			// 		case "three_with_sidebars":
			// 		}
			// 	}
			// }
			return "", true

		case "ac:layout-cell":
			return "", true

		case "ac:emoticon":
			for _, attr := range node.Attr {
				if attr.Key == "ac:name" {
					return ":" + attr.Val + ":", false
				}
			}

			return "", false

		case "ac:structured-macro":
			macroName := ""
			for _, attr := range node.Attr {
				if attr.Key == "ac:name" {
					macroName = attr.Val
					break
				}
			}

			switch macroName {
			case "status":
				for child := node.FirstChild; child != nil; child = child.NextSibling {
					if child.Data != "ac:parameter" {
						continue
					}

					parameterName := ""
					for _, attr := range child.Attr {
						if attr.Key == "ac:name" {
							parameterName = attr.Val
							break
						}
					}

					switch parameterName {
					case "colour":
						continue

					case "title":
						if child.FirstChild == nil {
							return "", false
						}

						return "[" + child.FirstChild.Data + "]", false

					default:
						continue
					}
				}

				return "", false

			default:
				return "", false
			}

		case "ri:page":
			spaceKey := ""
			contentTitle := ""

			for _, attr := range node.Attr {
				switch attr.Key {
				case "ri:space-key":
					spaceKey = attr.Val

				case "ri:content-title":
					contentTitle = attr.Val
				}
			}

			return path.Join(spaceKey, contentTitle), false

		case "ri:blog-post":
			spaceKey := ""
			contentTitle := ""
			postingDay := ""

			for _, attr := range node.Attr {
				switch attr.Key {
				case "ri:space-key":
					spaceKey = attr.Val

				case "ri:content-title":
					contentTitle = attr.Val

				case "ri:posting-day":
					postingDay = "(" + attr.Val + ")"
				}
			}

			return path.Join(spaceKey, contentTitle) + postingDay, false

		case "ri:attachment":
			for _, attr := range node.Attr {
				if attr.Key == "ri:filename" {
					return attr.Val, true
				}
			}

			return "", true

		case "ri:url":
			for _, attr := range node.Attr {
				if attr.Key == "ri:value" {
					return attr.Val, false
				}
			}

			return "", false

		case "ri:shortcut":
			key := ""
			parameter := ""

			for _, attr := range node.Attr {
				switch attr.Key {
				case "ri:key":
					key = attr.Val

				case "ri:parameter":
					parameter = attr.Val
				}
			}

			return "[" + parameter + "@" + key + "]", false

		case "ri:user":
			for _, attr := range node.Attr {
				if attr.Key == "ri:userkey" {
					return attr.Val, false
				}
			}

			return "", false

		case "ri:space":
			for _, attr := range node.Attr {
				if attr.Key == "ri:space-key" {
					return attr.Val, false
				}
			}

			return "", false

		case "ri:content-entity":
			for _, attr := range node.Attr {
				if attr.Key == "ri:content-id" {
					return attr.Val, false
				}
			}

			return "", false

		case "at:declarations", "at:string", "at:textarea", "at:list", "at:option":
			return "", false

		case "at:var":
			for _, attr := range node.Attr {
				if attr.Key == "at:name" {
					return "${" + attr.Val + "}", false
				}
			}

			return "", false

		default:
			return "", true
		}

	default:
		return "", true
	}
}

func ConfluenceViewFormatter(node *html.Node) (string, bool) {
	switch node.Type {
	case html.TextNode:
		return node.Data, true

	case html.ElementNode:
		switch node.DataAtom {
		case atom.Time:
			for _, attr := range node.Attr {
				if attr.Key == "datetime" {
					return attr.Val, false
				}
			}

			return "", true

		case atom.A:
			for _, attr := range node.Attr {
				if attr.Key != "data-username" {
					return attr.Val, false
				}
			}

			return "", true

		default:
			return "", true
		}

	default:
		return "", true
	}
}
