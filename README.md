htmlutil
====

html を雑に抽出できるようにする go ライブラリ

## Install

```
$ go get github.com/thamaji/htmlutil
```

## Example
Code
```
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/thamaji/htmlutil"
)

var html = `<table>
	<tbody>
		<tr>
			<td>1</td>
			<td>2</td>
			<td>3</td>
		</tr>
		<tr>
			<td>A</td>
			<td>B</td>
			<td>C</td>
		</tr>
		<tr>
			<td>ひとつめ</td>
			<td>ふたつめ</td>
			<td>みっつめ</td>
		</tr>
	</tbody>
</table>
`

func main() {
	doc, err := htmlutil.ParseFragment(strings.NewReader(html), htmlutil.ContextElem("body"))
	if err != nil {
		log.Fatal(err)
	}

	doc.FindAll(htmlutil.Tree(htmlutil.Elem("table"), htmlutil.Elem("tbody"), htmlutil.Elem("tr"))).Each(func(tr htmlutil.Selection) error {
		values, _ := tr.FindAll(htmlutil.Elem("td")).MapString(func(td htmlutil.Selection) (string, error) {
			return td.Text(htmlutil.DefaultFormatter), nil
		})

		fmt.Println(values)

		return nil
	})
}
```

Output
```
[1 2 3]
[A B C]
[ひとつめ ふたつめ みっつめ]
```