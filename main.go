package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"

	"vitess.io/vitess/go/vt/sqlparser"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	path := flag.String("schema", "create_commerce_schema.sql", "path to a file containing CREATE TABLE statements")

	flag.Parse()

	if *path == "" {
		log.Fatal("must specify -schema")
	}

	data, err := ioutil.ReadFile(*path)
	check(err)

	tok := sqlparser.NewStringTokenizer(string(data))

	var stmt sqlparser.Statement
	for err == nil {
		stmt, err = sqlparser.ParseNext(tok)
		if stmt == nil {
			continue
		}

		switch stmt := stmt.(type) {
		case *sqlparser.CreateTable:
			log.Print("hello in CreateTable case")
		default:
			buf := sqlparser.NewTrackedBuffer(nil)
			buf.Myprintf("%v", stmt)
			log.Printf("[warn] %s is not a CreateTable (type: %T), skipping ...", buf.String(), stmt)
		}
	}

	if err != io.EOF {
		check(err)
	}
}
