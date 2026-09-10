package reporter

import (
	"fmt"

	"compliance-report-generator/internal/models"

	"github.com/xuri/excelize/v2"
)

// ExcelReporter orchestrates multi-sheet workbook generation with styles and charts.
type ExcelReporter struct {
	file *excelize.File
}

// NewExcelReporter creates a new reporter instance.
func NewExcelReporter() *ExcelReporter {
	return &ExcelReporter{
		file: excelize.NewFile(),
	}
}

// GenerateWorkbook builds the complete 4-sheet compliance audit workbook.
func (r *ExcelReporter) GenerateWorkbook(records []models.EmployeeRecord, stats models.SummaryStats, outputPath string) error {
	f := r.file

	// Sheet 1: Summary
	summarySheet := "Summary"
	f.SetSheetName("Sheet1", summarySheet)
	if err := r.buildSummarySheet(f, summarySheet, stats); err != nil {
		return fmt.Errorf("error building Summary sheet: %w", err)
	}

	// Sheet 2: Detailed Data
	detailsSheet := "Detailed Data"
	f.NewSheet(detailsSheet)
	if err := r.buildDetailsSheet(f, detailsSheet, records); err != nil {
		return fmt.Errorf("error building Detailed Data sheet: %w", err)
	}

	// Sheet 3: Exceptions
	exceptionsSheet := "Exceptions"
	f.NewSheet(exceptionsSheet)
	if err := r.buildExceptionsSheet(f, exceptionsSheet, records); err != nil {
		return fmt.Errorf("error building Exceptions sheet: %w", err)
	}

	// Sheet 4: Charts
	chartsSheet := "Charts"
	f.NewSheet(chartsSheet)
	if err := r.buildChartsSheet(f, chartsSheet, stats); err != nil {
		return fmt.Errorf("error building Charts sheet: %w", err)
	}

	// Set active tab to Summary
	summaryIdx, _ := f.GetSheetIndex(summarySheet)
	f.SetActiveSheet(summaryIdx)

	// Save to destination
	if err := f.SaveAs(outputPath); err != nil {
		return fmt.Errorf("error saving workbook to %s: %w", outputPath, err)
	}

	return nil
}
