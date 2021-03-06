package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

type wellKnownProvider struct {
	goImportPath protogen.GoImportPath
	goName       string
	goPointer    bool

	convertImport protogen.GoImportPath
	convertFunc   string

	modifier func(string) string
}

func (w *wellKnownProvider) Name() (protogen.GoIdent, bool) {
	return protogen.GoIdent{
		GoName:       w.goName,
		GoImportPath: w.goImportPath,
	}, w.goPointer
}
func (w *wellKnownProvider) ConvertFunc() protogen.GoIdent {
	return protogen.GoIdent{
		GoImportPath: w.convertImport,
		GoName:       w.convertFunc,
	}
}

func (w *wellKnownProvider) Modify(s string) string {
	if w.modifier != nil {
		return w.modifier(s)
	} else {
		return s
	}
}

var CustomMap = map[string]OverrideProvider{
	"google.protobuf.Timestamp": &wellKnownProvider{
		goImportPath:  "go.starlark.net/starlark/lib/time",
		goName:        "Time",
		convertImport: "go.starlark.net/starlark/lib/time",
		convertFunc:   "NewDateTime",
		modifier:      func(s string) string { return s + ".AsTime()" },
	},
	"struct.Value": &wellKnownProvider{
		goImportPath: starlarkPackage,
		goName:       "Value",
		modifier:     func(s string) string { return s + ".ToStarlark()" },
	},
	"struct.Dict": &wellKnownProvider{
		goImportPath: starlarkPackage,
		goName:       "Dict",
		goPointer:    true,
		modifier:     func(s string) string { return s + ".ToStarlark()" },
	},
	"struct.List": &wellKnownProvider{
		goImportPath: starlarkPackage,
		goName:       "List",
		goPointer:    true,
		modifier:     func(s string) string { return s + ".ToStarlark()" },
	},
}

func IsWellKnownType(fullname string) bool {
	return strings.HasPrefix(fullname, "google.protobuf.")
}
