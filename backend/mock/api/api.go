package apimock

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v3"
)

func FromYAML(s string, message proto.Message) proto.Message {
	err := protojson.Unmarshal(yamlToJSON(s), message)
	if err != nil {
		panic(err)
	}

	return message
}

func AnyFromYAML(s string) *anypb.Any {
	a := &anypb.Any{}

	var obj map[string]interface{}
	if err := yaml.Unmarshal([]byte(s), &obj); err != nil {
		panic(err)
	}

	if err := protojson.Unmarshal(yamlToJSON(s), a); err != nil {
		panic(err)
	}

	return a
}

func yamlToJSON(s string) []byte {
	var obj map[string]interface{}
	if err := yaml.Unmarshal([]byte(s), &obj); err != nil {
		panic(err)
	}

	// Encode YAML to JSON.
	rawJSON, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return rawJSON
}
