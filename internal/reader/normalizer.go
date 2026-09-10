package reader

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	nonDigitsRegex     = regexp.MustCompile(`[^\d]`)
	cleanCharsRegex    = regexp.MustCompile(`[^a-zA-Z0-9]`)
	multipleSpaceRegex = regexp.MustCompile(`\s+`)
)

// NormalizeHeader maps messy/varied CSV header column names to standard keys.
func NormalizeHeader(header string) string {
	h := strings.ToLower(strings.TrimSpace(header))
	h = cleanCharsRegex.ReplaceAllString(h, "_")
	h = multipleSpaceRegex.ReplaceAllString(h, "_")
	h = strings.Trim(h, "_")

	switch {
	case strings.Contains(h, "emp") && strings.Contains(h, "id"), h == "id", h == "employee_number", h == "eid":
		return "employee_id"
	case strings.Contains(h, "first") || h == "fname" || h == "given_name":
		return "first_name"
	case strings.Contains(h, "last") || h == "lname" || h == "surname" || h == "family_name":
		return "last_name"
	case strings.Contains(h, "ssn") || strings.Contains(h, "social_security") || strings.Contains(h, "tin"):
		return "ssn"
	case strings.Contains(h, "dept") || strings.Contains(h, "department") || strings.Contains(h, "division") || strings.Contains(h, "cost_center"):
		return "department"
	case strings.Contains(h, "status") || strings.Contains(h, "employment_type") || strings.Contains(h, "ft_pt"):
		return "employment_status"
	case strings.Contains(h, "hire") || strings.Contains(h, "start_date") || h == "doh":
		return "hire_date"
	case strings.Contains(h, "hour") || strings.Contains(h, "hrs"):
		return "monthly_hours"
	case strings.Contains(h, "pay") || strings.Contains(h, "wage") || strings.Contains(h, "salary") || strings.Contains(h, "gross"):
		return "monthly_gross_pay"
	case strings.Contains(h, "line_14") || strings.Contains(h, "offer") || strings.Contains(h, "coverage_code") || strings.Contains(h, "mec_offer"):
		return "coverage_offer_code"
	case strings.Contains(h, "line_16") || strings.Contains(h, "safe_harbor") || strings.Contains(h, "harbor_code") || strings.Contains(h, "4980h"):
		return "safe_harbor_code"
	case strings.Contains(h, "contrib") || strings.Contains(h, "premium") || strings.Contains(h, "cost") || strings.Contains(h, "share"):
		return "employee_contribution"
	default:
		return h
	}
}

// NormalizeSSN strips formatting and formats valid 9-digit SSNs into XXX-XX-XXXX.
func NormalizeSSN(raw string) (formatted string, digits string) {
	d := nonDigitsRegex.ReplaceAllString(raw, "")
	if len(d) == 9 {
		return d[0:3] + "-" + d[3:5] + "-" + d[5:9], d
	}
	return strings.TrimSpace(raw), d
}

// NormalizeDate parses multiple common date formats and returns ISO 8601 YYYY-MM-DD.
func NormalizeDate(raw string) (string, time.Time, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", time.Time{}, nil
	}

	formats := []string{
		"2006-01-02",
		"01/02/2006",
		"1/2/2006",
		"2006/01/02",
		"02-Jan-2006",
		"January 2, 2006",
		"01-02-2006",
		"20060102",
	}

	for _, layout := range formats {
		if t, err := time.Parse(layout, trimmed); err == nil {
			return t.Format("2006-01-02"), t, nil
		}
	}

	return trimmed, time.Time{}, strconv.ErrSyntax
}

// NormalizeCurrency converts raw currency strings (e.g., "$1,250.00") to float64.
func NormalizeCurrency(raw string) float64 {
	cleaned := strings.ReplaceAll(raw, "$", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.TrimSpace(cleaned)
	val, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0.0
	}
	return val
}

// NormalizeFloat parses numeric values with comma handling.
func NormalizeFloat(raw string) float64 {
	cleaned := strings.ReplaceAll(raw, ",", "")
	cleaned = strings.TrimSpace(cleaned)
	val, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0.0
	}
	return val
}
