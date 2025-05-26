// main.go - Demonstrates recording scheduler job metrics
package main

import (
	"fmt"
	"time"

	"github.com/dapr/dapr/pkg/scheduler/monitoring"
)

func main() {
	fmt.Println("== Dapr Scheduler Metrics Demo ==")

	// Step 1: Initialize metrics system
	if err := monitoring.InitMetrics(); err != nil {
		fmt.Printf("Failed to initialize metrics: %v\n", err)
		return
	}
	fmt.Println("Metrics initialized successfully.")

	// Step 2: Record valid job status metrics
	fmt.Println("Recording job status metrics...")
	monitoring.RecordJobStatus(monitoring.JobStatusSuccess)
	monitoring.RecordJobStatus(monitoring.JobStatusFailed)
	monitoring.RecordJobStatus(monitoring.JobStatusUndelivered)

	// Step 3: Record invalid job status (should trigger log, not panic)
	fmt.Println("Recording unknown job status (expected: log only)...")
	monitoring.RecordJobStatus(999)

	// Step 4: Simulate trigger latency recording (1500ms delay)
	fmt.Println("Recording job trigger duration...")
	startTime := time.Now().Add(-1500 * time.Millisecond)
	monitoring.RecordTriggerDuration(startTime)

	fmt.Println("All metrics recorded successfully.")
}