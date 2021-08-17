package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

func GenerateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + "_starlark.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	f := newFileInfo(file)

	g.P("// Code generated by protoc-gen-starlark-go. DO NOT EDIT.")
	g.P()
	g.P("package ", f.GoPackageName)
	g.P()

	for i, imps := 0, f.Desc.Imports(); i < imps.Len(); i++ {
		genImport(gen, g, f, imps.Get(i))
	}

	for _, message := range f.allMessages {
		genMessage(gen, g, f, message)
	}

	return g
}

func genImport(gen *protogen.Plugin, g *protogen.GeneratedFile, f *fileInfo, imp protoreflect.FileImport) {
	impFile, ok := gen.FilesByPath[imp.Path()]
	if !ok {
		return
	}
	if impFile.GoImportPath == f.GoImportPath {
		// Don't generate imports or aliases for types in the same Go package.
		return
	}
	// Generate imports for all non-weak dependencies, even if they are not
	// referenced, because other code and tools depend on having the
	// full transitive closure of protocol buffer types in the binary.
	if !imp.IsWeak {
		g.Import(impFile.GoImportPath)
	}
	if !imp.IsPublic {
		return
	}

	// Generate public imports by generating the imported file, parsing it,
	// and extracting every symbol that should receive a forwarding declaration.
	impGen := GenerateFile(gen, impFile)
	impGen.Skip()
	b, err := impGen.Content()
	if err != nil {
		gen.Error(err)
		return
	}
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "", b, parser.ParseComments)
	if err != nil {
		gen.Error(err)
		return
	}
	genForward := func(tok token.Token, name string, expr ast.Expr) {
		// Don't import unexported symbols.
		r, _ := utf8.DecodeRuneInString(name)
		if !unicode.IsUpper(r) {
			return
		}
		// Don't import the FileDescriptor.
		if name == impFile.GoDescriptorIdent.GoName {
			return
		}
		// Don't import decls referencing a symbol defined in another package.
		// i.e., don't import decls which are themselves public imports:
		//
		//	type T = somepackage.T
		if _, ok := expr.(*ast.SelectorExpr); ok {
			return
		}
		g.P(tok, " ", name, " = ", impFile.GoImportPath.Ident(name))
	}
	g.P("// Symbols defined in public import of ", imp.Path(), ".")
	g.P()
	for _, decl := range astFile.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					genForward(decl.Tok, spec.Name.Name, spec.Type)
				case *ast.ValueSpec:
					for i, name := range spec.Names {
						var expr ast.Expr
						if i < len(spec.Values) {
							expr = spec.Values[i]
						}
						genForward(decl.Tok, name.Name, expr)
					}
				case *ast.ImportSpec:
				default:
					panic(fmt.Sprintf("can't generate forward for spec type %T", spec))
				}
			}
		}
	}
	g.P()
}

/*
var ResponseType = helper.NewStruct("Response",
	helper.StringArg("url").Short("原始请求 URL"),
	helper.IntArg("status_code").Short("返回值"),
	helper.DictArg("headers", helper.StringArg("header_name"), helper.ListArg("header_values", helper.StringArg("value"))),
	helper.FuncArg("body", nil, helper.StringArg("")),
	helper.FuncArg("json", nil, helper.StringArg("")).Short("body 信息，并进行 JSON decode"),
).WithAssertProvider(func(s *helper.Struct) (pass bool, failReason string) {
	statusCode := s.Values["status_code"]

	if helper.EqualInt(statusCode, 200) {
		return true, ""
	} else {
		return false, fmt.Sprintf("status code = %v != 200", statusCode)
	}
})
*/

func genMessage(gen *protogen.Plugin, g *protogen.GeneratedFile, f *fileInfo, m *messageInfo) {
	if m.Desc.IsMapEntry() {
		return
	}

	// Message type declaration.
	g.Annotate(m.GoIdent.GoName, m.Location)
	g.P("type ", m.GoIdent, `Starlark struct{`)
	genMessageFields(gen, g, f, m)
	g.P("}")
	g.P()

	//genMessageKnownFunctions(g, f, m)
	//genMessageDefaultDecls(g, f, m)
	//genMessageMethods(g, f, m)
	//genMessageOneofWrapperTypes(g, f, m)
}

func genMessageFields(gen *protogen.Plugin, g *protogen.GeneratedFile, f *fileInfo, m *messageInfo) {
	sf := f.allMessageFieldsByPtr[m]
	for _, field := range m.Fields {
		genMessageField(gen, g, f, m, field, sf)
	}
}

func genMessageField(gen *protogen.Plugin, g *protogen.GeneratedFile, f *fileInfo, m *messageInfo, field *protogen.Field, sf *structFields) {
	//if oneof := field.Oneof; oneof != nil && !oneof.Desc.IsSynthetic() {
	//	// It would be a bit simpler to iterate over the oneofs below,
	//	// but generating the field here keeps the contents of the Go
	//	// struct in the same order as the contents of the source
	//	// .proto file.
	//	//if oneof.Fields[0] != field {
	//	//	return // only generate for first appearance
	//	//}
	//	//
	//	//g.Annotate(m.GoIdent.GoName+"."+oneof.GoName, oneof.Location)
	//	//leadingComments := oneof.Comments.Leading
	//	//if leadingComments != "" {
	//	//	leadingComments += "\n"
	//	//}
	//	//ss := []string{fmt.Sprintf(" Types that are assignable to %s:\n", oneof.GoName)}
	//	//for _, field := range oneof.Fields {
	//	//	ss = append(ss, "\t*"+field.GoIdent.GoName+"\n")
	//	//}
	//	//leadingComments += protogen.Comments(strings.Join(ss, ""))
	//	//g.P(leadingComments,
	//	//	oneof.GoName, " ", oneofInterfaceName(oneof), tags)
	//	//sf.append(oneof.GoName)
	//
	//	// TODO
	//
	//	return
	//}

	g.P(field.GoName, " ", fieldStarlarkType(gen, g, f, field))
}

// fieldStarlarkType returns the Starlark type used for a field.
func fieldStarlarkType(gen *protogen.Plugin, g *protogen.GeneratedFile, f *fileInfo, field *protogen.Field) string {
	switch {
	case field.Desc.IsList():
		// TODO Enhance it
		return `*` + g.QualifiedGoIdent(starlarkList)
	case field.Desc.IsMap():
		return `*` + g.QualifiedGoIdent(starlarkMap)
	}

	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		return g.QualifiedGoIdent(starlarkBool)
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind, protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return g.QualifiedGoIdent(starlarkInt)
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return g.QualifiedGoIdent(starlarkFloat)
	case protoreflect.StringKind:
		return g.QualifiedGoIdent(starlarkString)
	case protoreflect.BytesKind:
		return g.QualifiedGoIdent(starlarkBytes)
	case protoreflect.EnumKind:
		// TODO: enhance support enum
		//goType = g.QualifiedGoIdent(field.Enum.GoIdent)
		return g.QualifiedGoIdent(starlarkInt)
	case protoreflect.MessageKind, protoreflect.GroupKind:
		full := string(field.Message.Desc.FullName())
		if custom := CustomMap[full]; custom != nil {
			name, pointer := custom.Name()
			prefix := ""
			if pointer {
				prefix = "*"
			}
			return prefix + g.QualifiedGoIdent(name)
		}
		if ext := GetMessageExtensionFor(field.Message); ext != nil {
			switch {
			case ext.GetDisable():
				return ""
			case ext.GetToString():
				return g.QualifiedGoIdent(starlarkString)
			}

			if custom := ext.GetCustom(); custom != nil {
				importPath := GetImportPath(custom.GetStarlarkTypePackage(), field.Message.GoIdent.GoImportPath)
				importName := custom.GetStarlarkTypeName()
				prefix := ""
				if strings.HasPrefix(importName, "*") {
					prefix = "*"
					importName = importName[1:]
				}
				return prefix + g.QualifiedGoIdent(protogen.GoIdent{GoImportPath: importPath, GoName: importName})
			}

			gen.Error(fmt.Errorf("invalid option"))
			return ""
		}
		if IsWellKnownType(full) {
			_, _ = os.Stderr.WriteString("unsupported well-known type " + full + "\n")
			//gen.Error(fmt.Errorf("unsupported well-known type %s", full))
			return g.QualifiedGoIdent(starlarkValue)
		}

		return `*` + g.QualifiedGoIdent(field.Message.GoIdent) + "Starlark"
	}

	gen.Error(fmt.Errorf("unknown type (not supported)"))
	return ""
}
