package main

import "google.golang.org/protobuf/compiler/protogen"

type OverrideProvider interface {
	Name() (protogen.GoIdent, bool) // starlark type name (include * if needed)
	ConvertFunc() protogen.GoIdent  // convert function
	Modify(s string) string
}

type CustomOverrideProvider struct {
	name        protogen.GoIdent
	pointer     bool
	convertFunc protogen.GoIdent
}

func (c *CustomOverrideProvider) Name() (protogen.GoIdent, bool) { return c.name, c.pointer }

func (c *CustomOverrideProvider) ConvertFunc() protogen.GoIdent { return c.convertFunc }

func (c *CustomOverrideProvider) Modify(s string) string { return s }
