package monitoring

import (
	"testing"

	"go.opencensus.io/stats/view"
)

func TestInitMetrics(t *testing.T) {
	err := InitMetrics()
	if err != nil {
		t.Fatalf("InitMetrics failed: %v", err)
	}

	metrics := []string{
		"scheduler/jobs_status_success_total",
		"scheduler/jobs_status_failed_total",
		"scheduler/jobs_status_undelivered_total",
	}

	for _, name := range metrics {
		v := view.Find(name)
		if v == nil {
			t.Errorf("Expected view %s not found", name)
		}
	}
}

func TestRecordJobStatus(t *testing.T) {
	// Initialization
	_ = InitMetrics()

	// Clean up old data
	view.Unregister(
		view.Find("scheduler/jobs_status_success_total"),
		view.Find("scheduler/jobs_status_failed_total"),
		view.Find("scheduler/jobs_status_undelivered_total"),
	)
	_ = InitMetrics()

	// Call each state
	RecordJobStatus(JobStatusSuccess)
	RecordJobStatus(JobStatusFailed)
	RecordJobStatus(JobStatusUndelivered)
	RecordJobStatus(999)

	tests := []struct {
		metricName string
		expected   int64
	}{
		{"scheduler/jobs_status_success_total", 1},
		{"scheduler/jobs_status_failed_total", 1},
		{"scheduler/jobs_status_undelivered_total", 1},
	}

	for _, test := range tests {
		rows, err := view.RetrieveData(test.metricName)
		if err != nil {
			t.Errorf("Failed to retrieve view data for %s: %v", test.metricName, err)
			continue
		}
		if len(rows) == 0 {
			t.Errorf("No data recorded for %s", test.metricName)
			continue
		}
		got := rows[0].Data.(*view.CountData).Value
		if int64(got) < test.expected {
			t.Errorf("Expected at least %d for %s, got %d", test.expected, test.metricName, got)
		}
	}
}
