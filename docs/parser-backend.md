# Parser backend

`mongodb-ftdcstat` keeps FTDC parsing behind `internal/ftdc.SampleReader`.
The default implementation is `NativeReader`, which decodes MongoDB diagnostic
files directly from concatenated BSON records and metric chunks.

`github.com/mongodb/ftdc` was checked as the requested candidate backend. The
current public module resolves to a pseudo-version that declares Go 1.24 and
pulls in a broad Evergreen/system metrics dependency tree, while this project is
built with Go 1.22. For this first version the dependency is therefore not wired
into the CLI. The rest of the tool depends only on `SampleReader`, `MetricSample`,
`Metadata`, and `Warning`, so a future backend can replace `NativeReader` without
changing metric formulas or renderers.

The native decoder implements:

- outer concatenated BSON record reading
- metadata record extraction
- `type: 1` metric chunk payload decoding
- zlib decompression after the four-byte uncompressed-size header
- BSON reference document parsing
- metric-count/delta-count handling, including reversed-count fallback
- signed varint delta decoding with FTDC zero-run compression
- per-sample reconstruction of selected metric paths
