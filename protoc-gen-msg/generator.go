package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	plugin "github.com/davyxu/pbmeta/proto/compiler"
)

type Generator struct {
	*bytes.Buffer

	Request  *plugin.CodeGeneratorRequest  // The input.
	Response *plugin.CodeGeneratorResponse // The output.
	indent   string
}

func New() *Generator {
	self := new(Generator)
	self.Buffer = new(bytes.Buffer)
	self.Request = new(plugin.CodeGeneratorRequest)
	self.Response = new(plugin.CodeGeneratorResponse)
	return self
}

// Error reports a problem, including an error, and exits the program.
func (self *Generator) Error(err error, msgs ...string) {
	s := strings.Join(msgs, " ") + ":" + err.Error()
	log.Errorln("protoc-gen-sharpnet, error:", s)
	os.Exit(1)
}

// Fail reports a problem and exits the program.
func (self *Generator) Fail(msgs ...string) {
	s := strings.Join(msgs, " ")
	log.Errorln("protoc-gen-sharpnet, error:", s)
	os.Exit(1)
}

func (self *Generator) Print(str ...interface{}) {

	for _, v := range str {
		switch s := v.(type) {
		case string:
			self.WriteString(s)
		case *string:
			self.WriteString(*s)
		case bool:
			fmt.Fprintf(self, "%t", s)
		case *bool:
			fmt.Fprintf(self, "%t", *s)
		case int, int32:
			fmt.Fprintf(self, "%d", s)
		case *int32:
			fmt.Fprintf(self, "%d", *s)
		case *int64:
			fmt.Fprintf(self, "%d", *s)
		case float64:
			fmt.Fprintf(self, "%g", s)
		case *float64:
			fmt.Fprintf(self, "%g", *s)
		default:
			panic("here")
			self.Fail(fmt.Sprintf("unknown type in printer: %T", v))

		}
	}

}

func (self *Generator) Println(str ...interface{}) {
	self.BeginLine()
	self.Print(str...)
	self.EndLine()
}

func (self *Generator) BeginLine() {
	self.WriteString(self.indent)
}

func (self *Generator) EndLine() {
	self.WriteByte('\n')
}

// In Indents the output one tab stop.
func (self *Generator) In() { self.indent += "\t" }

// Out unindents the output one tab stop.
func (self *Generator) Out() {
	if len(self.indent) > 0 {
		self.indent = self.indent[1:]
	}
}
