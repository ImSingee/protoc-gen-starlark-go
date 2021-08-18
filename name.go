package main

import (
	"google.golang.org/protobuf/compiler/protogen"
)

func StarlarkStructName(goName protogen.GoIdent) string {
	return goName.GoName + "Starlark"
}

func StarlarkFieldName(goName protogen.GoIdent) protogen.GoIdent {
	name := goName.GoName
	switch name {
	case "Name", "String", "Type", "Freeze", "Truth", "Hash", "Attr", "AttrNames",
		"Assert",
		"GetDefDoc", "GetSimpleDesc", "GetFullDesc":
		name += "_"
	}

	return protogen.GoIdent{
		GoName:       name,
		GoImportPath: goName.GoImportPath,
	}
}
