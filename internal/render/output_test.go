package render

import "testing"

func TestOutputModeFor(t *testing.T) {
	if got := OutputModeFor(Options{View: "summary"}); got != OutputTable {
		t.Fatalf("summary table mode=%v", got)
	}
	if got := OutputModeFor(Options{View: "summary", JSON: true}); got != OutputJSON {
		t.Fatalf("summary json mode=%v", got)
	}
}

func TestNeedsBufferedRowsOnlyForJSON(t *testing.T) {
	for _, view := range []string{"summary", "server", "wt", "system", "network", "repl"} {
		if NeedsBufferedRows(Options{View: view}) {
			t.Fatalf("view %s table output should stream rows", view)
		}
		if !NeedsBufferedRows(Options{View: view, JSON: true}) {
			t.Fatalf("view %s json output should buffer rows", view)
		}
	}
}
