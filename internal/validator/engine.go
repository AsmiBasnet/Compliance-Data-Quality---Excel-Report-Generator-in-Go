package validator

import (
	"fmt"
	"strings"

	"compliance-report-generator/internal/models"
)

// ValidationEngine runs multi-layer data quality & compliance audits on records.
type ValidationEngine struct{}

// NewValidationEngine creates a new validator instance.
func NewValidationEngine() *ValidationEngine {
	return &ValidationEngine{}
}

// ValidateRecords evaluates all employee records and records anomalies as Exceptions.
func (v *ValidationEngine) ValidateRecords(records []models.EmployeeRecord) []models.EmployeeRecord {
	for i := range records {
		v.ValidateRecord(&records[i])
	}
	return records
}

// ValidateRecord evaluates a single employee record.
func (v *ValidationEngine) ValidateRecord(rec *models.EmployeeRecord) {
	rec.Exceptions = make([]models.Exception, 0)
	fullName := strings.TrimSpace(rec.FirstName + " " + rec.LastName)
	if fullName == "" {
		fullName = "Unknown Employee"
	}

	addEx := func(field, problem, severity, invalidVal, suggestion string) {
		rec.Exceptions = append(rec.Exceptions, models.Exception{
			RecordID:            rec.RecordID,
			EmployeeID:          rec.EmployeeID,
			EmployeeName:        fullName,
			Field:               field,
			Problem:             problem,
			Severity:            severity,
			InvalidValue:        invalidVal,
			SuggestedCorrection: suggestion,
		})
		if severity == "CRITICAL" {
			rec.IsValid = false
		}
	}

	// 1. Employee ID
	if strings.TrimSpace(rec.EmployeeID) == "" {
		addEx("EmployeeID", "Missing Employee ID", "CRITICAL", "[EMPTY]", "Provide valid unique employee ID from HRIS roster")
	}

	// 2. Names
	if strings.TrimSpace(rec.FirstName) == "" {
		addEx("FirstName", "Missing First Name", "CRITICAL", "[EMPTY]", "Enter employee first name")
	}
	if strings.TrimSpace(rec.LastName) == "" {
		addEx("LastName", "Missing Last Name", "CRITICAL", "[EMPTY]", "Enter employee last name")
	}

	// 3. SSN / TIN
	if strings.TrimSpace(rec.SSN) == "" {
		addEx("SSN", "Missing SSN / Taxpayer Identification Number", "CRITICAL", "[EMPTY]", "Provide valid 9-digit SSN (XXX-XX-XXXX) or IRS ITIN")
	} else if valid, reason := IsValidSSN(rec.SSN); !valid {
		addEx("SSN", fmt.Sprintf("Invalid SSN format: %s", reason), "CRITICAL", rec.SSN, "Correct SSN digits to conform to SSA/IRS formatting guidelines")
	}

	// 4. Department
	if strings.TrimSpace(rec.Department) == "" {
		addEx("Department", "Missing Department Assignment", "INFO", "[EMPTY]", "Assign to valid department / cost center (e.g., Operations, IT, Sales)")
	}

	// 5. Hire Date
	if valid, reason := IsValidDate(rec.ParsedHireDate, rec.HireDate); !valid {
		addEx("HireDate", fmt.Sprintf("Invalid Hire Date: %s", reason), "CRITICAL", rec.HireDate, "Enter date in standard YYYY-MM-DD or MM/DD/YYYY format")
	}

	// 6. Monthly Hours & ACA Full-Time (FTE) Thresholds
	if rec.MonthlyHours < 0 {
		addEx("MonthlyHours", "Negative Monthly Hours Logged", "CRITICAL", fmt.Sprintf("%.2f", rec.MonthlyHours), "Ensure hours worked >= 0 or adjust payroll import")
	} else if rec.MonthlyHours > 744 {
		addEx("MonthlyHours", "Monthly Hours Exceed Calendar Month Maximum (744 hrs)", "WARNING", fmt.Sprintf("%.2f", rec.MonthlyHours), "Verify timesheet for duplicate shift submissions")
	}

	// ACA FTE Look-back / Monthly Threshold Alert
	statusLower := strings.ToLower(rec.EmploymentStatus)
	if (strings.Contains(statusLower, "part") || strings.Contains(statusLower, "variable")) && rec.MonthlyHours >= 130.0 {
		addEx("MonthlyHours", "ACA FTE Threshold Breach: Part-time/Variable employee exceeded 130 hrs/month (30 hrs/wk)", "WARNING", fmt.Sprintf("%.2f hrs (%s)", rec.MonthlyHours, rec.EmploymentStatus), "Review ACA Look-Back Measurement and evaluate mandatory health coverage offer")
	}

	// 7. Monthly Gross Pay
	if rec.MonthlyGrossPay < 0 {
		addEx("MonthlyGrossPay", "Negative Monthly Gross Pay", "CRITICAL", fmt.Sprintf("$%.2f", rec.MonthlyGrossPay), "Correct payroll gross pay to positive balance")
	}

	// 8. Line 14 Offer of Coverage Code
	if rec.CoverageOfferCode == "" {
		addEx("CoverageOfferCode", "Missing Form 1095-C Line 14 Offer Code", "WARNING", "[EMPTY]", "Specify IRS Line 14 Code (e.g., 1A Qualifying Offer, 1E MEC to Family, 1H No Offer)")
	} else if !IsValidOfferCode(rec.CoverageOfferCode) {
		addEx("CoverageOfferCode", fmt.Sprintf("Unrecognized Line 14 Offer Code '%s'", rec.CoverageOfferCode), "CRITICAL", rec.CoverageOfferCode, "Must be valid IRS code: 1A, 1B, 1C, 1D, 1E, 1F, 1G, 1H, 1J, 1K, 1L, etc.")
	}

	// 9. Line 16 Safe Harbor / Relief Code
	if rec.SafeHarborCode != "" && !IsValidSafeHarborCode(rec.SafeHarborCode) {
		addEx("SafeHarborCode", fmt.Sprintf("Unrecognized Section 4980H Safe Harbor Code '%s'", rec.SafeHarborCode), "WARNING", rec.SafeHarborCode, "Use valid code: 2A (Not Employed), 2B (Part-Time), 2C (Enrolled), 2F (W-2), 2G (FPL), 2H (Rate of Pay)")
	}

	// 10. Employee Contribution & Affordability Rules
	if rec.EmployeeContribution < 0 {
		addEx("EmployeeContribution", "Negative Employee Contribution Amount", "CRITICAL", fmt.Sprintf("$%.2f", rec.EmployeeContribution), "Set contribution >= $0.00")
	}

	// Affordability check: Qualifying Offer 1A requires lowest cost employee share <= FPL monthly limit (~$108-$115)
	if rec.CoverageOfferCode == "1A" && rec.EmployeeContribution > 115.00 {
		addEx("EmployeeContribution", "Code 1A (Qualifying Offer) Exceeds Federal Poverty Line Affordability Threshold", "WARNING", fmt.Sprintf("$%.2f", rec.EmployeeContribution), "Qualifying Offer (1A) requires employee cost <= 100% FPL safe harbor (~$108/mo)")
	}

	if len(rec.Exceptions) > 0 && rec.IsValid {
		// If there are only warnings, record is still valid for export, but logged
	}
}
