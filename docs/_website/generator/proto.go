package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/yoheimuta/go-protoparser/v4"
	"github.com/yoheimuta/go-protoparser/v4/interpret/unordered"
	"github.com/yoheimuta/go-protoparser/v4/parser"
)

type protoScope struct {
	messages map[string]*unordered.Message
}

func (p *protoScope) getSimpleMessageYAML(name string) string {
	m, ok := p.messages[name]
	if !ok {
		panic("could not find " + name)
	}

	fields := m.MessageBody.Fields

	// Flatten oneofs to regular fields since that's how they work in YAML/JSON (though mututally exclusive).
	for _, o := range m.MessageBody.Oneofs {
		for _, field := range o.OneofFields {
			fields = append(fields, &parser.Field{
				FieldName: field.FieldName,
				Type:      field.Type,
				Comments:  field.Comments,
			})
		}
	}

	var b strings.Builder
	for _, field := range fields {
		var val string
		val = fmt.Sprintf("<%s>", field.Type)

		switch field.Type {
		case "string":
			val = fmt.Sprintf("\"%s\"", val)
		}

		if isMessageType(field.Type) {
			// i.e. is a message type
			val = fmt.Sprintf("{%s}", val)
		}

		if field.IsRepeated {
			val = fmt.Sprintf("[%s, ...]", val)
		}

		fmt.Fprintf(&b, "%s: %s\n", field.FieldName, val)
	}
	return strings.TrimSpace(b.String())
}

func protoToMessages(p *unordered.Proto) map[string]*unordered.Message {
	packageName := p.ProtoBody.Packages[0].Name

	ret := make(map[string]*unordered.Message, len(p.ProtoBody.Messages))
	for _, m := range p.ProtoBody.Messages {
		qualifiedName := fmt.Sprintf("%s.%s", packageName, m.MessageName)
		ret[qualifiedName] = m
	}
	return ret
}

func newProtoScope(path string) (*protoScope, error) {
	files, err := getFiles(path, ".proto")
	if err != nil {
		return nil, err
	}

	ret := &protoScope{messages: map[string]*unordered.Message{}}
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil {
			return nil, err
		}
		pp, err := protoparser.Parse(fh)
		if err != nil {
			return nil, err
		}
		p, err := protoparser.UnorderedInterpret(pp)
		if err != nil {
			return nil, err
		}

		for k, v := range protoToMessages(p) {
			ret.messages[k] = v
		}
	}

	return ret, nil
}

func isMessageType(s string) bool {
	ss := strings.Split(s, ".")
	return startsWithUpper(ss[len(ss)-1])
}

func startsWithUpper(s string) bool {
	if s == "" {
		return false
	}
	return 'A' <= s[0] && s[0] <= 'Z'
}
