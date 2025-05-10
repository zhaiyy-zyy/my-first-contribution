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

package monitoring

import (
	"context"
	"sync/atomic"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/dapr/dapr/pkg/diagnostics/utils"
)

var (
	sidecarConnectionCount int64 = 0

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
	stats.RecordWithTags(context.Background(), utils.WithTags(sidecarsConnectedGauge.Name()), sidecarsConnectedGauge.M(current))
}

// RecordJobsScheduledCount records the number of jobs scheduled (jobType: "job", "actor", or "unknown")
func RecordJobsScheduledCount(jobType string) {
	stats.RecordWithTags(context.Background(), utils.WithTags(jobsScheduledTotal.Name(), jobType), jobsScheduledTotal.M(1))
}

// RecordJobsTriggeredCount records the number of successfully triggered jobs
func RecordJobsTriggeredCount(jobType string) {
	stats.RecordWithTags(context.Background(), utils.WithTags(jobsTriggeredTotal.Name(), jobType), jobsTriggeredTotal.M(1))
}

// RecordJobStatus records job status metrics using int values: 1=SUCCESS, 2=FAILED, 3=UNDELIVERED
func RecordJobStatus(status int64) {
	switch status {
	case 1:
		stats.Record(context.Background(), jobsStatusSuccessTotal.M(1))
	case 2:
		stats.Record(context.Background(), jobsStatusFailedTotal.M(1))
	case 3:
		stats.Record(context.Background(), jobsStatusUndeliveredTotal.M(1))
	}
}

// RecordTriggerDuration records the time it takes to send the job to dapr from the scheduler service
func RecordTriggerDuration(start time.Time) {
	elapsed := time.Since(start).Milliseconds()
	stats.RecordWithTags(context.Background(), utils.WithTags(triggerLatency.Name()), triggerLatency.M(float64(elapsed)))
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