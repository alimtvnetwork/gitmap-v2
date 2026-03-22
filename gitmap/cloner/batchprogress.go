package cloner

import (
	"fmt"
	"os"
	"time"

	"github.com/user/gitmap/constants"
)

// BatchProgress tracks progress for any batch operation (pull, exec, status).
type BatchProgress struct {
	total     int
	current   int
	start     time.Time
	quiet     bool
	succeeded int
	failed    int
	skipped   int
	operation string
}

// NewBatchProgress creates a progress tracker for a named operation.
func NewBatchProgress(total int, operation string, quiet bool) *BatchProgress {
	return &BatchProgress{
		total:     total,
		start:     time.Now(),
		quiet:     quiet,
		operation: operation,
	}
}

// BeginItem prints progress for starting an item.
func (p *BatchProgress) BeginItem(name string) {
	p.current++
	if p.quiet {
		return
	}

	fmt.Fprintf(os.Stderr, constants.BatchProgressBeginFmt, p.current, p.total, name)
}

// Succeed marks an item as successful.
func (p *BatchProgress) Succeed() {
	p.succeeded++
	if p.quiet {
		return
	}

	fmt.Fprintf(os.Stderr, constants.BatchProgressDoneFmt, formatDuration(time.Since(p.start)))
}

// Fail marks an item as failed.
func (p *BatchProgress) Fail() {
	p.failed++
	if p.quiet {
		return
	}

	fmt.Fprintf(os.Stderr, constants.BatchProgressFailFmt)
}

// Skip marks an item as skipped (e.g., missing directory).
func (p *BatchProgress) Skip() {
	p.skipped++
	if p.quiet {
		return
	}

	fmt.Fprintf(os.Stderr, constants.BatchProgressSkipFmt)
}

// PrintSummary prints the final summary.
func (p *BatchProgress) PrintSummary() {
	if p.quiet {
		return
	}

	elapsed := formatDuration(time.Since(p.start))
	fmt.Fprintf(os.Stderr, constants.BatchProgressSummaryFmt,
		p.operation, p.current, p.total, elapsed)
	fmt.Fprintf(os.Stderr, constants.BatchProgressDetailFmt,
		p.succeeded, p.failed, p.skipped)
}

// Succeeded returns the success count.
func (p *BatchProgress) Succeeded() int { return p.succeeded }

// Failed returns the failure count.
func (p *BatchProgress) Failed() int { return p.failed }

// Skipped returns the skip count.
func (p *BatchProgress) Skipped() int { return p.skipped }
