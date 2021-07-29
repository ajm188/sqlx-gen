package sqlxgen

import (
	"fmt"
	"strings"
)

type Table struct {
	Name    string
	Columns []*Column
}

type Column struct {
	Name string
	Type string
}

func (c *Column) GoName() (name string) {
	name = strings.ReplaceAll(c.Name, "_", " ") // hello_world => hello world
	name = strings.Title(name)                  // hello world => Hello World
	name = strings.ReplaceAll(name, " ", "")    // Hello World => HelloWorld
	return name
}

func (c *Column) StructTag() string {
	return fmt.Sprintf("`json:\"%s\" db:\"%s\"`", c.Name, c.Name)
}
