package reporter

import (
	"fmt"

	"compliance-report-generator/internal/models"

	"github.com/xuri/excelize/v2"
)

func (r *ExcelReporter) buildDetailsSheet(f *excelize.File, sheet string, records []models.EmployeeRecord) error {
	f.SetSheetView(sheet, 0, &excelize.ViewOptions{ShowGridLines: &[]bool{true}[0]})

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "#FFFFFF", Family: "Segoe UI"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#1E3A8A"}, Pattern: 1}, // Navy Blue
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#0F172A", Style: 2},
		},
	})

	cellLeft, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Family: "Segoe UI"},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})

	cellCenter, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Family: "Segoe UI"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	cellRight, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Family: "Segoe UI"},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"},
	})

	currencyStyle, _ := f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Size: 10, Family: "Segoe UI"},
		Alignment:    &excelize.Alignment{Horizontal: "right", Vertical: "center"},
		CustomNumFmt: &[]string{"$#,##0.00"}[0],
	})

	validPillStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "#15803D", Family: "Segoe UI"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#DCFCE7"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	flaggedPillStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "#B91C1C", Family: "Segoe UI"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FEE2E2"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	headers := []string{
		"Record #",
		"Employee ID",
		"First Name",
		"Last Name",
		"SSN / TIN",
		"Department",
		"Employment Status",
		"Hire Date",
		"Monthly Hours",
		"Monthly Gross Pay",
		"Line 14 (Offer Code)",
		"Line 16 (Safe Harbor)",
		"Employee Premium ($)",
		"Validation Status",
		"Anomalies Count",
	}

	cols := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O"}
	for i, h := range headers {
		cell := fmt.Sprintf("%s1", cols[i])
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}
	f.SetRowHeight(sheet, 1, 28)

	for i, rec := range records {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), rec.RecordID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), rec.EmployeeID)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), rec.FirstName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), rec.LastName)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), rec.SSN)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), rec.Department)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), rec.EmploymentStatus)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), rec.HireDate)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), rec.MonthlyHours)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), rec.MonthlyGrossPay)
		f.SetCellValue(sheet, fmt.Sprintf("K%d", row), rec.CoverageOfferCode)
		f.SetCellValue(sheet, fmt.Sprintf("L%d", row), rec.SafeHarborCode)
		f.SetCellValue(sheet, fmt.Sprintf("M%d", row), rec.EmployeeContribution)

		status := "VALID"
		pillStyle := validPillStyle
		if !rec.IsValid || len(rec.Exceptions) > 0 {
			if !rec.IsValid {
				status = "CRITICAL ERROR"
			} else {
				status = "WARNING"
			}
			pillStyle = flaggedPillStyle
		}

		f.SetCellValue(sheet, fmt.Sprintf("N%d", row), status)
		f.SetCellValue(sheet, fmt.Sprintf("O%d", row), len(rec.Exceptions))

		// Apply Cell Styles
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), cellCenter)
		f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), cellCenter)
		f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("D%d", row), cellLeft)
		f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), cellCenter)
		f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("G%d", row), cellLeft)
		f.SetCellStyle(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), cellCenter)
		f.SetCellStyle(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("I%d", row), cellRight)
		f.SetCellStyle(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("J%d", row), currencyStyle)
		f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("L%d", row), cellCenter)
		f.SetCellStyle(sheet, fmt.Sprintf("M%d", row), fmt.Sprintf("M%d", row), currencyStyle)
		f.SetCellStyle(sheet, fmt.Sprintf("N%d", row), fmt.Sprintf("N%d", row), pillStyle)
		f.SetCellStyle(sheet, fmt.Sprintf("O%d", row), fmt.Sprintf("O%d", row), cellCenter)

		f.SetRowHeight(sheet, row, 20)
	}

	// Column Widths
	widths := map[string]float64{
		"A": 10, "B": 14, "C": 16, "D": 16, "E": 16, "F": 22, "G": 20,
		"H": 14, "I": 16, "J": 18, "K": 20, "L": 20, "M": 20, "N": 18, "O": 16,
	}
	for col, w := range widths {
		f.SetColWidth(sheet, col, col, w)
	}

	return nil
}
