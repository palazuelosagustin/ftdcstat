package discovery

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"mongodb-ftdcstat/internal/model"
)

func TestDiscoverOrdersRotatedMetricsAndInterimLast(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{
		"metrics.interim",
		"metrics.2026-06-02T00-00-00Z-00000",
		"metrics.2026-06-01T00-00-00Z-00000",
		"export.json",
		"notes.txt",
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("{}"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	files, warnings, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %#v", warnings)
	}
	got := make([]string, len(files))
	for i, file := range files {
		got[i] = filepath.Base(file.Path)
	}
	want := []string{
		"metrics.2026-06-01T00-00-00Z-00000",
		"metrics.2026-06-02T00-00-00Z-00000",
		"export.json",
		"metrics.interim",
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("order mismatch at %d: got %q want %q (all=%v)", i, got[i], want[i], got)
		}
	}
}

func TestDiscoverRejectsNonDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "metrics.2026-06-01T00-00-00Z-00000")
	if err := os.WriteFile(path, []byte{}, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Discover(path); err == nil {
		t.Fatal("expected non-directory error")
	}
}

func TestFilterByTimeRangeSkipsNonOverlappingRotatedFiles(t *testing.T) {
	files := []MetricFile{
		{Path: "metrics.2026-06-01T00-00-00Z-00000", Kind: KindMetrics, Timestamp: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)},
		{Path: "metrics.2026-06-01T01-00-00Z-00000", Kind: KindMetrics, Timestamp: time.Date(2026, 6, 1, 1, 0, 0, 0, time.UTC)},
		{Path: "metrics.2026-06-01T02-00-00Z-00000", Kind: KindMetrics, Timestamp: time.Date(2026, 6, 1, 2, 0, 0, 0, time.UTC)},
		{Path: "metrics.2026-06-01T03-00-00Z-00000", Kind: KindMetrics, Timestamp: time.Date(2026, 6, 1, 3, 0, 0, 0, time.UTC)},
	}
	filtered := FilterByTimeRange(files, model.TimeRange{
		From: time.Date(2026, 6, 1, 1, 30, 0, 0, time.UTC),
		To:   time.Date(2026, 6, 1, 1, 45, 0, 0, time.UTC),
	})
	if len(filtered) != 2 {
		t.Fatalf("got %d files", len(filtered))
	}
	if filepath.Base(filtered[0].Path) != "metrics.2026-06-01T01-00-00Z-00000" {
		t.Fatalf("wrong file kept: %s", filtered[0].Path)
	}
	if filepath.Base(filtered[1].Path) != "metrics.2026-06-01T03-00-00Z-00000" {
		t.Fatalf("unknown-tail file should be kept conservatively: %s", filtered[1].Path)
	}
}
