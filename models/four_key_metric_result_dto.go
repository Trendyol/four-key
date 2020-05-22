package models

import "time"

type FourKeyMetricResultDto struct {
	RepoName                 string
	TeamName                 string
	DateRangeStart           time.Time      `json:"DateRangeStart"`
	DateRangeEnd             time.Time      `json:"DateRangeEnd"`
	CreationDate             time.Time      `json:"CreationDate"`
	MetricTags               []TagMetricDto `json:"MetricTags"`
	DeploymentFrequencyCount int64          `json:"DeploymentFrequencyCount"`
}
