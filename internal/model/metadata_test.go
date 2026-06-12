package model

import (
	"testing"
	"time"
)

func TestMetadataSelectsNewestAndSuppressesChangeWarnings(t *testing.T) {
	m := NewMetadata()
	oldTime := time.Unix(10, 0)
	newTime := time.Unix(20, 0)
	m.AddDocument(oldTime, "old", map[string]any{
		"buildInfo": map[string]any{"version": "8.0.1", "gitVersion": "aaa"},
	})
	m.AddDocument(newTime, "new", map[string]any{
		"buildInfo": map[string]any{"version": "8.0.2", "gitVersion": "bbb"},
	})
	doc, ok := m.LatestDoc("buildInfo")
	if !ok {
		t.Fatal("missing buildInfo")
	}
	if got, _ := Lookup(doc, "version"); got != "8.0.2" {
		t.Fatalf("got version %v", got)
	}
	if len(m.Warnings) != 0 {
		t.Fatalf("metadata-change warnings should be suppressed: %#v", m.Warnings)
	}
}
