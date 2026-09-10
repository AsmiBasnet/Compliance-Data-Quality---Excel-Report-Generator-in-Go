package reporter

import (
	"fmt"

	"compliance-report-generator/internal/models"

	"github.com/xuri/excelize/v2"
)

func (r *ExcelReporter) buildExceptionsSheet(f *excelize.File, sheet string, records []models.EmployeeRecord) error {
	f.SetSheetView(sheet, 0, &excelize.ViewOptions{ShowGridLines: &[]bool{true}[0]})

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "#FFFFFF", Family: "Segoe UI"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#881337"}, Pattern: 1}, // Deep Burgundy
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#4C0519", Style: 2},
		},
	})

	cellLeft, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Family: "Segoe UI"},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#E2E8F0", Style: 1},
			{Type: "left", Color: "#E2E8F0", Style: 1},
			{Type: "right", Color: "#E2E8F0", Style: 1},
		},
	})

	cellCenter, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Family: "Segoe UI"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#E2E8F0", Style: 1},
			{Type: "left", Color: "#E2E8F0", Style: 1},
			{Type: "right", Color: "#E2E8F0", Style: 1},
		},
	})

	criticalStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "#991B1B", Family: "Segoe UI"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FEE2E2"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#E2E8F0", Style: 1},
			{Type: "left", Color: "#E2E8F0", Style: 1},
			{Type: "right", Color: "#E2E8F0", Style: 1},
		},
	})

	warningStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "#92400E", Family: "Segoe UI"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FEF3C7"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#E2E8F0", Style: 1},
			{Type: "left", Color: "#E2E8F0", Style: 1},
			{Type: "right", Color: "#E2E8F0", Style: 1},
		},
	})

	infoStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "#1E40AF", Family: "Segoe UI"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#DBEAFE"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#E2E8F0", Style: 1},
			{Type: "left", Color: "#E2E8F0", Style: 1},
			{Type: "right", Color: "#E2E8F0", Style: 1},
		},
	})

	headers := []string{
		"Exception #",
		"Record ID",
		"Employee ID",
		"Employee Name",
		"Field Name",
		"Severity",
		"Problem Description",
		"Invalid / Detected Value",
		"Suggested Action & Correction",
	}

	cols := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}
	for i, h := range headers {
		cell := fmt.Sprintf("%s1", cols[i])
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}
	f.SetRowHeight(sheet, 1, 28)

	exCounter := 1
	currRow := 2

	for _, rec := range records {
		for _, ex := range rec.Exceptions {
			f.SetCellValue(sheet, fmt.Sprintf("A%d", currRow), exCounter)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", currRow), ex.RecordID)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", currRow), ex.EmployeeID)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", currRow), ex.EmployeeName)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", currRow), ex.Field)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", currRow), ex.Severity)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", currRow), ex.Problem)
			f.SetCellValue(sheet, fmt.Sprintf("H%d", currRow), ex.InvalidValue)
			f.SetCellValue(sheet, fmt.Sprintf("I%d", currRow), ex.SuggestedCorrection)

			sevStyle := infoStyle
			switch ex.Severity {
			case "CRITICAL":
				sevStyle = criticalStyle
			case "WARNING":
				sevStyle = warningStyle
			}

			f.SetCellStyle(sheet, fmt.Sprintf("A%d", currRow), fmt.Sprintf("C%d", currRow), cellCenter)
			f.SetCellStyle(sheet, fmt.Sprintf("D%d", currRow), fmt.Sprintf("E%d", currRow), cellLeft)
			f.SetCellStyle(sheet, fmt.Sprintf("F%d", currRow), fmt.Sprintf("F%d", currRow), sevStyle)
			f.SetCellStyle(sheet, fmt.Sprintf("G%d", currRow), fmt.Sprintf("I%d", currRow), cellLeft)

			f.SetRowHeight(sheet, currRow, 22)
			exCounter++
			currRow++
		}
	}

	// Auto-fit widths
	widths := map[string]float64{
		"A": 12, "B": 12, "C": 14, "D": 22, "E": 20, "F": 14, "G": 40, "H": 26, "I": 46,
	}
	for col, w := range widths {
		f.SetColWidth(sheet, col, col, w)
	}

	return nil
}
