package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"vitess.io/vitess/go/vt/sqlparser"

	"github.com/ajm188/sqlx-gen/sqlxgen"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	path := flag.String("schema", "create_commerce_schema.sql", "path to a file containing CREATE TABLE statements")
	pkgName := flag.String("pkg", "models", "package name to generate for")

	flag.Parse()

	if *path == "" {
		log.Fatal("must specify -schema")
	}

	data, err := ioutil.ReadFile(*path)
	check(err)

	tok := sqlparser.NewStringTokenizer(string(data))

	var (
		stmt   sqlparser.Statement
		tables []*sqlxgen.Table
	)

	for err == nil {
		stmt, err = sqlparser.ParseNext(tok)
		if err != nil {
			break
		}

		switch stmt := stmt.(type) {
		case *sqlparser.CreateTable:
			table := &sqlxgen.Table{
				Name: strings.Title(stmt.Table.Name.String()),
			}

			tables = append(tables, table)
		default:
			buf := sqlparser.NewTrackedBuffer(nil)
			buf.Myprintf("%v", stmt)
			log.Printf("[warn] %s is not a CreateTable (type: %T), skipping ...", buf.String(), stmt)
		}
	}

	if err != io.EOF {
		check(err)
	}

	err = sqlxgen.Generate(os.Stdout, &sqlxgen.Info{
		PackageName: *pkgName,
		Tables:      tables,
	})
	check(err)
}
