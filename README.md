# TerminalCharts

ASCII charts from CSV/JSON data — generate bar charts, line plots, pie charts, and histograms directly in your terminal.

## Features

- **Bar charts** — Horizontal bar charts with sorted values
- **Line plots** — ASCII line plots with axis labels
- **Pie charts** — Segmented pie charts showing proportions
- **Histograms** — Distribution histograms with configurable buckets
- **Multi-format input** — CSV and JSON data files
- **Multiple output** — Text (terminal) or JSON/JSONL for CI/CD
- **Configurable** — Width, height, colors, sorting, labels, titles

## Installation

### From source

```bash
go install github.com/EdgarOrtegaRamirez/terminalcharts/cmd/terminalcharts@latest
```

### Build from source

```bash
git clone https://github.com/EdgarOrtegaRamirez/terminalcharts
cd terminalcharts
go build -o terminalcharts ./cmd/terminalcharts/
```

## Usage

### Bar Chart

```bash
terminalcharts bar data.csv
terminalcharts bar data.csv --label-col category --value-col sales --title "Sales by Category"
```

### Line Plot

```bash
terminalcharts line data.json
terminalcharts line data.csv --title "Weekly Trend"
```

### Pie Chart

```bash
terminalcharts pie data.csv
terminalcharts pie data.json --title "Market Share"
```

### Histogram

```bash
terminalcharts hist data.csv
terminalcharts hist data.csv --height 10 --title "Value Distribution"
```

### File Info

```bash
terminalcharts info data.csv
```

## Options

| Option | Description | Default |
|--------|-------------|---------|
| `--header` | First row is header | auto-detected |
| `--columns <names>` | Comma-separated column names | — |
| `--label-col <name>` | Column for labels | auto-detected |
| `--value-col <name>` | Column for values | auto-detected |
| `--width <n>` | Chart width | 40 |
| `--height <n>` | Chart height | 15 |
| `--sort` / `--no-sort` | Sort by value | enabled |
| `--desc` / `--asc` | Sort direction | descending |
| `--format <text|json|jsonl>` | Output format | text |
| `--colors <list>` | Comma-separated color codes | default palette |
| `--no-legend` | Hide legend | shown |
| `--title <str>` | Chart title | — |
| `--decimals <n>` | Decimal places | 2 |

## Input Formats

### CSV

```csv
name,value
Electronics,500
Clothing,300
Books,100
```

### JSON — Array of objects

```json
[{"label":"A","value":10},{"label":"B","value":20}]
```

### JSON — Object

```json
{"a": 1, "b": 2, "c": 3}
```

### JSON — Simple array

```json
[10, 20, 30]
```

## JSON Output

```bash
terminalcharts bar data.csv --format json
```

```json
{
  "text": "Bar Chart\n\n  A | █████ 10\n  B | ████████ 20\n  C | ██████████ 30\n",
  "rows": 3,
  "width": 40,
  "height": 15
}
```

## CI/CD Integration

```yaml
- name: Check chart output
  run: |
    terminalcharts bar data.csv --format json | jq '.rows > 0'
```

## License

MIT — see [LICENSE](LICENSE) file.