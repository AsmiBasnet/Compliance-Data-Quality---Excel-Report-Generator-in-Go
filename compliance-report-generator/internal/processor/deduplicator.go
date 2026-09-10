package processor

import (
	"fmt"
	"strings"

	"compliance-report-generator/internal/models"
)

// Deduplicator scans records for duplicate employee IDs, SSNs, or conflicting demographic records.
type Deduplicator struct{}

// NewDeduplicator creates a new deduplicator instance.
func NewDeduplicator() *Deduplicator {
	return &Deduplicator{}
}

// ProcessDuplicates flags duplicate rows and records exceptions.
func (d *Deduplicator) ProcessDuplicates(records []models.EmployeeRecord) ([]models.EmployeeRecord, int) {
	seenEmpIDs := make(map[string]int) // EmpID -> first RecordID seen
	seenSSNs := make(map[string]int)   // SSN -> first RecordID seen
	duplicateCount := 0

	for i := range records {
		rec := &records[i]
		isDup := false

		empKey := strings.ToUpper(strings.TrimSpace(rec.EmployeeID))
		if empKey != "" {
			if origID, exists := seenEmpIDs[empKey]; exists {
				rec.Exceptions = append(rec.Exceptions, models.Exception{
					RecordID:            rec.RecordID,
					EmployeeID:          rec.EmployeeID,
					EmployeeName:        rec.FirstName + " " + rec.LastName,
					Field:               "EmployeeID",
					Problem:             fmt.Sprintf("Duplicate Employee ID found (Original in Record #%d)", origID),
					Severity:            "CRITICAL",
					InvalidValue:        rec.EmployeeID,
					SuggestedCorrection: "Deduplicate payroll feed or verify unique employee badge/roster identifier",
				})
				rec.IsValid = false
				isDup = true
			} else {
				seenEmpIDs[empKey] = rec.RecordID
			}
		}

		cleanSSN := strings.ReplaceAll(rec.SSN, "-", "")
		if cleanSSN != "" && len(cleanSSN) == 9 {
			if origID, exists := seenSSNs[cleanSSN]; exists {
				rec.Exceptions = append(rec.Exceptions, models.Exception{
					RecordID:            rec.RecordID,
					EmployeeID:          rec.EmployeeID,
					EmployeeName:        rec.FirstName + " " + rec.LastName,
					Field:               "SSN",
					Problem:             fmt.Sprintf("Duplicate SSN/TIN assigned to multiple records (First seen in Record #%d)", origID),
					Severity:            "CRITICAL",
					InvalidValue:        rec.SSN,
					SuggestedCorrection: "Verify employee tax identification documentation for SSN collision",
				})
				rec.IsValid = false
				isDup = true
			} else {
				seenSSNs[cleanSSN] = rec.RecordID
			}
		}

		if isDup {
			duplicateCount++
		}
	}

	return records, duplicateCount
}
