package models

import "time"

type TagMetricDto struct {
	TagName                string    `json:"Name"`
	TagDate                time.Time `json:"Date"`
	LeadTime               float64   `json:"LeadTime"`
	MeanTimeRestoreAverage float64   `json:"MeanTimeRestoreAverage"`
	ChangeFailPercentage   float64   `json:"ChangeFailPercentage"`
}
