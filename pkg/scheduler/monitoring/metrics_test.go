package monitoring

import (
	"testing"

	"go.opencensus.io/stats/view"
)

// TestInitMetrics verifies that all expected OpenCensus views
// are successfully registered without errors.
func TestInitMetrics(t *testing.T) {
	err := InitMetrics()
	if err != nil {
		t.Fatalf("InitMetrics failed: %v", err)
	}

	expectedViews := []string{
		"scheduler/jobs_status_success_total",
		"scheduler/jobs_status_failed_total",
		"scheduler/jobs_status_undelivered_total",
	}

	for _, viewName := range expectedViews {
		if view.Find(viewName) == nil {
			t.Errorf("Expected view %s not found", viewName)
		}
	}
}

// TestRecordJobStatus verifies that job status metrics are recorded
// correctly for each defined status, and gracefully handles unknown status.
func TestRecordJobStatus(t *testing.T) {
	// Step 1: Initialize metrics
	_ = InitMetrics()

	// Step 2: Unregister previous views to clear test state
	view.Unregister(
		view.Find("scheduler/jobs_status_success_total"),
		view.Find("scheduler/jobs_status_failed_total"),
		view.Find("scheduler/jobs_status_undelivered_total"),
	)
	_ = InitMetrics()

	// Step 3: Simulate metric recording for each status
	RecordJobStatus(JobStatusSuccess)
	RecordJobStatus(JobStatusFailed)
	RecordJobStatus(JobStatusUndelivered)

	// Step 4: Call with an unknown status (should log only, no panic)
	RecordJobStatus(999)

	// Step 5: Validate the recorded metric values
	testCases := []struct {
		metricName string
		expected  int64
	}{
		{"scheduler/jobs_status_success_total", 1},
		{"scheduler/jobs_status_failed_total", 1},
		{"scheduler/jobs_status_undelivered_total", 1},
	}

	for _, tc := range testCases {
		rows, err := view.RetrieveData(tc.metricName)
		if err != nil {
			t.Errorf("Failed to retrieve view data for %s: %v", tc.metricName, err)
			continue
		}
		if len(rows) == 0 {
			t.Errorf("No data recorded for %s", tc.metricName)
			continue
		}
		got := rows[0].Data.(*view.CountData).Value
		if got < tc.expected {
			t.Errorf("Expected at least %v for %s, got %v", tc.expected, tc.metricName, got)
		}
	}
}