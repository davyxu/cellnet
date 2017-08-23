package main

import (
	"github.com/davyxu/cellnet/util"
	"go/ast"
)

type Package struct {
	PackageName string

	Structs []*Struct
}

func (self *Package) Parse(fileNode *ast.File) error {

	self.PackageName = fileNode.Name.Name

	ast.Inspect(fileNode, func(n ast.Node) bool {

		switch typeSpec := n.(type) {
		case *ast.TypeSpec:

			switch typeSpecType := typeSpec.Type.(type) {
			case *ast.StructType:

				st := &Struct{
					Name: typeSpec.Name.Name,
				}

				st.Parse(typeSpecType)

				self.Structs = append(self.Structs, st)
			}
		}

		return true
	})

	return nil
}

func newFile() *Package {
	return &Package{}
}

type Struct struct {
	Name string

	Fields []*Field
}

func (self *Struct) MsgID() uint32 {
	return util.StringHash(self.Name)
}

func (self *Struct) Parse(n *ast.StructType) {

	for _, fd := range n.Fields.List {

		f := &Field{}
		f.Parse(fd)

		self.Fields = append(self.Fields, f)
	}
}

type Field struct {
	Name string
	Type string
}

func (self *Field) Parse(n *ast.Field) {

	self.Name = n.Names[0].Name

	switch x := n.Type.(type) {
	case *ast.Ident:
		self.Type = x.Name
	}
}
