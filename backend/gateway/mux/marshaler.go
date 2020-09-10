package mux

import (
	"io"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	runtimev2 "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// protojsonShim is a simple struct that allow us to move to the protojson marshaler of grpc-gateway v2 without fully migrating to v2.
// This replaces the original jsonpb-based marshaler which has a number of bugs.
// TODO: delete once migrated to grpc-gateway v2.
type protojsonShim struct {
	marshalerv2 *runtimev2.JSONPb
}

func (p *protojsonShim) Marshal(v interface{}) ([]byte, error) {
	return p.marshalerv2.Marshal(v)
}

func (p *protojsonShim) Unmarshal(data []byte, v interface{}) error {
	return p.marshalerv2.Unmarshal(data, v)
}

func (p *protojsonShim) NewDecoder(r io.Reader) runtime.Decoder {
	return p.marshalerv2.NewDecoder(r)
}

func (p *protojsonShim) NewEncoder(w io.Writer) runtime.Encoder {
	return p.marshalerv2.NewEncoder(w)
}

func (p *protojsonShim) ContentType() string {
	return p.marshalerv2.ContentType(nil)
}
