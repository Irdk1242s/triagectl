package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/plonxyz/triagectl/internal/analysis"
	"github.com/plonxyz/triagectl/internal/collectors"
	"github.com/plonxyz/triagectl/internal/models"
	"github.com/plonxyz/triagectl/internal/output"
	"github.com/plonxyz/triagectl/internal/progress"
	"github.com/plonxyz/triagectl/internal/report"
)

const version = "0.2.0"

func main() {
	// Command line flags
	outputDir := flag.String("output", "./triagectl-output", "Output directory for collected artifacts")
	listCollectors := flag.Bool("list", false, "List available collectors and exit")
	showVersion := flag.Bool("version", false, "Show version and exit")
	timeout := flag.Int("timeout", 300, "Global timeout in seconds for collection")
	collectorFilter := flag.String("collectors", "", "Comma-separated collector IDs to run (default: all)")
	collectorTimeout := flag.Int("collector-timeout", 60, "Per-collector timeout in seconds")
	concurrency := flag.Int("concurrency", 4, "Maximum number of collectors to run concurrently")
	enableCSV := flag.Bool("csv", false, "Enable CSV output")
	enableHTML := flag.Bool("html", false, "Generate HTML report")
	enableTimeline := flag.Bool("timeline", false, "Generate timeline.csv (Timesketch format)")
	iocFile := flag.String("ioc-file", "", "Path to IOC file (one indicator per line)")
	flag.Parse()

	if *showVersion {
		fmt.Printf("triagectl v%s\n", version)
		os.Exit(0)
	}

	if *listCollectors {
		printCollectors()
		os.Exit(0)
	}

	banner := fmt.Sprintf("triagectl v%s", version)
	const boxWidth = 39
	pad := boxWidth - len(banner)
	left := pad / 2
	right := pad - left
	fmt.Println("╔═══════════════════════════════════════╗")
	fmt.Printf("║%s%s%s║\n", strings.Repeat(" ", left), banner, strings.Repeat(" ", right))
	fmt.Println("╚═══════════════════════════════════════╝")
	fmt.Println()

	// 1. Create output directory
	ts := time.Now().Format("20060102-150405")
	hostname, _ := os.Hostname()
	collectionDir := filepath.Join(*outputDir, fmt.Sprintf("%s-%s", hostname, ts))

	if err := os.MkdirAll(collectionDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Output directory: %s\n\n", collectionDir)

	// 2. Init writers: SQLite + optionally CSV → MultiWriter
	sqlitePath := filepath.Join(collectionDir, "artifacts.db")

	sqliteWriter, err := output.NewSQLiteWriter(sqlitePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating SQLite writer: %v\n", err)
		os.Exit(1)
	}
	defer sqliteWriter.Close()

	writers := []output.Writer{sqliteWriter}

	var csvPath string
	if *enableCSV {
		csvPath = filepath.Join(collectionDir, "artifacts.csv")
		csvWriter, err := output.NewCSVWriter(csvPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating CSV writer: %v\n", err)
			os.Exit(1)
		}
		defer csvWriter.Close()
		writers = append(writers, csvWriter)
	}

	multiWriter := output.NewMultiWriter(writers...)
	defer multiWriter.Close()

	// 3. Filter collectors via --collectors
	activeCollectors := filterCollectors(*collectorFilter)

	// 4. Load IOCs if --ioc-file provided
	if *iocFile != "" {
		iocMatcher, err := analysis.NewIOCMatcher(*iocFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading IOC file: %v\n", err)
			os.Exit(1)
		}
		analysis.RegisterAnalyzer(iocMatcher)
		fmt.Printf("Loaded IOC file: %s\n\n", *iocFile)
	}

	// 5. Start progress tracker
	tracker := progress.NewTracker(len(activeCollectors))

	// 6. Launch collectors with semaphore + per-collector timeouts
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()

	startTime := time.Now()
	fmt.Println("Starting artifact collection...")
	fmt.Println()

	semaphore := make(chan struct{}, *concurrency)
	var wg sync.WaitGroup
	resultsCh := make(chan models.CollectionResult, len(activeCollectors))

	for _, collector := range activeCollectors {
		wg.Add(1)
		go func(c collectors.Collector) {
			defer wg.Done()

			semaphore <- struct{}{}        // acquire
			defer func() { <-semaphore }() // release

			tracker.Start(c.Name())

			collectorCtx, collectorCancel := context.WithTimeout(ctx, time.Duration(*collectorTimeout)*time.Second)
			defer collectorCancel()

			start := time.Now()
			artifacts, err := c.Collect(collectorCtx)
			dur := time.Since(start)

			result := models.CollectionResult{
				CollectorID: c.ID(),
				Artifacts:   artifacts,
				Error:       err,
				Duration:    dur,
				StartedAt:   start,
			}

			resultsCh <- result
		}(collector)
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// 7. Process results: write to MultiWriter, update progress
	var allArtifacts []models.Artifact
	var allResults []models.CollectionResult
	totalArtifacts := 0
	successfulCollectors := 0
	failedCollectors := 0

	for result := range resultsCh {
		allResults = append(allResults, result)

		if result.Error != nil {
			tracker.Fail(result.CollectorID, result.Error)
			failedCollectors++
			continue
		}

		if len(result.Artifacts) > 0 {
			allArtifacts = append(allArtifacts, result.Artifacts...)
			totalArtifacts += len(result.Artifacts)
			tracker.Success(result.CollectorID, len(result.Artifacts))
		} else {
			tracker.Success(result.CollectorID, 0)
		}
		successfulCollectors++
	}

	tracker.Finish()
	duration := time.Since(startTime)

	// 8. Run cross-artifact analyzers
	fmt.Println("\nRunning analysis...")
	allArtifacts = analysis.RunAll(allArtifacts)

	// 9. Write analyzed artifacts to all output formats
	if err := multiWriter.WriteMany(allArtifacts); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing artifacts: %v\n", err)
	}

	// Update severity counts
	findingsCount := 0
	for _, a := range allArtifacts {
		if a.RiskScore >= 40 {
			findingsCount++
		}
	}
	fmt.Printf("  Analysis complete: %d findings detected\n", findingsCount)

	// 9. If --timeline: generate timeline.csv
	if *enableTimeline {
		timelinePath := filepath.Join(collectionDir, "timeline.csv")
		if err := output.GenerateTimeline(allArtifacts, timelinePath, report.Summarize); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating timeline: %v\n", err)
		} else {
			fmt.Printf("  Timeline: %s\n", timelinePath)
		}
	}

	// 10. If --html: generate report.html
	if *enableHTML {
		reportPath := filepath.Join(collectionDir, "report.html")
		if err := report.GenerateHTMLReport(reportPath, allArtifacts, allResults, duration); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating HTML report: %v\n", err)
		} else {
			fmt.Printf("  HTML Report: %s\n", reportPath)
		}
	}

	// 11. Print summary
	fmt.Println()
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("Collection Summary")
	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("Duration: %v\n", duration.Round(time.Millisecond))
	fmt.Printf("Total Artifacts: %d\n", totalArtifacts)
	fmt.Printf("Successful Collectors: %d\n", successfulCollectors)
	fmt.Printf("Failed Collectors: %d\n", failedCollectors)

	if findingsCount > 0 {
		fmt.Printf("Findings (risk >= medium): %d\n", findingsCount)
	}

	fmt.Println()

	// Print detailed stats from SQLite
	if err := sqliteWriter.PrintStats(); err != nil {
		fmt.Fprintf(os.Stderr, "Error printing stats: %v\n", err)
	}

	fmt.Println()
	fmt.Println("Output files:")
	fmt.Printf("  - SQLite: %s\n", sqlitePath)
	if *enableCSV {
		fmt.Printf("  - CSV:    %s\n", csvPath)
	}
	if *enableTimeline {
		fmt.Printf("  - Timeline: %s\n", filepath.Join(collectionDir, "timeline.csv"))
	}
	if *enableHTML {
		fmt.Printf("  - Report: %s\n", filepath.Join(collectionDir, "report.html"))
	}
	fmt.Println()
	fmt.Println("Collection complete!")
}

func printCollectors() {
	fmt.Println("Available Collectors:")
	fmt.Println("====================")
	for _, collector := range collectors.Registry {
		rootRequired := ""
		if collector.RequiresRoot() {
			rootRequired = " [REQUIRES ROOT]"
		}
		fmt.Printf("  %-25s - %s%s\n", collector.ID(), collector.Description(), rootRequired)
	}
	fmt.Printf("\nTotal: %d collectors\n", len(collectors.Registry))
}

func filterCollectors(filter string) []collectors.Collector {
	if filter == "" {
		return collectors.Registry
	}

	allowed := make(map[string]bool)
	for _, id := range strings.Split(filter, ",") {
		id = strings.TrimSpace(id)
		if id != "" {
			allowed[id] = true
		}
	}

	var filtered []collectors.Collector
	for _, c := range collectors.Registry {
		if allowed[c.ID()] {
			filtered = append(filtered, c)
		}
	}

	if len(filtered) == 0 {
		fmt.Fprintf(os.Stderr, "Warning: no collectors matched filter '%s', running all\n", filter)
		return collectors.Registry
	}

	return filtered
}
