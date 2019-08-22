package htmlutil

import (
	"io"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ParseTables(selection Selection, formatter Formatter) []Table {
	tables := []Table{}

	selection.FindAll(Atom(atom.Table)).Each(func(selection Selection) error {
		table := Table{}

		if selection.Contains(Atom(atom.Tbody)) {
			selection = selection.Find(Atom(atom.Tbody))
		}

		selection.FindAll(Atom(atom.Tr)).Each(func(selection Selection) error {
			row, _ := selection.FindAll(Atom(atom.Td)).MapString(func(selection Selection) (string, error) {
				return selection.Text(formatter), nil
			})

			table = append(table, row)

			return nil
		})

		tables = append(tables, table)

		return nil
	})

	return tables
}

type Table [][]string

func (table Table) Render(w io.Writer) error {
	return html.Render(w, table.Node())
}

func (table Table) Node() *html.Node {
	node := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Table,
		Data:     atom.Table.String(),
	}

	node.AppendChild(func() *html.Node {
		node := &html.Node{
			Type:     html.ElementNode,
			DataAtom: atom.Tbody,
			Data:     atom.Tbody.String(),
		}

		for _, row := range table {
			node.AppendChild(func() *html.Node {
				node := &html.Node{
					Type:     html.ElementNode,
					DataAtom: atom.Tr,
					Data:     atom.Tr.String(),
				}

				for _, value := range row {
					node.AppendChild(func() *html.Node {
						node := &html.Node{
							Type:     html.ElementNode,
							DataAtom: atom.Td,
							Data:     atom.Td.String(),
						}

						node.AppendChild(&html.Node{
							Type: html.TextNode,
							Data: value,
						})

						return node
					}())
				}

				return node
			}())
		}

		return node
	}())

	return node
}
