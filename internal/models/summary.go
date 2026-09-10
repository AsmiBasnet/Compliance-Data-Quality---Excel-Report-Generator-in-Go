package models

// SummaryStats aggregates system-wide data quality and compliance metrics.
type SummaryStats struct {
	TotalRecords        int                        `json:"total_records"`
	ValidRecords        int                        `json:"valid_records"`
	ErrorRecords        int                        `json:"error_records"`
	MissingFieldCount   int                        `json:"missing_field_count"`
	DuplicateCount      int                        `json:"duplicate_count"`
	ErrorRate           float64                    `json:"error_rate"`
	TotalMonthlyHours   float64                    `json:"total_monthly_hours"`
	AverageMonthlyHours float64                    `json:"average_monthly_hours"`
	PotentialFTCount    int                        `json:"potential_ft_count"` // Employees >= 130 hours/mo
	ErrorsByType        map[string]int             `json:"errors_by_type"`
	ErrorsBySeverity    map[string]int             `json:"errors_by_severity"`
	DepartmentStats     map[string]*DepartmentStat `json:"department_stats"`
}

// DepartmentStat contains category-level breakdown for departments.
type DepartmentStat struct {
	Department   string  `json:"department"`
	TotalCount   int     `json:"total_count"`
	ValidCount   int     `json:"valid_count"`
	ErrorCount   int     `json:"error_count"`
	ErrorRate    float64 `json:"error_rate"`
	AverageHours float64 `json:"average_hours"`
}
