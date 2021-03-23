package resolver

import (
	"context"
	"reflect"
	"sync"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type FanoutResult struct {
	Messages []proto.Message
	Err      error
}

func NewFanoutResult(pbSlice interface{}, err error) FanoutResult {
	return FanoutResult{
		Messages: MessageSlice(pbSlice),
		Err:      err,
	}
}

func NewSingleFanoutResult(message proto.Message, err error) FanoutResult {
	return FanoutResult{
		Messages: []proto.Message{message},
		Err:      err,
	}
}

type FanoutHandler interface {
	Add(delta int)
	Done()

	Cancelled() <-chan struct{}
	Channel() chan<- FanoutResult

	Results(limit uint32) (*Results, error)
}

type fanoutHandler struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	ch chan FanoutResult
}

func (r *fanoutHandler) Add(delta int)                { r.wg.Add(delta) }
func (r *fanoutHandler) Done()                        { r.wg.Done() }
func (r *fanoutHandler) Cancelled() <-chan struct{}   { return r.ctx.Done() }
func (r *fanoutHandler) Channel() chan<- FanoutResult { return r.ch }
func (r *fanoutHandler) Results(limit uint32) (*Results, error) {
	results := make([]proto.Message, 0, limit)
	var failures []*status.Status

	done := make(chan struct{})
	go func() {
		r.wg.Wait()
		close(done)
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for limit == 0 || len(results) < int(limit) {
			select {
			case <-done:
				return
			case result := <-r.ch:
				if result.Err != nil {
					failures = append(failures, status.Convert(result.Err))
				} else {
					results = append(results, result.Messages...)
				}
			}
		}
		r.cancel()
	}()
	wg.Wait()

	return &Results{Messages: results, PartialFailures: failures}, nil
}

func NewFanoutHandler(ctx context.Context) (context.Context, FanoutHandler) {
	ctx, cancel := context.WithCancel(ctx)
	return ctx, &fanoutHandler{
		ctx:    ctx,
		ch:     make(chan FanoutResult),
		cancel: cancel,
	}
}

// MessageSlice takes a slice of protobuf objects and converts them to a slice of generic protobuf Messages.
func MessageSlice(s interface{}) []proto.Message {
	rs := reflect.ValueOf(s)
	ret := make([]proto.Message, rs.Len())
	for i := 0; i < rs.Len(); i++ {
		ret[i] = rs.Index(i).Interface().(proto.Message)
	}
	return ret
}
