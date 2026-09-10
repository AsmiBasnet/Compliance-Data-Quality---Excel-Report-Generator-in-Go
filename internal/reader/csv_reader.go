package reader

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"compliance-report-generator/internal/models"
)

// CSVReader handles ingestion of one or more compliance CSV files.
type CSVReader struct{}

// NewCSVReader creates a new CSV reader instance.
func NewCSVReader() *CSVReader {
	return &CSVReader{}
}

// ReadFiles reads multiple CSV paths or glob patterns and returns raw standardized records.
func (r *CSVReader) ReadFiles(filePaths ...string) ([]models.EmployeeRecord, error) {
	var allRecords []models.EmployeeRecord
	recordCounter := 1

	for _, pattern := range filePaths {
		matches, err := filepath.Glob(pattern)
		if err != nil || len(matches) == 0 {
			// Try direct path if glob has no matches
			matches = []string{pattern}
		}

		for _, filePath := range matches {
			records, err := r.ReadFile(filePath, recordCounter)
			if err != nil {
				return nil, fmt.Errorf("failed reading file %s: %w", filePath, err)
			}
			allRecords = append(allRecords, records...)
			recordCounter += len(records)
		}
	}

	return allRecords, nil
}

// ReadFile parses a single CSV file and normalizes its fields into EmployeeRecord structs.
func (r *CSVReader) ReadFile(filePath string, startID int) ([]models.EmployeeRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1 // Flexible field count

	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading header row: %w", err)
	}

	headerMap := make(map[string]int)
	for i, h := range headers {
		normalized := NormalizeHeader(h)
		headerMap[normalized] = i
	}

	var records []models.EmployeeRecord
	currID := startID

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// Skip completely corrupted lines or parse what's available
			continue
		}

		// Check if row is entirely empty
		isEmpty := true
		for _, col := range row {
			if strings.TrimSpace(col) != "" {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			continue
		}

		getVal := func(key string) string {
			idx, ok := headerMap[key]
			if !ok || idx >= len(row) {
				return ""
			}
			return strings.TrimSpace(row[idx])
		}

		rawSSN := getVal("ssn")
		formattedSSN, _ := NormalizeSSN(rawSSN)
		rawDate := getVal("hire_date")
		normDate, parsedDate, _ := NormalizeDate(rawDate)

		rec := models.EmployeeRecord{
			RecordID:             currID,
			EmployeeID:           getVal("employee_id"),
			FirstName:            getVal("first_name"),
			LastName:             getVal("last_name"),
			SSN:                  formattedSSN,
			Department:           getVal("department"),
			EmploymentStatus:     getVal("employment_status"),
			HireDate:             normDate,
			ParsedHireDate:       parsedDate,
			MonthlyHours:         NormalizeFloat(getVal("monthly_hours")),
			MonthlyGrossPay:      NormalizeCurrency(getVal("monthly_gross_pay")),
			CoverageOfferCode:    strings.ToUpper(getVal("coverage_offer_code")),
			SafeHarborCode:       strings.ToUpper(getVal("safe_harbor_code")),
			EmployeeContribution: NormalizeCurrency(getVal("employee_contribution")),
			IsValid:              true, // Will be validated by validator
		}

		records = append(records, rec)
		currID++
	}

	return records, nil
}
