package tests

import (
	"testing"

	"compliance-report-generator/internal/reader"
)

func TestNormalizeHeader(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Employee ID", "employee_id"},
		{"emp_id", "employee_id"},
		{"ID #", "employee_id"},
		{"First Name", "first_name"},
		{"LName", "last_name"},
		{"SSN #", "ssn"},
		{"Social_Security_Number", "ssn"},
		{"Department Name", "department"},
		{"Cost_Center", "department"},
		{"Status", "employment_status"},
		{"Hire Date", "hire_date"},
		{"Monthly Hours", "monthly_hours"},
		{"Gross Pay ($)", "monthly_gross_pay"},
		{"Offer Code (Line 14)", "coverage_offer_code"},
		{"Safe Harbor (Line 16)", "safe_harbor_code"},
		{"Employee Premium Share", "employee_contribution"},
	}

	for _, tt := range tests {
		actual := reader.NormalizeHeader(tt.input)
		if actual != tt.expected {
			t.Errorf("NormalizeHeader(%q) = %q, expected %q", tt.input, actual, tt.expected)
		}
	}
}

func TestNormalizeSSN(t *testing.T) {
	tests := []struct {
		input            string
		expectedFormat   string
		expectedDigits   string
	}{
		{"123456789", "123-45-6789", "123456789"},
		{"123-45-6789", "123-45-6789", "123456789"},
		{" 123 45 6789 ", "123-45-6789", "123456789"},
		{"123-45", "123-45", "12345"}, // Incomplete
	}

	for _, tt := range tests {
		formatted, digits := reader.NormalizeSSN(tt.input)
		if formatted != tt.expectedFormat || digits != tt.expectedDigits {
			t.Errorf("NormalizeSSN(%q) = (%q, %q), expected (%q, %q)", tt.input, formatted, digits, tt.expectedFormat, tt.expectedDigits)
		}
	}
}

func TestNormalizeDate(t *testing.T) {
	tests := []struct {
		input       string
		expectedISO string
		expectErr   bool
	}{
		{"2023-05-15", "2023-05-15", false},
		{"05/15/2023", "2023-05-15", false},
		{"5/15/2023", "2023-05-15", false},
		{"2023/05/15", "2023-05-15", false},
		{"invalid-date", "invalid-date", true},
	}

	for _, tt := range tests {
		formatted, _, err := reader.NormalizeDate(tt.input)
		if (err != nil) != tt.expectErr {
			t.Errorf("NormalizeDate(%q) error = %v, expectErr %v", tt.input, err, tt.expectErr)
		}
		if formatted != tt.expectedISO {
			t.Errorf("NormalizeDate(%q) formatted = %q, expected %q", tt.input, formatted, tt.expectedISO)
		}
	}
}

func TestNormalizeCurrency(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"$8,450.00", 8450.00},
		{"$120.50", 120.50},
		{"  $3,000 ", 3000.00},
		{"450.75", 450.75},
		{"", 0.0},
	}

	for _, tt := range tests {
		actual := reader.NormalizeCurrency(tt.input)
		if actual != tt.expected {
			t.Errorf("NormalizeCurrency(%q) = %v, expected %v", tt.input, actual, tt.expected)
		}
	}
}

func TestCSVReaderCleanFile(t *testing.T) {
	r := reader.NewCSVReader()
	records, err := r.ReadFile("../testdata/clean_benchmark.csv", 1)
	if err != nil {
		t.Fatalf("Failed reading clean benchmark: %v", err)
	}

	if len(records) != 8 {
		t.Errorf("Expected 8 records, got %d", len(records))
	}

	first := records[0]
	if first.EmployeeID != "EMP101" || first.FirstName != "Alexander" || first.MonthlyHours != 160.0 {
		t.Errorf("Unexpected record mapping: %+v", first)
	}
}
