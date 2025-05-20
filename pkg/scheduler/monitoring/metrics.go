/*
Copyright 2024 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package monitoring // Package monitoring provides the indicator definitions and recording methods required by Dapr Scheduler.

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/dapr/dapr/pkg/diagnostics/utils"
)

// Job status codes used in metric recording.
const (
	JobStatusSuccess     int64 = 1 // Job completed successfully
	JobStatusFailed      int64 = 2 // Job execution failed
	JobStatusUndelivered int64 = 3 // Job could not be delivered
)

var (
	sidecarConnectionCount int64

	sidecarsConnectedGauge = stats.Int64(
		"scheduler/sidecars_connected",
		"The number of dapr sidecars actively connected to the scheduler service.",
		stats.UnitDimensionless)

	jobsScheduledTotal = stats.Int64(
		"scheduler/jobs_created_total",
		"The total number of jobs scheduled.",
		stats.UnitDimensionless)

	jobsTriggeredTotal = stats.Int64(
		"scheduler/jobs_triggered_total",
		"The total number of successfully triggered jobs.",
		stats.UnitDimensionless)

	triggerLatency = stats.Float64(
		"scheduler/trigger_latency",
		"The total time it takes to trigger a job from the scheduler service.",
		stats.UnitMilliseconds)

	jobsStatusSuccessTotal = stats.Int64(
		"scheduler/jobs_status_success_total",
		"The total number of jobs with status SUCCESS.",
		stats.UnitDimensionless)

	jobsStatusFailedTotal = stats.Int64(
		"scheduler/jobs_status_failed_total",
		"The total number of jobs with status FAILED.",
		stats.UnitDimensionless)

	jobsStatusUndeliveredTotal = stats.Int64(
		"scheduler/jobs_status_undelivered_total",
		"The total number of jobs with status UNDELIVERED.",
		stats.UnitDimensionless)
)

// RecordSidecarsConnectedCount records the number of dapr sidecars connected to the scheduler service
func RecordSidecarsConnectedCount(change int) {
	current := atomic.AddInt64(&sidecarConnectionCount, int64(change))
	if err := stats.RecordWithTags(context.Background(), utils.WithTags(sidecarsConnectedGauge.Name()), sidecarsConnectedGauge.M(current)); err != nil {
		log.Printf("failed to record sidecars connected gauge: %v", err)
	}
}

// RecordJobsScheduledCount records the number of jobs scheduled (jobType: "job", "actor", or "unknown")
func RecordJobsScheduledCount(jobType string) {
	if err := stats.RecordWithTags(context.Background(), utils.WithTags(jobsScheduledTotal.Name(), jobType), jobsScheduledTotal.M(1)); err != nil {
		log.Printf("failed to record jobs scheduled count: %v", err)
	}
}

// RecordJobsTriggeredCount records the number of successfully triggered jobs
func RecordJobsTriggeredCount(jobType string) {
	if err := stats.RecordWithTags(context.Background(), utils.WithTags(jobsTriggeredTotal.Name(), jobType), jobsTriggeredTotal.M(1)); err != nil {
		log.Printf("failed to record jobs triggered count: %v", err)
	}
}

// RecordJobStatus records job status metrics using int values: 1=SUCCESS, 2=FAILED, 3=UNDELIVERED
func RecordJobStatus(status int64) {
	switch status {
	case JobStatusSuccess:
		stats.Record(context.Background(), jobsStatusSuccessTotal.M(1))
	case JobStatusFailed:
		stats.Record(context.Background(), jobsStatusFailedTotal.M(1))
	case JobStatusUndelivered:
		stats.Record(context.Background(), jobsStatusUndeliveredTotal.M(1))
	}
}

// RecordTriggerDuration records the time it takes to send the job to dapr from the scheduler service
func RecordTriggerDuration(start time.Time) {
	elapsed := time.Since(start).Milliseconds()
	if err := stats.RecordWithTags(context.Background(), utils.WithTags(triggerLatency.Name()), triggerLatency.M(float64(elapsed))); err != nil {
		log.Printf("failed to record trigger duration: %v", err)
	}
}

// InitMetrics initializes the scheduler service metrics
func InitMetrics() error {
	err := view.Register(
		utils.NewMeasureView(sidecarsConnectedGauge, []tag.Key{}, view.LastValue()),
		utils.NewMeasureView(jobsScheduledTotal, []tag.Key{}, view.Count()),
		utils.NewMeasureView(jobsTriggeredTotal, []tag.Key{}, view.Count()),
		utils.NewMeasureView(triggerLatency, []tag.Key{}, view.Distribution(0, 100, 500, 1000, 5000, 10000)),
		utils.NewMeasureView(jobsStatusSuccessTotal, []tag.Key{}, view.Count()),
		utils.NewMeasureView(jobsStatusFailedTotal, []tag.Key{}, view.Count()),
		utils.NewMeasureView(jobsStatusUndeliveredTotal, []tag.Key{}, view.Count()),
	)
	return err
}
