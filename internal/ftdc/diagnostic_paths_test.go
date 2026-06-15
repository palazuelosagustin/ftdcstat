package ftdc_test

import (
	"os"
	"strings"
	"testing"

	"ftdcstat/internal/discovery"
	"ftdcstat/internal/ftdc"
)

func TestDiagnosticDataVerboseReplicationPathsPresent(t *testing.T) {
	const dir = "../../diagnostic.data"
	if _, err := os.Stat(dir); err != nil {
		t.Skip("diagnostic.data not available")
	}
	files, _, err := discovery.Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	var file discovery.MetricFile
	for _, f := range files {
		if strings.Contains(f.Path, "metrics.") {
			file = f
			break
		}
	}
	opts := ftdc.ReaderOptionsFor("repl", true)
	capture, err := ftdc.NewNativeReader().ReadFiles([]discovery.MetricFile{file}, opts)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"serverStatus.metrics.repl.apply.ops",
		"serverStatus.metrics.repl.buffer.apply.count",
		"serverStatus.metrics.repl.buffer.apply.sizeBytes",
	}
	for _, path := range want {
		found := false
		for _, sample := range capture.Samples {
			if _, ok := sample.Values[path]; ok {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected path %q in FTDC samples", path)
		}
	}
}
