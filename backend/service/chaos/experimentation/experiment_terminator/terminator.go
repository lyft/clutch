package experiment_terminator

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	experimentationv1 "github.com/lyft/clutch/backend/api/chaos/experimentation/v1"
	"github.com/lyft/clutch/backend/service/chaos/experimentation/experimentstore"
)

type TerminationCriteria interface {
	ShouldTerminate(experimentStarted time.Time, config interface{}) error
}
type Terminator struct {
	store              experimentstore.Storer
	enabledConfigTypes []string
	criterias          []TerminationCriteria
}

func NewTerminator(store experimentstore.Storer, enabledConfigTypes []string, criterias []TerminationCriteria) Terminator {
	return Terminator{store: store, enabledConfigTypes: enabledConfigTypes, criterias: criterias}
}

func (t Terminator) Run() {
	go func() {
		trackedExperiments := map[uint64]struct{}{}
		ticker := time.NewTicker(time.Second)

		for range ticker.C {
			es, err := t.store.GetExperiments(context.Background(), t.enabledConfigTypes, experimentationv1.GetExperimentsRequest_STATUS_RUNNING)
			if err != nil {
				// log
			}

			for _, e := range es {
				if _, ok := trackedExperiments[e.Id]; !ok {
					trackedExperiments[e.Id] = struct{}{}
					go func() {
						ticker := time.NewTicker(time.Second)

						for range ticker.C {
							for _, c := range t.criterias {
								tt, _ := ptypes.Timestamp(e.StartTime)
								err := c.ShouldTerminate(tt, e)
								if err != nil {
									fmt.Println("TERMINAING")
									t.store.TerminateExperiment(context.Background(), e.Id, err.Error())
									return
								}
							}
						}
					}()
				}
			}
		}
	}()
}
