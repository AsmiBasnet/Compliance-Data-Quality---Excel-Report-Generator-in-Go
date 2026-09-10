package reporter

import (
	"fmt"
	"sort"
	"time"

	"compliance-report-generator/internal/models"

	"github.com/xuri/excelize/v2"
)

func (r *ExcelReporter) buildSummarySheet(f *excelize.File, sheet string, stats models.SummaryStats) error {
	// Set sheet gridlines visible
	f.SetSheetView(sheet, 0, &excelize.ViewOptions{ShowGridLines: &[]bool{true}[0]})

	// Styles
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 18, Color: "#FFFFFF", Family: "Segoe UI"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#1E293B"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", Indent: 1},
	})

	subtitleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Italic: true, Size: 10, Color: "#94A3B8", Family: "Segoe UI"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#1E293B"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", Indent: 1},
	})

	sectionHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 12, Color: "#0F172A", Family: "Segoe UI"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E2E8F0"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#94A3B8", Style: 2},
		},
	})

	tableHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "#FFFFFF", Family: "Segoe UI"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#334155"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#1E293B", Style: 1},
			{Type: "top", Color: "#1E293B", Style: 1},
			{Type: "left", Color: "#1E293B", Style: 1},
			{Type: "right", Color: "#1E293B", Style: 1},
		},
	})

	kpiLabelStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "#475569", Family: "Segoe UI"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#F8FAFC"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", Indent: 1},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#E2E8F0", Style: 1},
			{Type: "left", Color: "#E2E8F0", Style: 1},
		},
	})

	kpiValueStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "#0F172A", Family: "Segoe UI"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFFFFF"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#E2E8F0", Style: 1},
			{Type: "right", Color: "#E2E8F0", Style: 1},
		},
	})

	kpiAlertStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "#DC2626", Family: "Segoe UI"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FEF2F2"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#E2E8F0", Style: 1},
			{Type: "right", Color: "#E2E8F0", Style: 1},
		},
	})

	kpiSuccessStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "#16A34A", Family: "Segoe UI"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#F0FDF4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#E2E8F0", Style: 1},
			{Type: "right", Color: "#E2E8F0", Style: 1},
		},
	})

	cellLeft, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Family: "Segoe UI"},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#E2E8F0", Style: 1},
			{Type: "top", Color: "#E2E8F0", Style: 1},
			{Type: "left", Color: "#E2E8F0", Style: 1},
			{Type: "right", Color: "#E2E8F0", Style: 1},
		},
	})

	cellCenter, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Family: "Segoe UI"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#E2E8F0", Style: 1},
			{Type: "top", Color: "#E2E8F0", Style: 1},
			{Type: "left", Color: "#E2E8F0", Style: 1},
			{Type: "right", Color: "#E2E8F0", Style: 1},
		},
	})

	cellRight, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Family: "Segoe UI"},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#E2E8F0", Style: 1},
			{Type: "top", Color: "#E2E8F0", Style: 1},
			{Type: "left", Color: "#E2E8F0", Style: 1},
			{Type: "right", Color: "#E2E8F0", Style: 1},
		},
	})

	// Banner
	f.MergeCell(sheet, "A1", "G1")
	f.SetCellValue(sheet, "A1", "ACA COMPLIANCE & DATA QUALITY AUDIT REPORT")
	f.SetCellStyle(sheet, "A1", "G1", titleStyle)
	f.SetRowHeight(sheet, 1, 35)

	f.MergeCell(sheet, "A2", "G2")
	f.SetCellValue(sheet, "A2", fmt.Sprintf("Generated: %s | System: Compliance Data Quality Engine v1.0", time.Now().Format("2006-01-02 15:04:05 MST")))
	f.SetCellStyle(sheet, "A2", "G2", subtitleStyle)
	f.SetRowHeight(sheet, 2, 20)

	// Section 1: Executive KPI Summary
	f.MergeCell(sheet, "A4", "C4")
	f.SetCellValue(sheet, "A4", "📊 EXECUTIVE SUMMARY METRICS")
	f.SetCellStyle(sheet, "A4", "C4", sectionHeaderStyle)

	kpis := []struct {
		Label string
		Value string
		Style int
	}{
		{"Total Ingested Records", fmt.Sprintf("%d", stats.TotalRecords), kpiValueStyle},
		{"Clean / Valid Records", fmt.Sprintf("%d", stats.ValidRecords), kpiSuccessStyle},
		{"Records with Data Quality Errors", fmt.Sprintf("%d", stats.ErrorRecords), kpiAlertStyle},
		{"Missing Required Fields Count", fmt.Sprintf("%d", stats.MissingFieldCount), kpiAlertStyle},
		{"Duplicate Records Detected", fmt.Sprintf("%d", stats.DuplicateCount), kpiAlertStyle},
		{"Overall Error Rate", fmt.Sprintf("%.2f%%", stats.ErrorRate), kpiAlertStyle},
		{"Potential Full-Time Employees (>= 130 hrs)", fmt.Sprintf("%d", stats.PotentialFTCount), kpiValueStyle},
		{"Average Monthly Hours Worked", fmt.Sprintf("%.1f hrs", stats.AverageMonthlyHours), kpiValueStyle},
		{"Total Monthly Payroll Hours", fmt.Sprintf("%.1f hrs", stats.TotalMonthlyHours), kpiValueStyle},
	}

	row := 5
	for _, kpi := range kpis {
		f.MergeCell(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("B%d", row))
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), kpi.Label)
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("B%d", row), kpiLabelStyle)

		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), kpi.Value)
		f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), kpi.Style)
		f.SetRowHeight(sheet, row, 22)
		row++
	}

	// Section 2: Severity Breakdown
	f.MergeCell(sheet, "E4", "G4")
	f.SetCellValue(sheet, "E4", "⚠️ EXCEPTIONS BY SEVERITY")
	f.SetCellStyle(sheet, "E4", "G4", sectionHeaderStyle)

	severities := []struct {
		Severity string
		Count    int
		Desc     string
	}{
		{"CRITICAL", stats.ErrorsBySeverity["CRITICAL"], "Missing SSN/ID, negative hours/pay, invalid codes"},
		{"WARNING", stats.ErrorsBySeverity["WARNING"], "130+ hr PT breaches, FPL affordability threshold, 744+ hr"},
		{"INFO", stats.ErrorsBySeverity["INFO"], "Missing optional fields (e.g. Department)"},
	}

	sRow := 5
	for _, s := range severities {
		f.SetCellValue(sheet, fmt.Sprintf("E%d", sRow), s.Severity)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", sRow), s.Count)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", sRow), s.Desc)

		f.SetCellStyle(sheet, fmt.Sprintf("E%d", sRow), fmt.Sprintf("E%d", sRow), kpiLabelStyle)
		f.SetCellStyle(sheet, fmt.Sprintf("F%d", sRow), fmt.Sprintf("F%d", sRow), kpiValueStyle)
		f.SetCellStyle(sheet, fmt.Sprintf("G%d", sRow), fmt.Sprintf("G%d", sRow), cellLeft)
		f.SetRowHeight(sheet, sRow, 22)
		sRow++
	}

	// Section 3: Department Breakdown Table
	deptStartRow := 16
	f.MergeCell(sheet, fmt.Sprintf("A%d", deptStartRow), fmt.Sprintf("G%d", deptStartRow))
	f.SetCellValue(sheet, fmt.Sprintf("A%d", deptStartRow), "🏢 DATA QUALITY & COMPLIANCE BY DEPARTMENT")
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", deptStartRow), fmt.Sprintf("G%d", deptStartRow), sectionHeaderStyle)
	f.SetRowHeight(sheet, deptStartRow, 24)

	headers := []string{"Department Name", "Total Records", "Valid Records", "Error Records", "Error Rate", "Avg Monthly Hours", "Compliance Status"}
	headerRow := deptStartRow + 1
	cols := []string{"A", "B", "C", "D", "E", "F", "G"}
	for i, h := range headers {
		cell := fmt.Sprintf("%s%d", cols[i], headerRow)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, tableHeaderStyle)
	}
	f.SetRowHeight(sheet, headerRow, 24)

	// Sort departments alphabetically
	var deptNames []string
	for d := range stats.DepartmentStats {
		deptNames = append(deptNames, d)
	}
	sort.Strings(deptNames)

	dRow := headerRow + 1
	for _, dName := range deptNames {
		d := stats.DepartmentStats[dName]
		status := "HEALTHY"
		if d.ErrorRate > 15.0 {
			status = "ACTION REQUIRED"
		} else if d.ErrorRate > 0.0 {
			status = "ATTENTION"
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", dRow), d.Department)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", dRow), d.TotalCount)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", dRow), d.ValidCount)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", dRow), d.ErrorCount)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", dRow), fmt.Sprintf("%.2f%%", d.ErrorRate))
		f.SetCellValue(sheet, fmt.Sprintf("F%d", dRow), fmt.Sprintf("%.1f", d.AverageHours))
		f.SetCellValue(sheet, fmt.Sprintf("G%d", dRow), status)

		f.SetCellStyle(sheet, fmt.Sprintf("A%d", dRow), fmt.Sprintf("A%d", dRow), cellLeft)
		f.SetCellStyle(sheet, fmt.Sprintf("B%d", dRow), fmt.Sprintf("D%d", dRow), cellRight)
		f.SetCellStyle(sheet, fmt.Sprintf("E%d", dRow), fmt.Sprintf("F%d", dRow), cellRight)
		f.SetCellStyle(sheet, fmt.Sprintf("G%d", dRow), fmt.Sprintf("G%d", dRow), cellCenter)
		f.SetRowHeight(sheet, dRow, 20)
		dRow++
	}

	// Auto-fit column widths
	f.SetColWidth(sheet, "A", "A", 28)
	f.SetColWidth(sheet, "B", "B", 18)
	f.SetColWidth(sheet, "C", "C", 16)
	f.SetColWidth(sheet, "D", "D", 16)
	f.SetColWidth(sheet, "E", "E", 16)
	f.SetColWidth(sheet, "F", "F", 18)
	f.SetColWidth(sheet, "G", "G", 22)

	return nil
}
