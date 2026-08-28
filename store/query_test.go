package store

import (
	"drupal-scheduler/model"
	"testing"
	"time"
)

func TestFilterRecords(t *testing.T) {
	rows := []model.Record{{ID: "1", Title: "Alpha", Status: "queued", StartAt: time.Now()}, {ID: "2", Title: "Beta", Status: "published"}}
	if len(FilterRecords(rows, "queued", "alp")) != 1 {
		t.Fatal("filter")
	}
}
