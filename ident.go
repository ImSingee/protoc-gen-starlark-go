package main

import "google.golang.org/protobuf/compiler/protogen"

const (
	reflectPackage = protogen.GoImportPath("reflect")
	fmtPackage     = protogen.GoImportPath("fmt")
	sortPackage    = protogen.GoImportPath("sort")
	stringsPackage = protogen.GoImportPath("strings")
	syncPackage    = protogen.GoImportPath("sync")
	timePackage    = protogen.GoImportPath("time")
	utf8Package    = protogen.GoImportPath("unicode/utf8")
)

var (
	fmtSprintf = protogen.GoIdent{GoName: "Sprintf", GoImportPath: fmtPackage}
	fmtErrorf  = protogen.GoIdent{GoName: "Errorf", GoImportPath: fmtPackage}
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
	starlarkDict   = protogen.GoIdent{GoName: "Dict", GoImportPath: starlarkPackage}
	starlarkMap    = protogen.GoIdent{GoName: "Dict", GoImportPath: starlarkPackage}  // TODO Enhance it
	starlarkValue  = protogen.GoIdent{GoName: "Value", GoImportPath: starlarkPackage} // TODO Enhance it

	starlarkHasAttrs = protogen.GoIdent{GoName: "HasAttrs", GoImportPath: starlarkPackage}

	starlarkNone = protogen.GoIdent{GoName: "None", GoImportPath: starlarkPackage}

	starlarkMakeInt64  = protogen.GoIdent{GoName: "MakeInt64", GoImportPath: starlarkPackage}
	starlarkMakeUint64 = protogen.GoIdent{GoName: "MakeUint64", GoImportPath: starlarkPackage}
	starlarkNewList    = protogen.GoIdent{GoName: "NewList", GoImportPath: starlarkPackage}
	starlarkNewDict    = protogen.GoIdent{GoName: "NewDict", GoImportPath: starlarkPackage}

	starlarkhelperCanHelp          = protogen.GoIdent{GoName: "CanHelp", GoImportPath: starlarkhelperPackage}
	starlarkhelperHelpMode         = protogen.GoIdent{GoName: "HelpMode", GoImportPath: starlarkhelperPackage}
	starlarkhelperHelpModeTerminal = protogen.GoIdent{GoName: "HelpModeTerminal", GoImportPath: starlarkhelperPackage}
	starlarkhelperDocProvider      = protogen.GoIdent{GoName: "DocProvider", GoImportPath: starlarkhelperPackage}
	starlarkhelperNewDocObject     = protogen.GoIdent{GoName: "NewDocObject", GoImportPath: starlarkhelperPackage}
)

func GetImportPath(name string, default_ protogen.GoImportPath) protogen.GoImportPath {
	switch name {
	case "":
		return default_
	case "STARLARK":
		return starlarkPackage
	case "STARLARKHELPER", "HELPER":
		return starlarkhelperPackage
	default:
		return protogen.GoImportPath(name)
	}
}
