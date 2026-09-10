package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"compliance-report-generator/internal/models"
	"compliance-report-generator/internal/processor"
	"compliance-report-generator/internal/reader"
	"compliance-report-generator/internal/reporter"
	"compliance-report-generator/internal/validator"
)

const version = "1.0.0"

func main() {
	inputFlag := flag.String("input", "testdata/sample_compliance_data.csv", "Path to input CSV file(s) or glob pattern (e.g. 'testdata/*.csv')")
	outputFlag := flag.String("output", "compliance_audit_report.xlsx", "Path for the output Excel report (.xlsx)")
	verboseFlag := flag.Bool("verbose", true, "Enable verbose console logging")
	versionFlag := flag.Bool("version", false, "Print version and exit")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("Compliance Data Quality & Excel Report Generator v%s\n", version)
		os.Exit(0)
	}

	startTime := time.Now()

	printBanner()

	fmt.Printf("[1/5] Ingesting CSV dataset from: %s\n", *inputFlag)
	csvReader := reader.NewCSVReader()
	rawRecords, err := csvReader.ReadFiles(*inputFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Ingestion Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("      -> Successfully read %d raw records\n", len(rawRecords))

	fmt.Println("[2/5] Running Compliance & ACA Data Quality Validations...")
	valEngine := validator.NewValidationEngine()
	validatedRecords := valEngine.ValidateRecords(rawRecords)

	fmt.Println("[3/5] Running Deduplication & Cross-Record Identity Audit...")
	deduper := processor.NewDeduplicator()
	auditedRecords, dupCount := deduper.ProcessDuplicates(validatedRecords)
	fmt.Printf("      -> Found %d duplicate/colliding record instances\n", dupCount)

	fmt.Println("[4/5] Computing Summary Statistics & Department Breakdowns...")
	statsProcessor := processor.NewStatsProcessor()
	stats := statsProcessor.CalculateStats(auditedRecords, dupCount)

	if *verboseFlag {
		printConsoleSummary(stats)
	}

	fmt.Printf("[5/5] Generating 4-Sheet Formatted Excel Workbook: %s\n", *outputFlag)
	rep := reporter.NewExcelReporter()
	if err := rep.GenerateWorkbook(auditedRecords, stats, *outputFlag); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Reporting Error: %v\n", err)
		os.Exit(1)
	}

	elapsed := time.Since(startTime)
	fmt.Println("\n============================================================")
	fmt.Printf("✅ AUDIT COMPLETE in %v\n", elapsed)
	fmt.Printf("📁 Output Workbook: %s\n", *outputFlag)
	fmt.Printf("   • Sheet 1: Summary Dashboard (KPIs & Dept Breakdowns)\n")
	fmt.Printf("   • Sheet 2: Detailed Data (%d standardized records)\n", stats.TotalRecords)
	fmt.Printf("   • Sheet 3: Exceptions (%d total issues logged)\n", stats.ErrorsBySeverity["CRITICAL"]+stats.ErrorsBySeverity["WARNING"]+stats.ErrorsBySeverity["INFO"])
	fmt.Printf("   • Sheet 4: Visual Charts (Doughnut, Bar, Column Charts)\n")
	fmt.Println("============================================================")
}

func printBanner() {
	fmt.Println(`
  ╔═══════════════════════════════════════════════════════════╗
  ║    COMPLIANCE DATA QUALITY & EXCEL REPORT GENERATOR       ║
  ║         ACA / Benefits / Payroll Audit Pipeline           ║
  ╚═══════════════════════════════════════════════════════════╝`)
}

func printConsoleSummary(s models.SummaryStats) {
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("  Total Records:          %d\n", s.TotalRecords)
	fmt.Printf("  Clean Records:          %d\n", s.ValidRecords)
	fmt.Printf("  Records with Errors:    %d (%.2f%% Error Rate)\n", s.ErrorRecords, s.ErrorRate)
	fmt.Printf("  Missing Required Fields: %d\n", s.MissingFieldCount)
	fmt.Printf("  Duplicate Records:      %d\n", s.DuplicateCount)
	fmt.Printf("  Potential FT (>=130h):  %d\n", s.PotentialFTCount)
	fmt.Printf("  Avg Monthly Hours:      %.1f hrs\n", s.AverageMonthlyHours)
	fmt.Println("------------------------------------------------------------")
	fmt.Println("  Exceptions by Severity:")
	fmt.Printf("    • CRITICAL: %d\n", s.ErrorsBySeverity["CRITICAL"])
	fmt.Printf("    • WARNING:  %d\n", s.ErrorsBySeverity["WARNING"])
	fmt.Printf("    • INFO:     %d\n", s.ErrorsBySeverity["INFO"])
	fmt.Println("------------------------------------------------------------")
}
