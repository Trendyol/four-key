package models

import "time"

type MetricsRequest struct {
	StartDate         time.Time
	EndDate           time.Time
	ReleaseTagPattern string
	FixPatterns       []string
}
