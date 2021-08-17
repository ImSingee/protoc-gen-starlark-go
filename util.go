package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

func FullToIdent(s string) protogen.GoIdent {
	index := strings.LastIndexByte(s, '.')
	if index == -1 {
		return protogen.GoIdent{
			GoName: s,
		}
	} else {
		return protogen.GoIdent{
			GoImportPath: protogen.GoImportPath(s[:index]),
			GoName:       s[index+1:],
		}
	}
}
