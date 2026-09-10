package models

import "time"

// EmployeeRecord represents a standardized employee compliance data record.
type EmployeeRecord struct {
	RecordID             int       `json:"record_id"`
	EmployeeID           string    `json:"employee_id"`
	FirstName            string    `json:"first_name"`
	LastName             string    `json:"last_name"`
	SSN                  string    `json:"ssn"`
	Department           string    `json:"department"`
	EmploymentStatus     string    `json:"employment_status"` // Full-Time, Part-Time, Variable-Hour, Seasonal
	HireDate             string    `json:"hire_date"`
	ParsedHireDate       time.Time `json:"-"`
	MonthlyHours         float64   `json:"monthly_hours"`
	MonthlyGrossPay      float64   `json:"monthly_gross_pay"`
	CoverageOfferCode    string    `json:"coverage_offer_code"` // Line 14: 1A, 1B, 1C, 1D, 1E, 1F, 1G, 1H, 1J, 1K, 1L, etc.
	SafeHarborCode       string    `json:"safe_harbor_code"`    // Line 16: 2A, 2B, 2C, 2D, 2E, 2F, 2G, 2H, 2I
	EmployeeContribution float64   `json:"employee_contribution"`
	IsValid              bool      `json:"is_valid"`
	Exceptions           []Exception `json:"exceptions,omitempty"`
}

// Exception represents a data validation failure or compliance anomaly.
type Exception struct {
	RecordID            int    `json:"record_id"`
	EmployeeID          string `json:"employee_id"`
	EmployeeName        string `json:"employee_name"`
	Field               string `json:"field"`
	Problem             string `json:"problem"`
	Severity            string `json:"severity"` // CRITICAL, WARNING, INFO
	InvalidValue        string `json:"invalid_value"`
	SuggestedCorrection string `json:"suggested_correction"`
}
