package discovery

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"mongodb-ftdcstat/internal/model"
)

type FileKind string

const (
	KindMetrics FileKind = "metrics"
	KindInterim FileKind = "interim"
	KindJSON    FileKind = "json"
)

type MetricFile struct {
	Path      string
	Kind      FileKind
	Timestamp time.Time
	Sequence  int
}

var metricsTimestampRE = regexp.MustCompile(`^metrics\.(\d{4}-\d{2}-\d{2}T\d{2}-\d{2}-\d{2}Z)`)

func Discover(root string) ([]MetricFile, []model.Warning, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, nil, err
	}
	if !info.IsDir() {
		return nil, nil, errors.New("path must be a directory")
	}

	var files []MetricFile
	var warnings []model.Warning
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			warnings = append(warnings, model.Warning{Source: path, Message: walkErr.Error()})
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		kind, ok := classify(entry.Name())
		if !ok {
			return nil
		}
		ts, tsOK := parseMetricTimestamp(entry.Name())
		if !tsOK && kind == KindMetrics && strings.HasPrefix(entry.Name(), "metrics.") && !strings.Contains(entry.Name(), "interim") {
			warnings = append(warnings, model.Warning{Source: path, Message: "metrics file name does not contain a parseable rotation timestamp"})
		}
		files = append(files, MetricFile{Path: path, Kind: kind, Timestamp: ts})
		return nil
	})
	if err != nil {
		return nil, warnings, err
	}

	sort.SliceStable(files, func(i, j int) bool {
		a, b := files[i], files[j]
		ap, bp := sortPriority(a), sortPriority(b)
		if ap != bp {
			return ap < bp
		}
		if !a.Timestamp.IsZero() && !b.Timestamp.IsZero() && !a.Timestamp.Equal(b.Timestamp) {
			return a.Timestamp.Before(b.Timestamp)
		}
		return filepath.Base(a.Path) < filepath.Base(b.Path)
	})
	for i := range files {
		files[i].Sequence = i
	}
	if len(files) == 0 {
		return nil, warnings, errors.New("no FTDC metrics or exported JSON files found")
	}
	return files, warnings, nil
}

func FilterByTimeRange(files []MetricFile, tr model.TimeRange) []MetricFile {
	if tr.IsZero() || len(files) == 0 {
		return files
	}
	out := make([]MetricFile, 0, len(files))
	for i, file := range files {
		if file.Timestamp.IsZero() || file.Kind != KindMetrics {
			out = append(out, file)
			continue
		}
		next := nextTimestamp(files, i)
		if next.IsZero() {
			out = append(out, file)
			continue
		}
		if tr.Overlaps(file.Timestamp, next) {
			out = append(out, file)
		}
	}
	for i := range out {
		out[i].Sequence = i
	}
	return out
}

func nextTimestamp(files []MetricFile, start int) time.Time {
	for i := start + 1; i < len(files); i++ {
		if files[i].Kind == KindMetrics && !files[i].Timestamp.IsZero() {
			return files[i].Timestamp
		}
	}
	return time.Time{}
}

func classify(name string) (FileKind, bool) {
	lower := strings.ToLower(name)
	switch {
	case strings.HasPrefix(name, "metrics.interim"), strings.HasPrefix(name, "interim"), strings.Contains(name, ".interim"):
		return KindInterim, true
	case strings.HasPrefix(name, "metrics."):
		return KindMetrics, true
	case strings.HasSuffix(lower, ".json"), strings.HasSuffix(lower, ".jsonl"), strings.HasSuffix(lower, ".ndjson"):
		return KindJSON, true
	default:
		return "", false
	}
}

func parseMetricTimestamp(name string) (time.Time, bool) {
	matches := metricsTimestampRE.FindStringSubmatch(name)
	if len(matches) != 2 {
		return time.Time{}, false
	}
	t, err := time.Parse("2006-01-02T15-04-05Z", matches[1])
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

func sortPriority(file MetricFile) int {
	switch file.Kind {
	case KindMetrics:
		if file.Timestamp.IsZero() {
			return 1
		}
		return 0
	case KindJSON:
		return 2
	case KindInterim:
		return 3
	default:
		return 4
	}
}
