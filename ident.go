package main

import "google.golang.org/protobuf/compiler/protogen"

const (
	reflectPackage = protogen.GoImportPath("reflect")
	sortPackage    = protogen.GoImportPath("sort")
	stringsPackage = protogen.GoImportPath("strings")
	syncPackage    = protogen.GoImportPath("sync")
	timePackage    = protogen.GoImportPath("time")
	utf8Package    = protogen.GoImportPath("unicode/utf8")
)

const (
	starlarkPackage       = protogen.GoImportPath("go.starlark.net/starlark")
	starlarkhelperPackage = protogen.GoImportPath("github.com/ImSingee/starlarkhelper")
)

var (
	starlarkBool   = protogen.GoIdent{GoName: "Bool", GoImportPath: starlarkPackage}
	starlarkString = protogen.GoIdent{GoName: "String", GoImportPath: starlarkPackage}
	starlarkInt    = protogen.GoIdent{GoName: "Int", GoImportPath: starlarkPackage}
	starlarkFloat  = protogen.GoIdent{GoName: "Float", GoImportPath: starlarkPackage}
	starlarkBytes  = protogen.GoIdent{GoName: "Bytes", GoImportPath: starlarkPackage}
	starlarkList   = protogen.GoIdent{GoName: "List", GoImportPath: starlarkPackage}
	starlarkMap    = protogen.GoIdent{GoName: "Dict", GoImportPath: starlarkPackage}  // TODO Enhance it
	starlarkValue  = protogen.GoIdent{GoName: "Value", GoImportPath: starlarkPackage} // TODO Enhance it
)
