package middleware

import (
	"strings"

	"github.com/gobwas/glob"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Factory map[string]func(*any.Any, *zap.Logger, tally.Scope) (Middleware, error)

type Middleware interface {
	UnaryInterceptor() grpc.UnaryServerInterceptor
}

func SplitFullMethod(fullMethod string) (service string, method string, ok bool) {
	s := strings.SplitN(fullMethod, "/", 3)
	if len(s) != 3 {
		return "serviceUnknown", "methodUnknown", false
	}
	return s[1], s[2], true
}

func MatchMethodOrResource(pattern, input string) bool {
	if pattern == input || pattern == "*" {
		return true
	}

	g, err := glob.Compile(pattern, '/')
	if err != nil {
		return false
	}
	return g.Match(input)
}
