package validator

import (
	"regexp"
	"strings"
	"time"
)

var (
	ssnPattern = regexp.MustCompile(`^(?:\d{3}-\d{2}-\d{4}|\d{9})$`)
)

// Valid Line 14 Offer of Coverage Codes
var validOfferCodes = map[string]bool{
	"1A": true, "1B": true, "1C": true, "1D": true,
	"1E": true, "1F": true, "1G": true, "1H": true,
	"1J": true, "1K": true, "1L": true, "1M": true,
	"1N": true, "1O": true, "1P": true, "1Q": true,
	"1R": true, "1S": true, "1T": true, "1U": true,
}

// Valid Line 16 Section 4980H Safe Harbor and Relief Codes
var validSafeHarborCodes = map[string]bool{
	"":   true, // Can be blank if no safe harbor applies
	"2A": true, "2B": true, "2C": true, "2D": true,
	"2E": true, "2F": true, "2G": true, "2H": true, "2I": true,
}

// IsValidSSN checks if a string matches SSN syntax and passes SSA issuance sanity checks.
func IsValidSSN(ssn string) (bool, string) {
	clean := strings.ReplaceAll(strings.TrimSpace(ssn), "-", "")
	if len(clean) != 9 {
		return false, "SSN must contain exactly 9 digits"
	}
	if !ssnPattern.MatchString(clean) {
		return false, "SSN contains invalid characters"
	}

	area := clean[0:3]
	group := clean[3:5]
	serial := clean[5:9]

	// SSA invalid rules
	if area == "000" || area == "666" || area >= "900" {
		return false, "Area number cannot be 000, 666, or 900-999"
	}
	if group == "00" {
		return false, "Group number cannot be 00"
	}
	if serial == "0000" {
		return false, "Serial number cannot be 0000"
	}
	if clean == "111111111" || clean == "123456789" || clean == "999999999" {
		return false, "Recognized dummy/test SSN sequence"
	}

	return true, ""
}

// IsValidDate checks if a hire date is logically plausible.
func IsValidDate(t time.Time, raw string) (bool, string) {
	if strings.TrimSpace(raw) == "" {
		return false, "Date cannot be empty"
	}
	if t.IsZero() {
		return false, "Unrecognized date format (expected YYYY-MM-DD or MM/DD/YYYY)"
	}
	if t.Year() < 1950 {
		return false, "Hire date year precedes 1950"
	}
	// Allow up to 1 year in future for future hires, otherwise flag
	oneYearFuture := time.Now().AddDate(1, 0, 0)
	if t.After(oneYearFuture) {
		return false, "Hire date is more than 1 year in the future"
	}
	return true, ""
}

// IsValidOfferCode checks Form 1095-C Line 14 codes.
func IsValidOfferCode(code string) bool {
	return validOfferCodes[strings.ToUpper(strings.TrimSpace(code))]
}

// IsValidSafeHarborCode checks Form 1095-C Line 16 codes.
func IsValidSafeHarborCode(code string) bool {
	return validSafeHarborCodes[strings.ToUpper(strings.TrimSpace(code))]
}
