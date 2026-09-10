package processor

import (
	"math"
	"strings"

	"compliance-report-generator/internal/models"
)

// StatsProcessor calculates summary metrics and category groupings.
type StatsProcessor struct{}

// NewStatsProcessor creates a new stats processor instance.
func NewStatsProcessor() *StatsProcessor {
	return &StatsProcessor{}
}

// CalculateStats processes records and generates a SummaryStats object.
func (p *StatsProcessor) CalculateStats(records []models.EmployeeRecord, dupCount int) models.SummaryStats {
	stats := models.SummaryStats{
		TotalRecords:     len(records),
		DuplicateCount:   dupCount,
		ErrorsByType:     make(map[string]int),
		ErrorsBySeverity: make(map[string]int),
		DepartmentStats:  make(map[string]*models.DepartmentStat),
	}

	var totalHours float64
	validCount := 0
	errorCount := 0
	missingCount := 0
	potentialFT := 0

	for _, rec := range records {
		totalHours += rec.MonthlyHours
		if rec.MonthlyHours >= 130.0 {
			potentialFT++
		}

		deptName := strings.TrimSpace(rec.Department)
		if deptName == "" {
			deptName = "Unassigned / General"
		}

		dept, exists := stats.DepartmentStats[deptName]
		if !exists {
			dept = &models.DepartmentStat{
				Department: deptName,
			}
			stats.DepartmentStats[deptName] = dept
		}
		dept.TotalCount++
		dept.AverageHours += rec.MonthlyHours

		if rec.IsValid {
			validCount++
			dept.ValidCount++
		} else {
			errorCount++
			dept.ErrorCount++
		}

		for _, ex := range rec.Exceptions {
			stats.ErrorsBySeverity[ex.Severity]++
			stats.ErrorsByType[ex.Field]++

			if ex.InvalidValue == "[EMPTY]" || strings.Contains(strings.ToLower(ex.Problem), "missing") {
				missingCount++
			}
		}
	}

	stats.ValidRecords = validCount
	stats.ErrorRecords = errorCount
	stats.MissingFieldCount = missingCount
	stats.TotalMonthlyHours = math.Round(totalHours*100) / 100
	stats.PotentialFTCount = potentialFT

	if stats.TotalRecords > 0 {
		stats.AverageMonthlyHours = math.Round((totalHours/float64(stats.TotalRecords))*100) / 100
		stats.ErrorRate = math.Round((float64(errorCount)/float64(stats.TotalRecords))*10000) / 100 // e.g. 15.25%
	}

	// Calculate department rates & averages
	for _, dept := range stats.DepartmentStats {
		if dept.TotalCount > 0 {
			dept.AverageHours = math.Round((dept.AverageHours/float64(dept.TotalCount))*100) / 100
			dept.ErrorRate = math.Round((float64(dept.ErrorCount)/float64(dept.TotalCount))*10000) / 100
		}
	}

	return stats
}
