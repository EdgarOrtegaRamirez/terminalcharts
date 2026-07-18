package chart

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/EdgarOrtegaRamirez/terminalcharts/internal/data"
)

// ChartOptions holds configuration for chart rendering.
type ChartOptions struct {
	Width    int
	Height   int
	Sort     bool
	Desc     bool
	Colors   []string
	NoLegend bool
	Title    string
	Decimals int
}

func (o ChartOptions) getColor(i int, max int) string {
	if len(o.Colors) == 0 {
		colors := []string{
			"\033[32m", "\033[32m", "\033[92m", "\033[33m",
			"\033[33m", "\033[36m", "\033[34m", "\033[35m",
		}
		return colors[i%len(colors)]
	}
	return o.Colors[i%len(o.Colors)]
}

func sortDataset(ds *data.Dataset, desc bool) *data.Dataset {
	rows := make([]data.Row, len(ds.Rows))
	copy(rows, ds.Rows)
	sort.Slice(rows, func(i, j int) bool {
		if desc {
			return rows[i].Value > rows[j].Value
		}
		return rows[i].Value < rows[j].Value
	})
	return &data.Dataset{
		Rows:   rows,
		Header: ds.Header,
	}
}

func maxFloat(rows []data.Row) float64 {
	if len(rows) == 0 {
		return 1
	}
	m := rows[0].Value
	for _, r := range rows[1:] {
		if r.Value > m {
			m = r.Value
		}
	}
	return m
}

func minFloat(rows []data.Row) float64 {
	if len(rows) == 0 {
		return 0
	}
	m := rows[0].Value
	for _, r := range rows[1:] {
		if r.Value < m {
			m = r.Value
		}
	}
	return m
}

func sumFloat(rows []data.Row) float64 {
	s := 0.0
	for _, r := range rows {
		s += r.Value
	}
	return s
}

func formatNum(v float64, decimals int) string {
	format := "%." + strconv.Itoa(decimals) + "f"
	return fmt.Sprintf(format, v)
}

// BarChart generates a horizontal bar chart.
func BarChart(ds *data.Dataset, opts ChartOptions) (*ChartResult, error) {
	ds = sortDataset(ds, opts.Desc)
	rows := ds.Rows

	if len(rows) == 0 {
		return &ChartResult{Text: "(no data)\n"}, nil
	}

	maxVal := maxFloat(rows)
	if maxVal == 0 {
		maxVal = 1
	}

	maxLabelLen := 0
	for _, r := range rows {
		if len(r.Label) > maxLabelLen {
			maxLabelLen = len(r.Label)
		}
	}

	barAreaWidth := opts.Width - maxLabelLen - 4
	if barAreaWidth < 4 {
		barAreaWidth = 4
	}

	var sb strings.Builder

	if opts.Title != "" {
		sb.WriteString(center(opts.Title, opts.Width) + "\n\n")
	}

	sb.WriteString(" Bar Chart\n\n")

	for i, r := range rows {
		label := r.Label
		shortLabel := label
		if maxLabelLen > 5 && len(label) > maxLabelLen-2 {
			shortLabel = label[:maxLabelLen-4] + ".."
		}
		barLen := int(math.Round(float64(r.Value/maxVal) * float64(barAreaWidth)))
		if barLen < 1 && r.Value > 0 {
			barLen = 1
		}

		col := opts.getColor(i, len(rows))
		bar := strings.Repeat("█", barLen)

		sb.WriteString(fmt.Sprintf("  %s%s| %s%s %s\n",
			shortLabel,
			strings.Repeat(" ", maxLabelLen-len(shortLabel)),
			col+bar+"\033[0m",
			strings.Repeat(" ", barAreaWidth-barLen),
			formatNum(r.Value, opts.Decimals)))
	}

	if !opts.NoLegend {
		sb.WriteString(fmt.Sprintf("\n  scale: 0 to %s\n", formatNum(maxVal, opts.Decimals)))
		sb.WriteString(fmt.Sprintf("  %d items\n", len(rows)))
	}

	return &ChartResult{
		Text:   sb.String(),
		Rows:   len(rows),
		Width:  opts.Width,
		Height: opts.Height,
	}, nil
}

// LineChart generates an ASCII line plot.
func LineChart(ds *data.Dataset, opts ChartOptions) (*ChartResult, error) {
	rows := ds.Rows
	if len(rows) == 0 {
		return &ChartResult{Text: "(no data)\n"}, nil
	}

	maxVal := maxFloat(rows)
	minVal := minFloat(rows)
	valRange := maxVal - minVal
	if valRange == 0 {
		valRange = 1
	}

	height := opts.Height - 2
	if height < 1 {
		height = 1
	}

	width := opts.Width
	if width > len(rows)+2 {
		width = len(rows) + 2
	}
	if width < 4 {
		width = 4
	}

	var sb strings.Builder

	if opts.Title != "" {
		sb.WriteString(center(opts.Title, opts.Width) + "\n\n")
	}

	sb.WriteString(" Line Plot\n\n")

	for row := height; row >= 0; row-- {
		yVal := minVal + (valRange * float64(row) / float64(height))

		if row == height {
			sb.WriteString(fmt.Sprintf("  %s |", formatNum(yVal, opts.Decimals)))
		} else if row == 0 {
			sb.WriteString(fmt.Sprintf("  %s +", formatNum(yVal, opts.Decimals)))
		} else {
			sb.WriteString("    |")
		}

		step := float64(len(rows)) / float64(width-2)
		for col := 0; col < width-2; col++ {
			dataIdx := int(step*float64(col) + step/2)
			if dataIdx >= len(rows) {
				dataIdx = len(rows) - 1
			}

			dataVal := rows[dataIdx].Value
			norm := (dataVal - minVal) / valRange
			yPos := int(math.Round(norm * float64(height)))

			if yPos == row {
				sb.WriteString("*")
			} else {
				sb.WriteString(" ")
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("  %s+", strings.Repeat("-", width-2)))
	sb.WriteString("\n")

	if len(rows) > 0 {
		step := float64(len(rows)) / float64(width-2)
		xLabels := ""
		for col := 0; col < width-2; col++ {
			dataIdx := int(step * float64(col))
			if dataIdx >= len(rows) {
				dataIdx = len(rows) - 1
			}
			label := rows[dataIdx].Label
			if len(label) > 1 {
				label = label[:1]
			}
			xLabels += label
		}
		sb.WriteString(fmt.Sprintf("    %s\n", xLabels))
	}

	if !opts.NoLegend {
		sb.WriteString(fmt.Sprintf("\n  min: %s  max: %s\n", formatNum(minVal, opts.Decimals), formatNum(maxVal, opts.Decimals)))
		sb.WriteString(fmt.Sprintf("  %d data points\n", len(rows)))
	}

	return &ChartResult{
		Text:   sb.String(),
		Rows:   len(rows),
		Width:  opts.Width,
		Height: opts.Height,
	}, nil
}

// PieChart generates a horizontal pie chart (segmented bar).
func PieChart(ds *data.Dataset, opts ChartOptions) (*ChartResult, error) {
	ds = sortDataset(ds, opts.Desc)
	rows := ds.Rows

	if len(rows) == 0 {
		return &ChartResult{Text: "(no data)\n"}, nil
	}

	total := sumFloat(rows)
	if total == 0 {
		total = 1
	}

	var sb strings.Builder

	if opts.Title != "" {
		sb.WriteString(center(opts.Title, opts.Width) + "\n\n")
	}

	sb.WriteString(" Pie Chart\n\n")

	pieWidth := opts.Width - 2
	if pieWidth < 10 {
		pieWidth = 10
	}
	if pieWidth > 40 {
		pieWidth = 40
	}

	cumulative := 0.0
	type seg struct {
		label string
		value float64
		pct   float64
	}

	segments := make([]seg, 0, len(rows))
	for _, r := range rows {
		pct := r.Value / total * 100
		segments = append(segments, seg{
			label: r.Label,
			value: r.Value,
			pct:   pct,
		})
		cumulative += pct / 100
	}

	for i, seg := range segments {
		barLen := int(math.Round(seg.pct / 100 * float64(pieWidth)))
		if barLen < 1 && seg.pct > 0 {
			barLen = 1
		}

		label := seg.label
		if len(label) > 20 {
			label = label[:18] + ".."
		}

		col := opts.getColor(i, len(segments))

		sb.WriteString(fmt.Sprintf("  %s%-20s %s%s %5.1f%%\n",
			label,
			strings.Repeat(" ", 20-len(label)),
			col+strings.Repeat("█", barLen)+"\033[0m",
			strings.Repeat(" ", pieWidth-barLen),
			seg.pct))
	}

	if !opts.NoLegend {
		sb.WriteString(fmt.Sprintf("\n  total: %s\n", formatNum(total, opts.Decimals)))
		sb.WriteString(fmt.Sprintf("  %d segments\n", len(segments)))
	}

	return &ChartResult{
		Text:   sb.String(),
		Rows:   len(rows),
		Width:  opts.Width,
		Height: opts.Height,
	}, nil
}

// HistogramChart generates a histogram chart.
func HistogramChart(ds *data.Dataset, opts ChartOptions) (*ChartResult, error) {
	rows := ds.Rows
	if len(rows) == 0 {
		return &ChartResult{Text: "(no data)\n"}, nil
	}

	bucketCount := opts.Height - 2
	if bucketCount < 3 {
		bucketCount = 3
	}
	if bucketCount > 20 {
		bucketCount = 20
	}

	nBuckets := len(rows)
	if nBuckets > bucketCount {
		nBuckets = bucketCount
	}
	if nBuckets < 2 {
		nBuckets = 2
	}

	maxVal := maxFloat(rows)
	minVal := minFloat(rows)
	valRange := maxVal - minVal
	if valRange == 0 {
		valRange = 1
	}

	bucketSize := valRange / float64(nBuckets)

	buckets := make([]int, nBuckets)
	labels := make([]string, nBuckets)
	for _, r := range rows {
		bucketIdx := int((r.Value - minVal) / bucketSize)
		if bucketIdx >= nBuckets {
			bucketIdx = nBuckets - 1
		}
		if bucketIdx < 0 {
			bucketIdx = 0
		}
		buckets[bucketIdx]++

		low := minVal + float64(bucketIdx)*bucketSize
		high := low + bucketSize
		labels[bucketIdx] = formatNum(low, 0) + "-" + formatNum(high, 0)
	}

	maxBucket := 0
	for _, b := range buckets {
		if b > maxBucket {
			maxBucket = b
		}
	}
	if maxBucket == 0 {
		maxBucket = 1
	}

	barAreaWidth := opts.Width - 15
	if barAreaWidth < 4 {
		barAreaWidth = 4
	}

	var sb strings.Builder

	if opts.Title != "" {
		sb.WriteString(center(opts.Title, opts.Width) + "\n\n")
	}

	sb.WriteString(" Histogram\n\n")

	for bi := 0; bi < len(buckets); bi++ {
		b := buckets[bi]
		label := labels[bi]
		barLen := int(math.Round(float64(b) / float64(maxBucket) * float64(barAreaWidth)))
		if barLen < 1 && b > 0 {
			barLen = 1
		}

		col := opts.getColor(bi, len(buckets))

		sb.WriteString(fmt.Sprintf("  %-10s %s%s %d\n",
			label,
			col+strings.Repeat("█", barLen)+"\033[0m",
			strings.Repeat(" ", barAreaWidth-barLen),
			b))
	}

	if !opts.NoLegend {
		sb.WriteString(fmt.Sprintf("\n  range: %s to %s\n", formatNum(minVal, opts.Decimals), formatNum(maxVal, opts.Decimals)))
		sb.WriteString(fmt.Sprintf("  %d buckets, %d values\n", nBuckets, len(rows)))
	}

	return &ChartResult{
		Text:   sb.String(),
		Rows:   len(rows),
		Width:  opts.Width,
		Height: opts.Height,
	}, nil
}

func center(s string, width int) string {
	if len(s) >= width {
		return s
	}
	pad := (width - len(s)) / 2
	return strings.Repeat(" ", pad) + s + strings.Repeat(" ", width-pad-len(s))
}

// ChartResult holds the output of a chart generation.
type ChartResult struct {
	Text   string `json:"text"`
	Rows   int    `json:"rows"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
