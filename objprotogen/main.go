package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
)

var paramOut = flag.String("out", "", "output file name")

func main() {

	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		fmt.Println("Require go source file list")
		os.Exit(1)
	}

	fs := token.NewFileSet()

	p := &Package{}

	for _, filename := range files {

		file, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = p.Parse(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	genCode(*paramOut, p)

}
