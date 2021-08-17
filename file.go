package main

import (
	"unicode"
	"unicode/utf8"

	"google.golang.org/protobuf/compiler/protogen"
)

type fileInfo struct {
	*protogen.File

	allEnums    []*enumInfo
	allMessages []*messageInfo

	allEnumsByPtr         map[*enumInfo]int    // value is index into allEnums
	allMessagesByPtr      map[*messageInfo]int // value is index into allMessages
	allMessageFieldsByPtr map[*messageInfo]*structFields
}

type structFields struct {
	count      int
	unexported map[int]string
}

func (sf *structFields) append(name string) {
	if r, _ := utf8.DecodeRuneInString(name); !unicode.IsUpper(r) {
		if sf.unexported == nil {
			sf.unexported = make(map[int]string)
		}
		sf.unexported[sf.count] = name
	}
	sf.count++
}

func newFileInfo(file *protogen.File) *fileInfo {
	f := &fileInfo{File: file}

	// Collect all enums, messages, and extensions in "flattened ordering".
	// See filetype.TypeBuilder.
	var walkMessages func([]*protogen.Message, func(*protogen.Message))
	walkMessages = func(messages []*protogen.Message, f func(*protogen.Message)) {
		for _, m := range messages {
			f(m)
			walkMessages(m.Messages, f)
		}
	}
	initEnumInfos := func(enums []*protogen.Enum) {
		for _, enum := range enums {
			f.allEnums = append(f.allEnums, newEnumInfo(f, enum))
		}
	}
	initMessageInfos := func(messages []*protogen.Message) {
		for _, message := range messages {
			f.allMessages = append(f.allMessages, newMessageInfo(f, message))
		}
	}

	initEnumInfos(f.Enums)
	initMessageInfos(f.Messages)

	walkMessages(f.Messages, func(m *protogen.Message) {
		initEnumInfos(m.Enums)
		initMessageInfos(m.Messages)
	})

	// Derive a reverse mapping of enum and message pointers to their index
	// in allEnums and allMessages.
	if len(f.allEnums) > 0 {
		f.allEnumsByPtr = make(map[*enumInfo]int)
		for i, e := range f.allEnums {
			f.allEnumsByPtr[e] = i
		}
	}
	if len(f.allMessages) > 0 {
		f.allMessagesByPtr = make(map[*messageInfo]int)
		f.allMessageFieldsByPtr = make(map[*messageInfo]*structFields)
		for i, m := range f.allMessages {
			f.allMessagesByPtr[m] = i
			f.allMessageFieldsByPtr[m] = new(structFields)
		}
	}

	return f
}

type enumInfo struct {
	*protogen.Enum
}

func newEnumInfo(f *fileInfo, enum *protogen.Enum) *enumInfo {
	e := &enumInfo{Enum: enum}
	return e
}

type messageInfo struct {
	*protogen.Message

	hasWeak bool
}

func newMessageInfo(f *fileInfo, message *protogen.Message) *messageInfo {
	m := &messageInfo{Message: message}
	for _, field := range m.Fields {
		m.hasWeak = m.hasWeak || field.Desc.IsWeak()
	}
	return m
}
