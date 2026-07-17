package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/EdgarOrtegaRamirez/terminalcharts/internal/chart"
	"github.com/EdgarOrtegaRamirez/terminalcharts/internal/data"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("terminalcharts — ASCII charts from CSV/JSON data")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  terminalcharts bar <file> [options]    Bar chart")
		fmt.Println("  terminalcharts line <file> [options]   Line plot")
		fmt.Println("  terminalcharts pie <file> [options]    Pie chart")
		fmt.Println("  terminalcharts hist <file> [options]   Histogram")
		fmt.Println("  terminalcharts info <file>             Show file info")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  --header                       First row is header")
		fmt.Println("  --columns <names>              Comma-separated column names")
		fmt.Println("  --label-col <name>             Column for labels")
		fmt.Println("  --value-col <name>             Column for values")
		fmt.Println("  --width <n>                    Chart width (default: 40)")
		fmt.Println("  --height <n>                   Chart height (default: 15)")
		fmt.Println("  --sort                         Sort by value (default: yes)")
		fmt.Println("  --desc                         Sort descending (default)")
		fmt.Println("  --format <text|json|jsonl>     Output format (default: text)")
		fmt.Println("  --colors <list>                Comma-separated color codes")
		fmt.Println("  --no-legend                    Hide legend")
		fmt.Println("  --title <str>                  Chart title")
		fmt.Println("  --decimals <n>                 Decimal places (default: 2)")
		os.Exit(0)
	}

	command := strings.ToLower(os.Args[1])

	switch command {
	case "bar", "line", "pie", "hist":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Error: missing file argument\n")
			os.Exit(1)
		}
		file := os.Args[2]
		flags := parseFlags(os.Args[3:])

		ds, err := loadData(file, flags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		result, err := buildChart(command, ds, flags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if flags.format == "json" || flags.format == "jsonl" {
			writeOutput(result, flags.format == "jsonl")
			return
		}
		fmt.Print(result.Text)

	case "info":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Error: missing file argument\n")
			os.Exit(1)
		}
		file := os.Args[2]
		info, err := data.ReadInfo(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("File: %s\n", info.Path)
		fmt.Printf("Type: %s\n", info.Format)
		fmt.Printf("Rows: %d\n", info.Rows)
		fmt.Printf("Columns: %d\n", info.Cols)
		if len(info.Header) > 0 {
			fmt.Printf("Columns: %s\n", strings.Join(info.Header, ", "))
		}

	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n", command)
		os.Exit(1)
	}
}

type configFlags struct {
	hasHeader bool
	columns   []string
	labelCol  string
	valueCol  string
	width     int
	height    int
	sort      bool
	desc      bool
	format    string
	colors    []string
	noLegend  bool
	title     string
	decimals  int
}

func parseFlags(args []string) configFlags {
	cfg := configFlags{
		width:    40,
		height:   15,
		sort:     true,
		desc:     true,
		format:   "text",
		decimals: 2,
	}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--header":
			cfg.hasHeader = true
		case "--columns":
			if i+1 < len(args) {
				cfg.columns = strings.Split(args[i+1], ",")
				i++
			}
		case "--label-col":
			if i+1 < len(args) {
				cfg.labelCol = args[i+1]
				i++
			}
		case "--value-col":
			if i+1 < len(args) {
				cfg.valueCol = args[i+1]
				i++
			}
		case "--width":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &cfg.width)
				i++
			}
		case "--height":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &cfg.height)
				i++
			}
		case "--sort":
			cfg.sort = true
		case "--no-sort":
			cfg.sort = false
		case "--desc":
			cfg.desc = true
		case "--asc":
			cfg.desc = false
		case "--format":
			if i+1 < len(args) {
				cfg.format = args[i+1]
				i++
			}
		case "--colors":
			if i+1 < len(args) {
				cfg.colors = strings.Split(args[i+1], ",")
				i++
			}
		case "--no-legend":
			cfg.noLegend = true
		case "--title":
			if i+1 < len(args) {
				cfg.title = args[i+1]
				i++
			}
		case "--decimals":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &cfg.decimals)
				i++
			}
		}
	}
	return cfg
}

func loadData(path string, cfg configFlags) (*data.Dataset, error) {
	return data.Load(path, data.LoadConfig{
		HasHeader: cfg.hasHeader,
		Columns:   cfg.columns,
		LabelCol:  cfg.labelCol,
		ValueCol:  cfg.valueCol,
	})
}

func buildChart(cmd string, ds *data.Dataset, cfg configFlags) (*chart.ChartResult, error) {
	opts := chart.ChartOptions{
		Width:    cfg.width,
		Height:   cfg.height,
		Sort:     cfg.sort,
		Desc:     cfg.desc,
		Colors:   cfg.colors,
		NoLegend: cfg.noLegend,
		Title:    cfg.title,
		Decimals: cfg.decimals,
	}

	switch cmd {
	case "bar":
		return chart.BarChart(ds, opts)
	case "line":
		return chart.LineChart(ds, opts)
	case "pie":
		return chart.PieChart(ds, opts)
	case "hist":
		return chart.HistogramChart(ds, opts)
	}
	return nil, fmt.Errorf("unknown chart type: %s", cmd)
}

func writeOutput(result *chart.ChartResult, asJSONL bool) {
	if asJSONL {
		encoder := json.NewEncoder(os.Stdout)
		encoder.Encode(result)
		return
	}
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON marshal error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}