package reporting

import (
	"drupal-scheduler/model"
	"testing"
)

func TestSummarize(t *testing.T) {
	s := Summarize([]model.Record{{Status: "queued"}, {Status: "published"}})
	if s.Total != 2 || s.Queued != 1 {
		t.Fatal(s)
	}
}
