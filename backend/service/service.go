package service

import (
	"github.com/golang/protobuf/ptypes/any"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

type Service interface{}

type Factory map[string]func(*any.Any, *zap.Logger, tally.Scope) (Service, error)

var Registry = map[string]Service{}

// TODO: create a one-way registry that errors on duplicates and can be locked after instantiation for additional safety.
