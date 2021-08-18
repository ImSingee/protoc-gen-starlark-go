package main

import "google.golang.org/protobuf/compiler/protogen"

type MessageDesc struct {
	*protogen.Message

	Fields []*FieldDesc
}

type FieldDesc struct {
	*protogen.Field

	Name string
}
