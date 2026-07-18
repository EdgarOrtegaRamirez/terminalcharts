package chart

import (
	"strings"
	"testing"

	"github.com/EdgarOrtegaRamirez/terminalcharts/internal/data"
)

func TestBarChart(t *testing.T) {
	ds := &data.Dataset{
		Rows: []data.Row{
			{Label: "A", Value: 10},
			{Label: "B", Value: 20},
			{Label: "C", Value: 30},
		},
	}

	result, err := BarChart(ds, ChartOptions{Width: 30, Sort: false, NoLegend: true})
	if err != nil {
		t.Fatalf("BarChart: %v", err)
	}
	if result.Rows != 3 {
		t.Errorf("expected 3 rows, got %d", result.Rows)
	}
	if !strings.Contains(result.Text, "Bar Chart") {
		t.Errorf("expected 'Bar Chart' in output, got: %s", result.Text)
	}
}

func TestBarChartEmpty(t *testing.T) {
	ds := &data.Dataset{}
	result, err := BarChart(ds, ChartOptions{Sort: false, NoLegend: true})
	if err != nil {
		t.Fatalf("BarChart: %v", err)
	}
	if !strings.Contains(result.Text, "(no data)") {
		t.Errorf("expected '(no data)', got: %s", result.Text)
	}
}

func TestLineChart(t *testing.T) {
	ds := &data.Dataset{
		Rows: []data.Row{
			{Label: "A", Value: 10},
			{Label: "B", Value: 20},
			{Label: "C", Value: 30},
			{Label: "D", Value: 25},
		},
	}

	result, err := LineChart(ds, ChartOptions{Width: 20, Sort: false, NoLegend: true})
	if err != nil {
		t.Fatalf("LineChart: %v", err)
	}
	if result.Rows != 4 {
		t.Errorf("expected 4 rows, got %d", result.Rows)
	}
	if !strings.Contains(result.Text, "Line Plot") {
		t.Errorf("expected 'Line Plot' in output, got: %s", result.Text)
	}
}

func TestLineChartEmpty(t *testing.T) {
	ds := &data.Dataset{}
	result, err := LineChart(ds, ChartOptions{Sort: false, NoLegend: true})
	if err != nil {
		t.Fatalf("LineChart: %v", err)
	}
	if !strings.Contains(result.Text, "(no data)") {
		t.Errorf("expected '(no data)', got: %s", result.Text)
	}
}

func TestPieChart(t *testing.T) {
	ds := &data.Dataset{
		Rows: []data.Row{
			{Label: "A", Value: 50},
			{Label: "B", Value: 30},
			{Label: "C", Value: 20},
		},
	}

	result, err := PieChart(ds, ChartOptions{Width: 30, Sort: false, NoLegend: true})
	if err != nil {
		t.Fatalf("PieChart: %v", err)
	}
	if result.Rows != 3 {
		t.Errorf("expected 3 rows, got %d", result.Rows)
	}
	if !strings.Contains(result.Text, "Pie Chart") {
		t.Errorf("expected 'Pie Chart' in output, got: %s", result.Text)
	}
}

func TestPieChartEmpty(t *testing.T) {
	ds := &data.Dataset{}
	result, err := PieChart(ds, ChartOptions{Sort: false, NoLegend: true})
	if err != nil {
		t.Fatalf("PieChart: %v", err)
	}
	if !strings.Contains(result.Text, "(no data)") {
		t.Errorf("expected '(no data)', got: %s", result.Text)
	}
}

func TestHistogramChart(t *testing.T) {
	ds := &data.Dataset{
		Rows: []data.Row{
			{Label: "A", Value: 10},
			{Label: "B", Value: 20},
			{Label: "C", Value: 30},
			{Label: "D", Value: 40},
			{Label: "E", Value: 50},
		},
	}

	result, err := HistogramChart(ds, ChartOptions{Width: 30, Height: 10, NoLegend: true})
	if err != nil {
		t.Fatalf("HistogramChart: %v", err)
	}
	if result.Rows != 5 {
		t.Errorf("expected 5 rows, got %d", result.Rows)
	}
	if !strings.Contains(result.Text, "Histogram") {
		t.Errorf("expected 'Histogram' in output, got: %s", result.Text)
	}
}

func TestHistogramEmpty(t *testing.T) {
	ds := &data.Dataset{}
	result, err := HistogramChart(ds, ChartOptions{NoLegend: true})
	if err != nil {
		t.Fatalf("HistogramChart: %v", err)
	}
	if !strings.Contains(result.Text, "(no data)") {
		t.Errorf("expected '(no data)', got: %s", result.Text)
	}
}

func TestCenter(t *testing.T) {
	result := center("hello", 10)
	if result != "  hello   " {
		t.Errorf("expected '  hello   ', got '%s'", result)
	}
}

func TestFormatNum(t *testing.T) {
	result := formatNum(3.14159, 2)
	if result != "3.14" {
		t.Errorf("expected '3.14', got '%s'", result)
	}
}
