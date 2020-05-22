package helpers

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"regexp"
	"strings"
	"time"
)

func GetTagFixAndFeatureCommits(fixPatterns []string, tagDateRangeTotalCommits []object.Commit, tagCommits []tagCommit) (metricTags []tagMetricData) {

	for i := 0; i < len(tagCommits); i++ {
		var featureCommits []object.Commit
		var fixCommits []object.Commit
		var tagTotalCommits []object.Commit

		if tagCommits[i].isDateRange {
			var startDate time.Time
			var endDate time.Time
			var baseDate = tagCommits[i].commit.Committer.When
			if i == 0 {
				endDate = tagCommits[i+1].commit.Committer.When
				featureCommits = FetchFeatureCommitsInDateRange(fixPatterns, tagDateRangeTotalCommits, baseDate, endDate)
				tagTotalCommits = GetTagTotalCommitsInDateRange(tagDateRangeTotalCommits, baseDate, endDate)
				fixCommits = FetchFixFirstsCommitsInDateRange(fixPatterns, tagDateRangeTotalCommits, baseDate)
			} else if i == (len(tagCommits) - 1) {
				startDate = tagCommits[i-1].commit.Committer.When
				featureCommits = FetchFeatureLastCommitsInDateRange(fixPatterns, tagDateRangeTotalCommits, baseDate)
				tagTotalCommits = GetTagTotalCommitsInDateRange(tagDateRangeTotalCommits, startDate, baseDate)
				fixCommits = FetchFixCommitsInDateRange(fixPatterns, tagDateRangeTotalCommits, startDate, baseDate)
			} else {
				startDate = tagCommits[i-1].commit.Committer.When
				endDate = tagCommits[i+1].commit.Committer.When
				featureCommits = FetchFeatureCommitsInDateRange(fixPatterns, tagDateRangeTotalCommits, baseDate, endDate)
				tagTotalCommits = GetTagTotalCommitsInDateRange(tagDateRangeTotalCommits, baseDate, endDate)
				fixCommits = FetchFixCommitsInDateRange(fixPatterns, tagDateRangeTotalCommits, startDate, baseDate)
			}

			tagMetricData := tagMetricData{
				tagDate:      tagCommits[i].tagDate,
				tag:          tagCommits[i].tag,
				fixCommits:   fixCommits,
				featCommits:  featureCommits,
				totalCommits: tagTotalCommits,
			}

			metricTags = append(metricTags, tagMetricData)
		}
	}

	return metricTags
}

func GetTagTotalCommitsInDateRange(tagDateRangeTotalCommits []object.Commit, startDate, endDate time.Time) (totalCommits []object.Commit) {
	for i := 0; i < len(tagDateRangeTotalCommits); i++ {
		if IsDateWithinRange(tagDateRangeTotalCommits[i].Committer.When, startDate, endDate) {
			totalCommits = append(totalCommits, tagDateRangeTotalCommits[i])
		}
	}

	return totalCommits
}

func FetchFeatureCommitsInDateRange(fixPatterns []string, tagDateRangeTotalCommits []object.Commit, startDate, endDate time.Time) (featureCommits []object.Commit) {
	for i := 0; i < len(tagDateRangeTotalCommits); i++ {
		if IsDateWithinRange(tagDateRangeTotalCommits[i].Committer.When, startDate, endDate) {
			if !IsFix(fixPatterns, tagDateRangeTotalCommits[i].Message) {
				featureCommits = append(featureCommits, tagDateRangeTotalCommits[i])
			}
		}
	}

	return featureCommits
}

func FetchFeatureLastCommitsInDateRange(fixPatterns []string, tagDateRangeTotalCommits []object.Commit, startDate time.Time) (featureCommits []object.Commit) {
	for i := 0; i < len(tagDateRangeTotalCommits); i++ {
		if tagDateRangeTotalCommits[i].Committer.When.Before(startDate) {
			if !IsFix(fixPatterns, tagDateRangeTotalCommits[i].Message) {
				featureCommits = append(featureCommits, tagDateRangeTotalCommits[i])
			}
		}
	}

	return featureCommits
}

func FetchFixFirstsCommitsInDateRange(fixPatterns []string, tagDateRangeTotalCommits []object.Commit, endDate time.Time) (fixCommits []object.Commit) {

	for i := 0; i < len(tagDateRangeTotalCommits); i++ {
		if tagDateRangeTotalCommits[i].Committer.When.After(endDate) {
			if IsFix(fixPatterns, tagDateRangeTotalCommits[i].Message) {
				fixCommits = append(fixCommits, tagDateRangeTotalCommits[i])
			}
		}
	}

	return fixCommits
}

func FetchFixCommitsInDateRange(fixPatterns []string, tagDateRangeTotalCommits []object.Commit, startDate, endDate time.Time) (fixCommits []object.Commit) {
	for i := 0; i < len(tagDateRangeTotalCommits); i++ {
		if IsDateWithinRange(tagDateRangeTotalCommits[i].Committer.When, startDate, endDate) {
			if IsFix(fixPatterns, tagDateRangeTotalCommits[i].Message) {
				fixCommits = append(fixCommits, tagDateRangeTotalCommits[i])
			}
		}
	}

	return fixCommits
}

func IsMergeCommit(commitMessage string) bool {
	mergePatterns := [2]string{"merge pull request", "merge branch"}
	var lowerCommitMessage = strings.ToLower(commitMessage)
	for _, mergePattern := range mergePatterns {
		matched, err := regexp.MatchString(mergePattern, lowerCommitMessage)
		if err != nil {
			println(err)
		}
		if matched {
			return true
		}
	}

	return false
}

func IsFix(fixPatterns []string, commitMessage string) bool {
	var lowerCommitMessage = strings.ToLower(commitMessage)
	for _, fixPattern := range fixPatterns {
		matched, err := regexp.MatchString("\\b"+fixPattern+"\\b", lowerCommitMessage)
		if err != nil {
			println(err)
		}
		if matched {
			return true
		}
	}

	return false
}

func GetCommitFromTagHash(repo *git.Repository, tagHash plumbing.Hash) (*object.Commit, error) {
	tag, err := repo.TagObject(tagHash)
	if err != nil {
		//fmt.Println(err)
	}

	if tag != nil {
		cm, err := tag.Commit()
		if err != nil {
			return nil, err
		}
		return cm, nil
	}

	commit, err := repo.CommitObject(tagHash)
	if err != nil {
		return nil, err
	}

	return commit, nil
}
