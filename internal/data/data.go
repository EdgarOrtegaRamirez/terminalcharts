package data

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Dataset holds parsed data rows with labels and values.
type Dataset struct {
	Header   []string
	Rows     []Row
	FileType string
}

// Row represents a single data row.
type Row struct {
	Label string
	Value float64
}

// LoadConfig holds options for loading data.
type LoadConfig struct {
	HasHeader bool
	Columns   []string
	LabelCol  string
	ValueCol  string
}

// ReadInfo returns basic file metadata without full parsing.
type FileReadInfo struct {
	Path   string
	Format string
	Rows   int
	Cols   int
	Header []string
}

// ReadInfo returns metadata about a data file.
func ReadInfo(path string) (*FileReadInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	info := &FileReadInfo{Path: path}

	content := strings.TrimSpace(string(data))
	if len(content) == 0 {
		return nil, fmt.Errorf("empty file")
	}

	if content[0] == '{' || content[0] == '[' {
		info.Format = "json"
		lines := strings.Split(content, "\n")
		info.Rows = len(lines)
		var arr []json.RawMessage
		if err := json.Unmarshal([]byte(content), &arr); err == nil {
			info.Cols = len(arr)
		}
		return info, nil
	}

	info.Format = "csv"
	info.Rows = 0
	info.Cols = 0

	r := csv.NewReader(strings.NewReader(content))
	allRows, err := r.ReadAll()
	if err != nil {
		content = strings.ReplaceAll(content, "\r\n", "\n")
		content = strings.ReplaceAll(content, "\r", "\n")
		allRows, err = csv.NewReader(strings.NewReader(content)).ReadAll()
		if err != nil {
			return nil, fmt.Errorf("parse CSV: %w", err)
		}
	}

	info.Rows = len(allRows)
	if len(allRows) > 0 {
		info.Cols = len(allRows[0])
		info.Header = allRows[0]
	}

	return info, nil
}

// Load reads a data file and returns a Dataset.
func Load(path string, cfg LoadConfig) (*Dataset, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	trimmed := strings.TrimSpace(string(data))
	if len(trimmed) == 0 {
		return nil, fmt.Errorf("empty file: %s", path)
	}

	if trimmed[0] == '{' || trimmed[0] == '[' {
		return loadJSON(trimmed, cfg)
	}

	return loadCSV(trimmed, cfg)
}

func loadJSON(content string, cfg LoadConfig) (*Dataset, error) {
	var objects []map[string]json.RawMessage
	if err := json.Unmarshal([]byte(content), &objects); err == nil {
		return parseJSONArray(objects, cfg)
	}

	var simple []float64
	if err := json.Unmarshal([]byte(content), &simple); err == nil {
		rows := make([]Row, len(simple))
		for i, v := range simple {
			rows[i] = Row{Label: strconv.Itoa(i + 1), Value: v}
		}
		return &Dataset{
			Rows:     rows,
			FileType: "json",
		}, nil
	}

	var obj map[string]float64
	if err := json.Unmarshal([]byte(content), &obj); err == nil {
		return &Dataset{
			Header:   stringKeys(obj),
			Rows:     mapToRows(obj),
			FileType: "json",
		}, nil
	}

	return nil, fmt.Errorf("unsupported JSON format")
}

func parseJSONArray(objects []map[string]json.RawMessage, cfg LoadConfig) (*Dataset, error) {
	if len(objects) == 0 {
		return &Dataset{FileType: "json"}, nil
	}

	labelKey := ""
	valueKey := ""

	if cfg.LabelCol != "" {
		labelKey = cfg.LabelCol
	}
	if cfg.ValueCol != "" {
		valueKey = cfg.ValueCol
	}

	if labelKey == "" || valueKey == "" {
		for k := range objects[0] {
			switch k {
			case "label", "name", "x", "key":
				if labelKey == "" {
					labelKey = k
				}
			case "value", "y", "count", "amount", "data":
				if valueKey == "" {
					valueKey = k
				}
			}
		}
		if labelKey == "" {
			for k := range objects[0] {
				if labelKey == "" {
					labelKey = k
				}
				break
			}
		}
		if valueKey == "" {
			for k, v := range objects[0] {
				var num float64
				if err := json.Unmarshal(v, &num); err == nil {
					valueKey = k
					break
				}
			}
		}
	}

	if valueKey == "" {
		return nil, fmt.Errorf("could not determine value column in JSON array")
	}

	ds := &Dataset{FileType: "json"}
	for _, obj := range objects {
		row := Row{}
		if labelKey != "" {
			if raw, ok := obj[labelKey]; ok {
				var s string
				if err := json.Unmarshal(raw, &s); err == nil {
					row.Label = s
				}
			}
		}
		if raw, ok := obj[valueKey]; ok {
			var v float64
			if err := json.Unmarshal(raw, &v); err != nil {
				return nil, fmt.Errorf("value is not a number: %s", string(raw))
			}
			row.Value = v
		}
		ds.Rows = append(ds.Rows, row)
	}

	return ds, nil
}

func mapToRows(obj map[string]float64) []Row {
	rows := make([]Row, 0, len(obj))
	for label, value := range obj {
		rows = append(rows, Row{Label: label, Value: value})
	}
	return rows
}

func stringKeys(obj map[string]float64) []string {
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	return keys
}

func loadCSV(content string, cfg LoadConfig) (*Dataset, error) {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")

	r := csv.NewReader(strings.NewReader(content))
	allRows, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parse CSV: %w", err)
	}

	ds := &Dataset{FileType: "csv"}

	if len(allRows) == 0 {
		return ds, nil
	}

	hasHeader := cfg.HasHeader
	if !hasHeader && len(allRows) > 1 {
		allStr := true
		for _, cell := range allRows[0] {
			if _, err := strconv.ParseFloat(cell, 64); err == nil {
				allStr = false
				break
			}
		}
		if cfg.ValueCol != "" {
			hasHeader = false
		} else if allStr && len(allRows) > 1 {
			hasHeader = true
		}
	}

	start := 0
	if hasHeader {
		ds.Header = allRows[0]
		start = 1
	}

	valueColIdx := -1
	labelColIdx := -1

	if cfg.ValueCol != "" {
		for i, h := range allRows[0] {
			if h == cfg.ValueCol {
				valueColIdx = i
				break
			}
		}
		if valueColIdx < 0 {
			for i, h := range allRows[0] {
				if strings.EqualFold(h, cfg.ValueCol) {
					valueColIdx = i
					break
				}
			}
		}
	}

	if cfg.LabelCol != "" {
		for i, h := range allRows[0] {
			if h == cfg.LabelCol {
				labelColIdx = i
				break
			}
		}
		if labelColIdx < 0 {
			for i, h := range allRows[0] {
				if strings.EqualFold(h, cfg.LabelCol) {
					labelColIdx = i
					break
				}
			}
		}
	}

	if len(cfg.Columns) >= 2 {
		ds.Header = cfg.Columns
	}

	if valueColIdx < 0 {
		if labelColIdx >= 0 {
			valueColIdx = 1
			if valueColIdx == labelColIdx {
				valueColIdx = 0
			}
		} else {
			valueColIdx = 1
			if valueColIdx >= len(allRows[0]) {
				valueColIdx = 0
			}
		}
	}

	for i := start; i < len(allRows); i++ {
		row := allRows[i]
		if len(row) == 0 {
			continue
		}

		var label string
		if labelColIdx >= 0 && labelColIdx < len(row) {
			label = row[labelColIdx]
		} else if labelColIdx < 0 {
			if valueColIdx > 0 {
				label = row[0]
			} else {
				label = strconv.Itoa(i)
			}
		}

		var value float64
		if valueColIdx < len(row) {
			v, err := strconv.ParseFloat(row[valueColIdx], 64)
			if err != nil {
				i64, err2 := strconv.Atoi(row[valueColIdx])
				if err2 != nil {
					return nil, fmt.Errorf("row %d: invalid value '%s': %v", i, row[valueColIdx], err2)
				}
				value = float64(i64)
			} else {
				value = v
			}
		}

		ds.Rows = append(ds.Rows, Row{Label: label, Value: value})
	}

	return ds, nil
}