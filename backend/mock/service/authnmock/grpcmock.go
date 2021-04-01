package authnmock

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type MockServerTransportStream struct {
	SentHeaders []metadata.MD
}

func (MockServerTransportStream) Method() string                 { return "" }
func (MockServerTransportStream) SetHeader(md metadata.MD) error { return nil }
func (m *MockServerTransportStream) SendHeader(md metadata.MD) error {
	m.SentHeaders = append(m.SentHeaders, md)
	return nil
}
func (MockServerTransportStream) SetTrailer(md metadata.MD) error { return nil }

var _ grpc.ServerTransportStream = &MockServerTransportStream{}
