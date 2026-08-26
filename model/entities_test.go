package model

import (
	"testing"
	"time"
)

func TestRecordValidation(t *testing.T) {
	r := NewRecord("1", "title", "body", time.Now().Add(time.Hour))
	if err := ValidateRecord(r); err != nil {
		t.Fatal(err)
	}
}
