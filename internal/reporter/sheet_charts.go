package reporter

import (
	"fmt"
	"sort"

	"compliance-report-generator/internal/models"

	"github.com/xuri/excelize/v2"
)

func (r *ExcelReporter) buildChartsSheet(f *excelize.File, sheet string, stats models.SummaryStats) error {
	f.SetSheetView(sheet, 0, &excelize.ViewOptions{ShowGridLines: &[]bool{true}[0]})

	// Title
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 16, Color: "#FFFFFF", Family: "Segoe UI"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#0F172A"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", Indent: 1},
	})
	f.MergeCell(sheet, "A1", "P1")
	f.SetCellValue(sheet, "A1", "AUDIT & COMPLIANCE VISUAL ANALYTICS DASHBOARD")
	f.SetCellStyle(sheet, "A1", "P1", titleStyle)
	f.SetRowHeight(sheet, 1, 32)

	tableHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "#FFFFFF", Family: "Segoe UI"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#334155"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	cellLeft, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 9, Family: "Segoe UI"},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})
	cellRight, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 9, Family: "Segoe UI"},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"},
	})

	// 1. Data Table for Chart 1: Status Breakdown (Cols AA - AB to keep charts clean, or in rows 35+)
	f.SetCellValue(sheet, "AA1", "Status")
	f.SetCellValue(sheet, "AB1", "Count")
	f.SetCellStyle(sheet, "AA1", "AB1", tableHeaderStyle)

	f.SetCellValue(sheet, "AA2", "Clean Records")
	f.SetCellValue(sheet, "AB2", stats.ValidRecords)
	f.SetCellValue(sheet, "AA3", "Records with Errors")
	f.SetCellValue(sheet, "AB3", stats.ErrorRecords)
	f.SetCellStyle(sheet, "AA2", "AA3", cellLeft)
	f.SetCellStyle(sheet, "AB2", "AB3", cellRight)

	// Chart 1: Donut / Doughnut Chart: Records by Status
	statusChart := &excelize.Chart{
		Type: excelize.Doughnut,
		Series: []excelize.ChartSeries{
			{
				Name:       "Status Distribution",
				Categories: fmt.Sprintf("'%s'!$AA$2:$AA$3", sheet),
				Values:     fmt.Sprintf("'%s'!$AB$2:$AB$3", sheet),
			},
		},
		Title: []excelize.RichTextRun{
			{Text: "Record Compliance Status (Clean vs Flagged)"},
		},
		PlotArea: excelize.ChartPlotArea{
			ShowBubbleSize:  false,
			ShowCatName:     false,
			ShowLeaderLines: false,
			ShowPercent:     true,
			ShowVal:         true,
		},
		Dimension: excelize.ChartDimension{
			Width:  520,
			Height: 320,
		},
	}
	if err := f.AddChart(sheet, "A3", statusChart); err != nil {
		return err
	}

	// 2. Data Table for Chart 2: Top Error Fields
	f.SetCellValue(sheet, "AD1", "Field")
	f.SetCellValue(sheet, "AE1", "Errors")
	f.SetCellStyle(sheet, "AD1", "AE1", tableHeaderStyle)

	type kv struct {
		Key string
		Val int
	}
	var errList []kv
	for k, v := range stats.ErrorsByType {
		errList = append(errList, kv{k, v})
	}
	sort.Slice(errList, func(i, j int) bool {
		return errList[i].Val > errList[j].Val
	})

	if len(errList) == 0 {
		errList = append(errList, kv{"No Errors Detected", 0})
	}

	for i, item := range errList {
		r := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("AD%d", r), item.Key)
		f.SetCellValue(sheet, fmt.Sprintf("AE%d", r), item.Val)
		f.SetCellStyle(sheet, fmt.Sprintf("AD%d", r), fmt.Sprintf("AD%d", r), cellLeft)
		f.SetCellStyle(sheet, fmt.Sprintf("AE%d", r), fmt.Sprintf("AE%d", r), cellRight)
	}

	maxFieldRow := len(errList) + 1
	if maxFieldRow < 2 {
		maxFieldRow = 2
	}

	// Chart 2: Bar / Column Chart: Exceptions by Field
	fieldChart := &excelize.Chart{
		Type: excelize.Bar,
		Series: []excelize.ChartSeries{
			{
				Name:       "Exception Count",
				Categories: fmt.Sprintf("'%s'!$AD$2:$AD$%d", sheet, maxFieldRow),
				Values:     fmt.Sprintf("'%s'!$AE$2:$AE$%d", sheet, maxFieldRow),
			},
		},
		Title: []excelize.RichTextRun{
			{Text: "Data Quality Exceptions by Field"},
		},
		Dimension: excelize.ChartDimension{
			Width:  520,
			Height: 320,
		},
	}
	if err := f.AddChart(sheet, "I3", fieldChart); err != nil {
		return err
	}

	// 3. Data Table for Chart 3: Department Error Rate
	f.SetCellValue(sheet, "AG1", "Department")
	f.SetCellValue(sheet, "AH1", "Error Rate (%)")
	f.SetCellStyle(sheet, "AG1", "AH1", tableHeaderStyle)

	var deptKeys []string
	for d := range stats.DepartmentStats {
		deptKeys = append(deptKeys, d)
	}
	sort.Strings(deptKeys)

	dRow := 2
	for _, d := range deptKeys {
		stat := stats.DepartmentStats[d]
		f.SetCellValue(sheet, fmt.Sprintf("AG%d", dRow), stat.Department)
		f.SetCellValue(sheet, fmt.Sprintf("AH%d", dRow), stat.ErrorRate)
		f.SetCellStyle(sheet, fmt.Sprintf("AG%d", dRow), fmt.Sprintf("AG%d", dRow), cellLeft)
		f.SetCellStyle(sheet, fmt.Sprintf("AH%d", dRow), fmt.Sprintf("AH%d", dRow), cellRight)
		dRow++
	}
	maxDeptRow := dRow - 1
	if maxDeptRow < 2 {
		maxDeptRow = 2
	}

	// Chart 3: Column Chart: Error Rate by Department
	deptChart := &excelize.Chart{
		Type: excelize.Col,
		Series: []excelize.ChartSeries{
			{
				Name:       "Department Error Rate (%)",
				Categories: fmt.Sprintf("'%s'!$AG$2:$AG$%d", sheet, maxDeptRow),
				Values:     fmt.Sprintf("'%s'!$AH$2:$AH$%d", sheet, maxDeptRow),
			},
		},
		Title: []excelize.RichTextRun{
			{Text: "Data Quality Error Rate by Department (%)"},
		},
		Dimension: excelize.ChartDimension{
			Width:  1060,
			Height: 340,
		},
	}
	if err := f.AddChart(sheet, "A20", deptChart); err != nil {
		return err
	}

	return nil
}
