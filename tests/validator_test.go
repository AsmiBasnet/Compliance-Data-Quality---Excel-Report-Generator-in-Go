package tests

import (
	"testing"
	"time"

	"compliance-report-generator/internal/models"
	"compliance-report-generator/internal/validator"
)

func TestIsValidSSN(t *testing.T) {
	tests := []struct {
		ssn      string
		expected bool
	}{
		{"452-88-1920", true},
		{"554128874", true},
		{"123-45-6789", false}, // dummy sequential series
		{"000-12-3456", false}, // 000 area
		{"666-45-7890", false}, // 666 area
		{"999-99-9999", false}, // 999 area / dummy
		{"123-00-6789", false}, // 00 group
		{"123-45-0000", false}, // 0000 serial
		{"123-45", false},      // short
	}

	for _, tt := range tests {
		valid, _ := validator.IsValidSSN(tt.ssn)
		if valid != tt.expected {
			t.Errorf("IsValidSSN(%q) = %v, expected %v", tt.ssn, valid, tt.expected)
		}
	}
}

func TestValidationEngine_CriticalErrors(t *testing.T) {
	valEngine := validator.NewValidationEngine()

	// Missing SSN and Employee ID
	rec := models.EmployeeRecord{
		RecordID:         1,
		EmployeeID:       "",
		FirstName:        "Jane",
		LastName:         "Doe",
		SSN:              "",
		HireDate:         "2022-01-01",
		ParsedHireDate:   time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		MonthlyHours:     160.0,
		MonthlyGrossPay:  5000.0,
		EmploymentStatus: "Full-Time",
		CoverageOfferCode: "1E",
		SafeHarborCode:   "2C",
		IsValid:          true,
	}

	valEngine.ValidateRecord(&rec)

	if rec.IsValid {
		t.Errorf("Expected record with missing ID and SSN to be marked invalid")
	}

	if len(rec.Exceptions) < 2 {
		t.Errorf("Expected at least 2 exceptions for missing fields, got %d", len(rec.Exceptions))
	}
}

func TestValidationEngine_FTEBreachWarning(t *testing.T) {
	valEngine := validator.NewValidationEngine()

	// Part-time employee working 145 hours
	rec := models.EmployeeRecord{
		RecordID:         2,
		EmployeeID:       "EMP050",
		FirstName:        "John",
		LastName:         "Smith",
		SSN:              "452-88-1920",
		HireDate:         "2023-01-01",
		ParsedHireDate:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		MonthlyHours:     145.0,
		MonthlyGrossPay:  3500.0,
		EmploymentStatus: "Part-Time",
		CoverageOfferCode: "1H",
		SafeHarborCode:   "2B",
		IsValid:          true,
	}

	valEngine.ValidateRecord(&rec)

	// Warnings don't invalidate record, but are captured in exceptions
	hasFTEWarning := false
	for _, ex := range rec.Exceptions {
		if ex.Field == "MonthlyHours" && ex.Severity == "WARNING" {
			hasFTEWarning = true
			break
		}
	}

	if !hasFTEWarning {
		t.Errorf("Expected FTE threshold breach warning for PT employee with >= 130 hours")
	}
}
