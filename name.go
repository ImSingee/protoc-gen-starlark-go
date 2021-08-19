package main

import (
	"google.golang.org/protobuf/compiler/protogen"
)

func StarlarkStructName(goName protogen.GoIdent) protogen.GoIdent {
	return protogen.GoIdent{
		GoName:       goName.GoName + "Starlark",
		GoImportPath: goName.GoImportPath,
	}
}

func StarlarkFieldName(goIdent protogen.GoIdent, name string) protogen.GoIdent {
	switch name {
	case "Name", "String", "Type", "Freeze", "Truth", "Hash", "Attr", "AttrNames",
		"Assert",
		"GetDefDoc", "GetSimpleDesc", "GetFullDesc":
		name += "_"
	}

	return protogen.GoIdent{
		GoName:       name,
		GoImportPath: goIdent.GoImportPath,
	}
}
