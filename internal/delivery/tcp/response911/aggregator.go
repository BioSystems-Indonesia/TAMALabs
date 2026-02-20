package response911

import (
	"context"
	"sync"
	"time"

	"log/slog"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
)

// aggregator collects partial ORU_R01 messages for the same barcode and
// flushes them as a single ORU_R01 after a debounce timeout.
type aggregator struct {
	mu      sync.Mutex
	items   map[string]*aggItem
	timeout time.Duration
	// process is a callback to deliver aggregated ORU_R01
	process func(context.Context, entity.ORU_R01) error
}

type aggItem struct {
	barcode string
	// map test_code -> ObservationResult (last-wins)
	results map[string]entity.ObservationResult
	timer   *time.Timer
}

func newAggregator(process func(context.Context, entity.ORU_R01) error, timeout time.Duration) *aggregator {
	if timeout <= 0 {
		timeout = 1 * time.Second
	}
	return &aggregator{
		items:   make(map[string]*aggItem),
		timeout: timeout,
		process: process,
	}
}

// Add merges results from an ORU_R01 into the aggregator for its barcode.
// If barcode is empty, the ORU_R01 is flushed immediately.
func (a *aggregator) Add(oru *entity.ORU_R01) {
	if oru == nil || len(oru.Patient) == 0 || len(oru.Patient[0].Specimen) == 0 {
		return
	}

	spec := oru.Patient[0].Specimen[0]
	barcode := spec.Barcode
	if barcode == "" {
		// Nothing we can aggregate reliably — forward immediately
		_ = a.process(context.Background(), *oru)
		return
	}

	a.mu.Lock()
	it, ok := a.items[barcode]
	if !ok {
		it = &aggItem{barcode: barcode, results: make(map[string]entity.ObservationResult)}
		a.items[barcode] = it
	}

	// merge results (last-wins)
	for _, r := range spec.ObservationResult {
		if r.TestCode == "" {
			continue
		}
		it.results[r.TestCode] = r
	}

	// reset timer
	if it.timer != nil {
		if !it.timer.Stop() {
			// drain if needed
			select {
			case <-it.timer.C:
			default:
			}
		}
	}
	it.timer = time.AfterFunc(a.timeout, func() {
		a.flush(barcode)
	})
	a.mu.Unlock()
}

func (a *aggregator) flush(barcode string) {
	a.mu.Lock()
	it, ok := a.items[barcode]
	if !ok {
		a.mu.Unlock()
		return
	}
	delete(a.items, barcode)
	a.mu.Unlock()

	// build ORU from aggregated results
	var obs []entity.ObservationResult
	for _, r := range it.results {
		obs = append(obs, r)
	}

	if len(obs) == 0 {
		slog.Warn("aggregator flush: no observations to flush", "barcode", barcode)
		return
	}

	spec := entity.Specimen{
		ReceivedDate:      time.Now(),
		Barcode:           barcode,
		ObservationResult: obs,
	}

	oru := &entity.ORU_R01{
		MSH: entity.MSH{},
		Patient: []entity.Patient{{
			Specimen: []entity.Specimen{spec},
		}},
	}

	if err := a.process(context.Background(), *oru); err != nil {
		slog.Error("aggregator: ProcessORUR01 failed", "error", err, "barcode", barcode)
	}
}
