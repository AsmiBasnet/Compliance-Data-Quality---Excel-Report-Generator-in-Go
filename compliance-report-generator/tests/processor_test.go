package tests

import (
	"testing"

	"compliance-report-generator/internal/models"
	"compliance-report-generator/internal/processor"
)

func TestDeduplicator(t *testing.T) {
	deduper := processor.NewDeduplicator()

	records := []models.EmployeeRecord{
		{RecordID: 1, EmployeeID: "EMP001", SSN: "123-45-6789", FirstName: "Alice", LastName: "Smith", IsValid: true},
		{RecordID: 2, EmployeeID: "EMP002", SSN: "987-65-4321", FirstName: "Bob", LastName: "Jones", IsValid: true},
		{RecordID: 3, EmployeeID: "EMP001", SSN: "123-45-6789", FirstName: "Alice", LastName: "Smith", IsValid: true}, // Duplicate
	}

	audited, dupCount := deduper.ProcessDuplicates(records)

	if dupCount != 1 {
		t.Errorf("Expected 1 duplicate count, got %d", dupCount)
	}

	if audited[2].IsValid {
		t.Errorf("Expected duplicate record #3 to be marked invalid")
	}
}

func TestStatsProcessor(t *testing.T) {
	statsProc := processor.NewStatsProcessor()

	records := []models.EmployeeRecord{
		{RecordID: 1, Department: "IT", MonthlyHours: 160, IsValid: true},
		{RecordID: 2, Department: "IT", MonthlyHours: 140, IsValid: false, Exceptions: []models.Exception{{Severity: "CRITICAL", Field: "SSN"}}},
		{RecordID: 3, Department: "Sales", MonthlyHours: 80, IsValid: true},
	}

	stats := statsProc.CalculateStats(records, 0)

	if stats.TotalRecords != 3 {
		t.Errorf("Expected 3 total records, got %d", stats.TotalRecords)
	}
	if stats.ValidRecords != 2 {
		t.Errorf("Expected 2 valid records, got %d", stats.ValidRecords)
	}
	if stats.ErrorRecords != 1 {
		t.Errorf("Expected 1 error record, got %d", stats.ErrorRecords)
	}
	if stats.PotentialFTCount != 2 {
		t.Errorf("Expected 2 potential FT records (>=130 hrs), got %d", stats.PotentialFTCount)
	}

	itDept, ok := stats.DepartmentStats["IT"]
	if !ok || itDept.TotalCount != 2 || itDept.ErrorCount != 1 {
		t.Errorf("Unexpected IT department stats: %+v", itDept)
	}
}
