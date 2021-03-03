# front

Extracts front matter.

Modified fork of github.com/gernest/front
## Features
* Custom delimiters (any three character string. e.g `+++`,  `$$$`,  `---`,  `%%%`)
* Supports YAML frontmatter
* Supports JSON frontmatter
* JSON or Map output
* Supports additional delimiters in the body

## Installation

	go get github.com/barshociaj/front

## How to use

```go
package main

import (
	"fmt"
	"strings"

	"github.com/barshociaj/front"
)

var txt = `+++
{
    "title":"front"
}
+++

# Body
Over my dead body
`

func main() {
	m := front.NewMatter("+++")
	front, body, err := m.JSONToMap(strings.NewReader(txt))
	if err != nil {
		panic(err)
	}

	fmt.Printf("The front matter is:\n%#v\n", front)
	fmt.Printf("The body is:\n%q\n", body)
}
```

### `m.YAMLToJSON(fileReader)`

Convert YAML front matter to JSON (`[]byte`)
```go
m := front.NewMatter("---")
front, body, err := m.YAMLToJSON(fileReader)
```

### `m.YAMLToMap(fileReader)`

Convert YAML front matter to a map (`map[string]interface{}`)
```go
m := front.NewMatter("---")
front, body, err := m.YAMLToMap(fileReader)
```

### `m.JSONToMap(fileReader)`

Convert JSON front matter to a map (`map[string]interface{}`)
```go
m := front.NewMatter("---")
front, body, err := m.JSONToMap(fileReader)
```

## Licence

This project is under the MIT Licence. See the [LICENCE](LICENCE) file for the full license text.

