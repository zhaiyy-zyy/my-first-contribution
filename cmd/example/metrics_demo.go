package main

import (
	"fmt"
	"time"

	"github.com/dapr/dapr/pkg/scheduler/monitoring"
)

func main() {
	fmt.Println("Initializing metrics...")

	err := monitoring.InitMetrics()
	if err != nil {
		fmt.Printf("Failed to init metrics: %v\n", err)
		return
	}

	fmt.Println("Recording metrics...")
	monitoring.RecordJobStatus(monitoring.JobStatusSuccess)
	monitoring.RecordJobStatus(monitoring.JobStatusFailed)
	monitoring.RecordJobStatus(monitoring.JobStatusUndelivered)
	monitoring.RecordJobStatus(999)

	monitoring.RecordTriggerDuration(time.Now().Add(-1500 * time.Millisecond))

	fmt.Println("Metrics recorded.")
}
