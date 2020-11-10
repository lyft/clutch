package stats

import (
	"context"
	"runtime"
	"time"

	"github.com/uber-go/tally"
)

type RuntimeStatCollector struct {
	scope tally.Scope
}

func NewRuntimeStats(scope tally.Scope) *RuntimeStatCollector {
	return &RuntimeStatCollector{
		scope: scope.SubScope("runtime"),
	}
}

func (r *RuntimeStatCollector) Collect(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 10)
	for {
		r.collectCPUStats()
		r.collectMemStats()
		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (r *RuntimeStatCollector) collectCPUStats() {
	r.scope.Gauge("cpu.goroutines").Update(float64(runtime.NumGoroutine()))
	r.scope.Gauge("cpu.cgo_calls").Update(float64(runtime.NumCgoCall()))
}

func (r *RuntimeStatCollector) collectMemStats() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// General Stats
	r.scope.Gauge("memory.alloc").Update(float64(memStats.Alloc))
	r.scope.Gauge("memory.total").Update(float64(memStats.TotalAlloc))
	r.scope.Gauge("memory.sys").Update(float64(memStats.Sys))
	r.scope.Gauge("memory.lookups").Update(float64(memStats.Lookups))
	r.scope.Gauge("memory.malloc").Update(float64(memStats.Mallocs))
	r.scope.Gauge("memory.frees").Update(float64(memStats.Frees))

	// Heap Stats
	r.scope.Gauge("memory.heap.alloc").Update(float64(memStats.HeapAlloc))
	r.scope.Gauge("memory.heap.sys").Update(float64(memStats.HeapSys))
	r.scope.Gauge("memory.heap.idle").Update(float64(memStats.HeapIdle))
	r.scope.Gauge("memory.heap.inuse").Update(float64(memStats.HeapInuse))
	r.scope.Gauge("memory.heap.released").Update(float64(memStats.HeapReleased))
	r.scope.Gauge("memory.heap.objects").Update(float64(memStats.HeapObjects))

	// Stack Stats
	r.scope.Gauge("memory.stack.inuse").Update(float64(memStats.StackInuse))
	r.scope.Gauge("memory.stack.sys").Update(float64(memStats.StackSys))
	r.scope.Gauge("memory.stack.mspan_inuse").Update(float64(memStats.MSpanInuse))
	r.scope.Gauge("memory.stack.mspan_sys").Update(float64(memStats.MSpanSys))
	r.scope.Gauge("memory.stack.mcache_inuse").Update(float64(memStats.MCacheInuse))
	r.scope.Gauge("memory.stack.mcache_sys").Update(float64(memStats.MCacheSys))

	r.scope.Gauge("memory.othersys").Update(float64(memStats.OtherSys))

	// GC Stats
	r.scope.Gauge("memory.gc.sys").Update(float64(memStats.GCSys))
	r.scope.Gauge("memory.gc.next").Update(float64(memStats.NextGC))
	r.scope.Gauge("memory.gc.last").Update(float64(memStats.LastGC))
	r.scope.Gauge("memory.gc.pause_total").Update(float64(memStats.PauseTotalNs))
	r.scope.Gauge("memory.gc.count").Update(float64(memStats.NumGC))
}
