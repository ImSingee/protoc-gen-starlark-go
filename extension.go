package main

import (
	"github.com/ImSingee/protoc-gen-starlark-go/options"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

type MessageExtension struct {
	*options.MessageOption
}

func GetMessageExtensionFor(message *protogen.Message) *MessageExtension {
	opts := message.Desc.Options().(*descriptorpb.MessageOptions)
	if opts == nil || !proto.HasExtension(opts, options.E_MessageOption) {
		return nil
	}

	ext := proto.GetExtension(opts, options.E_MessageOption).(*options.MessageOption)

	return &MessageExtension{ext}
}
