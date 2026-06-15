package derive

import "testing"

func TestViewNeedsVerboseReplication(t *testing.T) {
	cases := []struct {
		view    string
		verbose bool
		want    bool
	}{
		{"repl", true, true},
		{"all", true, true},
		{"server", true, false},
		{"system", true, false},
		{"repl", false, false},
		{"all", false, false},
	}
	for _, tc := range cases {
		if got := ViewNeedsVerboseReplication(tc.view, tc.verbose); got != tc.want {
			t.Fatalf("ViewNeedsVerboseReplication(%q, %v)=%v want %v", tc.view, tc.verbose, got, tc.want)
		}
	}
}

func TestRequiredPathsForVerboseReplication(t *testing.T) {
	paths, _ := RequiredPathsFor("repl", true)
	for _, path := range verboseReplicationPaths {
		if !paths[path] {
			t.Fatalf("expected verbose path %q", path)
		}
	}
	plain, _ := RequiredPathsFor("all", false)
	for _, path := range verboseReplicationPaths {
		if plain[path] {
			t.Fatalf("non-verbose should not include %q", path)
		}
	}
	systemVerbose, _ := RequiredPathsFor("system", true)
	for _, path := range verboseReplicationPaths {
		if systemVerbose[path] {
			t.Fatalf("system verbose should not include %q", path)
		}
	}
}

func TestInterestingVerboseReplicationPaths(t *testing.T) {
	paths, prefixes := RequiredPathsFor("repl", true)
	if !Interesting("replSetGetStatus.members.0.pingMs", paths, prefixes, true) {
		t.Fatal("expected pingMs to be interesting with verbose replication")
	}
	if Interesting("replSetGetStatus.members.0.pingMs", paths, prefixes, false) {
		t.Fatal("pingMs should not be interesting without verbose replication")
	}
	if !Interesting("serverStatus.metrics.repl.apply.ops", paths, prefixes, true) {
		t.Fatal("expected serverStatus.metrics.repl.apply.ops to be interesting")
	}
	if Interesting("serverStatus.metrics.repl.buffer.count", paths, prefixes, true) {
		t.Fatal("broad metrics.repl path should not be interesting")
	}
	if Interesting("replSetGetStatus.set", paths, prefixes, true) {
		t.Fatal("broad replSetGetStatus path should not be interesting")
	}
}
