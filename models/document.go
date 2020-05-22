package models

type Document struct {
	Output       string                `json:"output"`
	Repositories []*DocumentRepository `json:"repositories"`
}

type DocumentRepository struct {
	TeamName          string   `json:"teamName"`
	CloneAddress      string   `json:"cloneAddress"`
	ReleaseTagPattern string   `json:"releaseTagPattern"`
	FixCommitPatterns []string `json:"fixCommitPatterns"`
}
