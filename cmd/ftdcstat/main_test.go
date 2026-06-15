package main

import "testing"

func TestParseArgsDefaultIntervalIsSixty(t *testing.T) {
	opts, err := parseArgs([]string{"diagnostic.data"})
	if err != nil {
		t.Fatal(err)
	}
	if opts.Interval != 60 {
		t.Fatalf("interval=%d", opts.Interval)
	}
	if opts.View != "summary" {
		t.Fatalf("view=%s", opts.View)
	}
}

func TestParseArgsSummaryViewAccepted(t *testing.T) {
	opts, err := parseArgs([]string{"diagnostic.data", "--view", "summary"})
	if err != nil {
		t.Fatal(err)
	}
	if opts.View != "summary" {
		t.Fatalf("view=%s", opts.View)
	}
}

func TestParseArgsDiskAliasesToSystem(t *testing.T) {
	opts, err := parseArgs([]string{"diagnostic.data", "--view", "disk"})
	if err != nil {
		t.Fatal(err)
	}
	if opts.View != "system" {
		t.Fatalf("view=%s", opts.View)
	}
}

func TestParseArgsVerbose(t *testing.T) {
	opts, err := parseArgs([]string{"diagnostic.data", "--verbose"})
	if err != nil {
		t.Fatal(err)
	}
	if !opts.Verbose {
		t.Fatal("expected verbose=true")
	}
}

func TestParseArgsFromTo(t *testing.T) {
	opts, err := parseArgs([]string{"diagnostic.data", "--from", "2026-06-04T19:00:00", "--to", "2026-06-04T20:00:00"})
	if err != nil {
		t.Fatal(err)
	}
	if opts.Range.From.IsZero() || opts.Range.To.IsZero() {
		t.Fatalf("range not set: %#v", opts.Range)
	}
}
