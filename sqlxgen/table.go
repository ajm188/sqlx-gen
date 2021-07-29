package sqlxgen

import (
	"fmt"
	"strings"

	"vitess.io/vitess/go/sqltypes"
	"vitess.io/vitess/go/vt/sqlparser"
)

type Table struct {
	Name    string
	Columns []*Column
}

type Column struct {
	Name string
	Type sqlparser.ColumnType
}

func (c *Column) GoName() (name string) {
	name = strings.ReplaceAll(c.Name, "_", " ") // hello_world => hello world
	name = strings.Title(name)                  // hello world => Hello World
	name = strings.ReplaceAll(name, " ", "")    // Hello World => HelloWorld
	return name
}

func (c *Column) GoType() string {
	t := c.Type.SQLType()
	switch {
	case c.Type.Type == "tinyint": // special handling for booleans
		return "bool"
	case sqltypes.IsSigned(t):
		return "int64"
	case sqltypes.IsUnsigned(t):
		return "uint64"
	case sqltypes.IsFloat(t) || t == sqltypes.Decimal:
		return "float64"
	case sqltypes.IsText(t):
		return "string"
	case sqltypes.IsBinary(t):
		return "[]byte"
	default:
		panic(fmt.Sprintf("unsupported query.Type: %s (%d)", t.String(), t))
	}
}

func (c *Column) StructTag() string {
	return fmt.Sprintf("`json:\"%s\" db:\"%s\"`", c.Name, c.Name)
}
