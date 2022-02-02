package htmlutil

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Selection []*html.Node

func (selection Selection) Nodes() []*html.Node {
	return selection
}

func (selection Selection) Text(formatter Formatter) string {
	buf := &strings.Builder{}

	for _, node := range selection {
		text, ok := formatter(node)
		fmt.Fprint(buf, text)

		if !ok {
			continue
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			fmt.Fprint(buf, Selection([]*html.Node{child}).Text(formatter))
		}
	}

	return buf.String()
}

func (selection Selection) Render(w io.Writer) error {
	for _, node := range selection {
		if err := html.Render(w, node); err != nil {
			return err
		}
	}
	return nil
}

func (selection Selection) String() string {
	buf := &strings.Builder{}
	selection.Render(buf)
	return buf.String()
}

func (selection Selection) Children() Selection {
	results := make([]*html.Node, 0, len(selection))

	for _, node := range selection {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			results = append(results, child)
		}
	}

	return results
}

func (selection Selection) Parents() Selection {
	results := make([]*html.Node, 0, len(selection))

	for _, node := range selection {
		results = append(results, node.Parent)
	}

	return results
}

func (selection Selection) AppendPrev(nodes ...*html.Node) {
	for _, selectionNode := range selection {
		for _, node := range nodes {
			node := *node
			node.NextSibling = selectionNode
			if prev := selectionNode.PrevSibling; prev != nil {
				prev.NextSibling = selectionNode
			}
			selectionNode.PrevSibling = &node
		}
	}
}

func (selection Selection) AppendNext(nodes ...*html.Node) {
	for _, selectionNode := range selection {
		for i := range nodes {
			node := *nodes[i]
			node.PrevSibling = selectionNode
			if next := selectionNode.NextSibling; next != nil {
				next.PrevSibling = &node
			}
			selectionNode.NextSibling = &node
		}
	}
}

func (selection Selection) AppendChild(nodes ...*html.Node) {
	for _, selectionNode := range selection {
		for i := range nodes {
			node := *nodes[i]
			selectionNode.AppendChild(&node)
		}
	}
}

func (selection Selection) Contains(selector Selector) bool {
	return selection.Find(selector).Len() > 0
}

func (selection Selection) Find(selector Selector) Selection {
	for _, node := range selection {
		if selector(node) {
			return []*html.Node{node}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if result := Selection([]*html.Node{child}).Find(selector); len(result) > 0 {
				return result
			}
		}
	}

	return []*html.Node{}
}

func (selection Selection) FindAll(selector Selector) Selection {
	result := []*html.Node{}

	for _, node := range selection {
		if selector(node) {
			result = append(result, node)
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			result = append(result, Selection([]*html.Node{child}).FindAll(selector)...)
		}
	}

	return result
}

func (selection Selection) Len() int {
	return len(selection)
}

func (selection Selection) Get(index int) Selection {
	if index >= selection.Len() || index < 0 {
		return []*html.Node{}
	}

	return []*html.Node{selection[index]}
}

func (selection Selection) First() Selection {
	return selection.Get(0)
}

func (selection Selection) Last() Selection {
	return selection.Get(selection.Len() - 1)
}

func (selection Selection) Each(f func(Selection) error) error {
	for _, node := range selection {
		if err := f(Selection([]*html.Node{node})); err != nil {
			return err
		}
	}

	return nil
}

func (selection Selection) Filter(selector Selector) Selection {
	result := make([]*html.Node, 0, len(selection))

	for _, node := range selection {
		if selector(node) {
			result = append(result, node)
		}
	}

	return result
}

func MapSelection[T any](selection Selection, f func(Selection) (T, error)) ([]T, error) {
	result := make([]T, 0, len(selection))

	for _, node := range selection {
		val, err := f(Selection([]*html.Node{node}))
		if err != nil {
			return result, err
		}

		result = append(result, val)
	}

	return result, nil
}

func (selection Selection) MapBool(f func(Selection) (bool, error)) ([]bool, error) {
	return MapSelection(selection, f)
}

func (selection Selection) MapString(f func(Selection) (string, error)) ([]string, error) {
	return MapSelection(selection, f)
}

func (selection Selection) MapByte(f func(Selection) (byte, error)) ([]byte, error) {
	return MapSelection(selection, f)
}

func (selection Selection) MapBytes(f func(Selection) ([]byte, error)) ([][]byte, error) {
	return MapSelection(selection, f)
}

func (selection Selection) MapRune(f func(Selection) (rune, error)) ([]rune, error) {
	return MapSelection(selection, f)
}

func (selection Selection) MapRunes(f func(Selection) ([]rune, error)) ([][]rune, error) {
	return MapSelection(selection, f)
}

func (selection Selection) MapInt(f func(Selection) (int, error)) ([]int, error) {
	return MapSelection(selection, f)
}

func (selection Selection) MapInt8(f func(Selection) (int8, error)) ([]int8, error) {
	return MapSelection(selection, f)
}

func (selection Selection) MapInt16(f func(Selection) (int16, error)) ([]int16, error) {
	return MapSelection(selection, f)
}

func (selection Selection) MapInt32(f func(Selection) (int32, error)) ([]int32, error) {
	return MapSelection(selection, f)
}

func (selection Selection) MapInt64(f func(Selection) (int64, error)) ([]int64, error) {
	return MapSelection(selection, f)
}

func (selection Selection) MapUint(f func(Selection) (uint, error)) ([]uint, error) {
	return MapSelection(selection, f)
}

func (selection Selection) MapUint8(f func(Selection) (uint8, error)) ([]uint8, error) {
	return MapSelection(selection, f)
}

func (selection Selection) MapUint16(f func(Selection) (uint16, error)) ([]uint16, error) {
	return MapSelection(selection, f)
}

func (selection Selection) MapUint32(f func(Selection) (uint32, error)) ([]uint32, error) {
	return MapSelection(selection, f)
}

func (selection Selection) MapUint64(f func(Selection) (uint64, error)) ([]uint64, error) {
	return MapSelection(selection, f)
}

func (selection Selection) MapUintptr(f func(Selection) (uintptr, error)) ([]uintptr, error) {
	return MapSelection(selection, f)
}

func (selection Selection) MapFloat32(f func(Selection) (float32, error)) ([]float32, error) {
	return MapSelection(selection, f)
}

func (selection Selection) MapFloat64(f func(Selection) (float64, error)) ([]float64, error) {
	return MapSelection(selection, f)
}

func (selection Selection) MapComplex64(f func(Selection) (complex64, error)) ([]complex64, error) {
	return MapSelection(selection, f)
}

func (selection Selection) MapComplex128(f func(Selection) (complex128, error)) ([]complex128, error) {
	return MapSelection(selection, f)
}
