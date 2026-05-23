# logslice

A tool to extract and filter structured log ranges from large files by timestamp or pattern.

---

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git && cd logslice && go build -o logslice .
```

---

## Usage

```
logslice [flags] <logfile>
```

### Examples

Extract logs between two timestamps:
```bash
logslice --from "2024-01-15T08:00:00" --to "2024-01-15T09:00:00" app.log
```

Filter by pattern within a time range:
```bash
logslice --from "2024-01-15T08:00:00" --pattern "ERROR" app.log
```

Read from stdin:
```bash
cat app.log | logslice --from "2024-01-15T08:00:00" --to "2024-01-15T09:00:00"
```

### Flags

| Flag | Description |
|------|-------------|
| `--from` | Start timestamp (RFC3339 or custom format) |
| `--to` | End timestamp (RFC3339 or custom format) |
| `--pattern` | Regex pattern to match log lines |
| `--format` | Timestamp format layout (default: RFC3339) |
| `--output` | Output file (default: stdout) |

---

## License

MIT — see [LICENSE](LICENSE) for details.