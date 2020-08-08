package helpers

import (
	_ "container/list"
	"errors"
	. "four-key/models"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"regexp"
	_ "sort"
	"strings"
	"time"
)

type tagCommit struct {
	commit      object.Commit
	tag         *plumbing.Reference
	tagDate     time.Time
	isDateRange bool
	tagType     string
	tagComment  string
}

type FourKeyMetricsDto struct {
	TagMetrics []tagMetricData
	Df         DeploymentFrequencyDto
	RepoName   string
	TeamName   string
}

type tagMetricData struct {
	fixCommits                       []object.Commit
	featCommits                      []object.Commit
	totalCommits                     []object.Commit
	tag                              *plumbing.Reference
	tagDate                          time.Time
	tagMeanTimeRestoreAverageSeconds float64
	tagLeadTimeSeconds               float64
	tagChangeFailPercentage          float64
	deploymentFrequency              float64
}

type DeploymentTag struct {
	When  time.Time `json:"date"`
	Name  string
	Count int `json:"value"`
}

type DeploymentFrequencyDto struct {
	Tags []DeploymentTag `json:"Deployments"`
}

func CalculateMetrics(repo *git.Repository, request MetricsRequest) (FourKeyMetricResultDto, error) {

	var keyMetrics FourKeyMetricResultDto

	RepoCheck(repo)

	tagCommits, err := getTagCommitBetweenDates(repo, request)

	if len(tagCommits) > 1 {

		var tagDateRangeTotalCommits = GetTagDateRangeTotalCommits(repo, tagCommits)

		//fix and fea commits found
		var tagFixAndFeatureCommits = GetTagFixAndFeatureCommits(request.FixPatterns, tagDateRangeTotalCommits, tagCommits)

		//added MeanTimeToRestore
		tagFixAndFeatureCommits = GetMeanTimeToRestore(tagFixAndFeatureCommits)

		//added ChangeFailPercentage
		tagFixAndFeatureCommits = GetChangeFailPercentage(tagFixAndFeatureCommits)

		//added LeadTime
		tagFixAndFeatureCommits = GetLeadTime(tagFixAndFeatureCommits)

		if err != nil {
			return keyMetrics, err
		}

		var tagMetricDtoList []TagMetricDto
		for _, tagMetricDateRange := range tagFixAndFeatureCommits {
			tagNameParse := strings.Split(string(tagMetricDateRange.tag.Name()), "/")
			tagMetricDto := TagMetricDto{
				TagName:                tagNameParse[len(tagNameParse)-1],
				TagDate:                tagMetricDateRange.tagDate,
				MeanTimeRestoreAverage: tagMetricDateRange.tagMeanTimeRestoreAverageSeconds,
				LeadTime:               tagMetricDateRange.tagLeadTimeSeconds,
				ChangeFailPercentage:   tagMetricDateRange.tagChangeFailPercentage,
			}
			tagMetricDtoList = append(tagMetricDtoList, tagMetricDto)
		}

		keyMetrics.CreationDate = time.Now()
		keyMetrics.MetricTags = tagMetricDtoList
		keyMetrics.DateRangeEnd = request.EndDate
		keyMetrics.DateRangeStart = request.StartDate
		keyMetrics.DeploymentFrequencyCount = int64(len(tagMetricDtoList))

		return keyMetrics, nil
	}

	return keyMetrics, errors.New("metrics could not be calculated because there is no release")
}

func GetTagDateRangeTotalCommits(repo *git.Repository, tagCommits []tagCommit) []object.Commit {

	var descendingSortCommits = GetDescendingCommits(repo)

	var lastTagCommit = tagCommits[len(tagCommits)-1]
	var firstTagCommit = tagCommits[0]

	var isFistTagDateRangeCommit = firstTagCommit.isDateRange
	var isLastTagDateRangeCommit = lastTagCommit.isDateRange
	var tagDateRangeTotalCommits []object.Commit

	var isDateRangeCommitFounds = false
	for _, sortCommit := range descendingSortCommits {
		var isDateRangeCommit = IsDateWithinRange(sortCommit.Commit.Committer.When, firstTagCommit.tagDate, lastTagCommit.tagDate)

		var isCommitAfterTagDate = sortCommit.Commit.Committer.When.After(firstTagCommit.tagDate)
		if isCommitAfterTagDate && isFistTagDateRangeCommit {
			isDateRangeCommitFounds = true
		}

		if isDateRangeCommit {
			isDateRangeCommitFounds = true
		}
		var isCommitBeforeTagDate = sortCommit.Commit.Committer.When.Before(lastTagCommit.tagDate)
		if isCommitBeforeTagDate && isLastTagDateRangeCommit {
			isDateRangeCommitFounds = true
		}

		if isDateRangeCommitFounds {
			if !IsMergeCommit(sortCommit.Commit.Message) {
				tagDateRangeTotalCommits = append(tagDateRangeTotalCommits, sortCommit.Commit)
			}
		}

		if !isDateRangeCommit {
			isDateRangeCommitFounds = false
		}
	}

	return tagDateRangeTotalCommits
}

func IsReleaseTag(tagName, releaseTagPattern string) bool {
	var lowerTagName = strings.ToLower(tagName)
	var lowerReleaseTagPattern = strings.ToLower(releaseTagPattern)
	matched, err := regexp.MatchString(lowerReleaseTagPattern, lowerTagName)
	if err != nil {
		println(err)
	}

	return matched
}

func getTagCommitBetweenDates(r *git.Repository, request MetricsRequest) ([]tagCommit, error) {
	var commitTags []tagCommit

	var sortedTagList = GetAscendingOrderByTagDate(r)

	var prevTag *tagCommit
	var lastTag *tagCommit
	var firstTag *tagCommit
	lastTagFound := false
	firstTagFound := false
	for _, t := range sortedTagList {
		if !IsReleaseTag(t.tag.Name().Short(), request.ReleaseTagPattern) {
			continue
		}

		tagCmt, err := GetCommitFromTagHash(r, t.tag.Hash())
		if err != nil {
			return commitTags, err
		}

		cTag := tagCommit{
			commit:      *tagCmt,
			tag:         t.tag,
			isDateRange: true,
			tagDate:     tagCmt.Committer.When,
			tagType:     "TagIsDateRange",
			tagComment:  "Tag in the date range",
		}

		tagIsDateRange := inTimeSpan(request.StartDate, request.EndDate, tagCmt.Committer.When)

		if tagIsDateRange {
			if !lastTagFound {
				lastTag = prevTag
				lastTagFound = true
			}

			commitTags = append(commitTags, cTag)
		}

		if lastTagFound && !firstTagFound && !tagIsDateRange {
			firstTag = &cTag
			firstTagFound = true
		}

		prevTag = &cTag
	}

	var sortedCommitTags []tagCommit

	if firstTag != nil {
		firstTag.isDateRange = false
		firstTag.tagType = "firstTag"
		firstTag.tagComment = "First tag not in the date range"
		sortedCommitTags = append(sortedCommitTags, *firstTag)
	}

	for i := len(commitTags) - 1; i >= 0; i-- {
		sortedCommitTags = append(sortedCommitTags, commitTags[i])
	}

	if lastTag != nil {
		lastTag.isDateRange = false
		lastTag.tagType = "lastTag"
		lastTag.tagComment = "Tag before last tag in date range"
		sortedCommitTags = append(sortedCommitTags, *lastTag)
	}

	return sortedCommitTags, nil
}
