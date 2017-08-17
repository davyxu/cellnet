package main

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type generateOption struct {
	formatGoCode bool

	outputData []byte
}

// 字段首字母大写
func publicFieldName(name string) string {
	return strings.ToUpper(string(name[0])) + name[1:]
}

func generateCode(templateName, templateStr, output string, model interface{}, opt *generateOption) {

	var err error

	if opt == nil {
		opt = &generateOption{}
	}

	var bf bytes.Buffer

	tpl, err := template.New(templateName).Parse(templateStr)
	if err != nil {
		goto OnError
	}

	err = tpl.Execute(&bf, model)
	if err != nil {
		goto OnError
	}

	if opt.formatGoCode {
		if err = formatCode(&bf); err != nil {
			fmt.Println("format golang code err", err)
		}
	}

	opt.outputData = bf.Bytes()

	if output != "" {

		os.MkdirAll(filepath.Dir(output), 666)

		err = ioutil.WriteFile(output, bf.Bytes(), 0666)

		if err != nil {
			goto OnError
		}
	}
	return

OnError:
	fmt.Println(err)
	os.Exit(1)
}

func formatCode(bf *bytes.Buffer) error {

	fset := token.NewFileSet()

	ast, err := parser.ParseFile(fset, "", bf, parser.ParseComments)
	if err != nil {
		return err
	}

	bf.Reset()

	err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(bf, fset, ast)
	if err != nil {
		return err
	}

	return nil
}
