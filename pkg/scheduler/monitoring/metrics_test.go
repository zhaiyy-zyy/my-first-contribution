package monitoring

import (
	"testing"

	"go.opencensus.io/stats/view"
)

func TestInitMetrics(t *testing.T) {
	err := InitMetrics()
	if err != nil {
		t.Errorf("InitMetrics failed: %v", err)
	}

	views := []*view.View{
		view.Find("scheduler/jobs_status_success_total"),
		view.Find("scheduler/jobs_status_failed_total"),
		view.Find("scheduler/jobs_status_undelivered_total"),
	}

	for _, v := range views {
		if v == nil {
			t.Errorf("Expected view not found")
		}
	}
}

func TestRecordJobStatus(t *testing.T) {
	t.Run("SUCCESS", func(t *testing.T) {
		RecordJobStatus(1)
	})
	t.Run("FAILED", func(t *testing.T) {
		RecordJobStatus(2)
	})
	t.Run("UNDELIVERED", func(t *testing.T) {
		RecordJobStatus(3)
	})
	t.Run("UNKNOWN", func(t *testing.T) {
		RecordJobStatus(999)
	})
}