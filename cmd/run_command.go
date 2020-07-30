package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	Command "four-key/command"
	. "four-key/helpers"
	. "four-key/models"
	"four-key/settings"
	"four-key/template"
	"github.com/spf13/cobra"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var runCommand = &cobra.Command{
	Use:   "run",
	Short: "run repository",
	Long:  "run repository",
	Run:   onRun,
}

type ChartItem struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
	Name  string  `json:"name"`
}

type ChartItems struct {
	meanTimes             []ChartItem
	leadTimes             []ChartItem
	failPercentages       []ChartItem
	deploymentFrequencies []ChartItem
	teamName              string
}

type SeparatedTeamItems struct {
	teams map[string][]ChartItems
}

var allChartItems []ChartItems
var separatedTeamItems SeparatedTeamItems
var commander Command.ICommand
var s *settings.Settings

const deploymentInitialValue = 1

func init() {
	rootCmd.AddCommand(runCommand)
	runCommand.Flags().StringP("startDate", "s", "", "Set a start date of range")
	runCommand.Flags().StringP("endDate", "e", "", "Set a end date of range")
	runCommand.Flags().StringP("repository", "r", "", "Set a name of the specific repository")
	commander = Command.ACommander()
}

func onRun(cmd *cobra.Command, args []string) {
	repositoryName, err := cmd.Flags().GetString("repository")
	startDateInput, err := cmd.Flags().GetString("startDate")
	startDate, err := time.Parse(settings.DefaultDateFormat, startDateInput)

	if err != nil {
		fmt.Println(commander.Warn("invalid start date, start date will be like -s YYYY-MM-DD"))
		return
	}

	endDateInput, err := cmd.Flags().GetString("endDate")
	endDate, err := time.Parse(settings.DefaultDateFormat, endDateInput)

	if err != nil {
		fmt.Println(commander.Warn("invalid end date, end date will be like -s YYYY-MM-DD"))
		return
	}

	err = settings.Initialize(commander)
	if err != nil {
		fmt.Println(commander.Fatal(err.Error()))
		return
	}

	s, err = settings.Get()
	if err != nil || s == nil {
		if err != nil {
			fmt.Println(commander.Fatal(err.Error()))
			return
		}

		fmt.Println(commander.Fatal("configurations didn't loaded"))
		return
	}

	var repositories []RepositoryWrapper
	if repositoryName != "" {
		repository, err := GetRepositoryByName(s, repositoryName)

		if err != nil {
			fmt.Println(commander.Fatal(err.Error()))
			return
		}

		repositories = append(repositories, repository)

	} else {
		repositories, err = GetRepositories(s)
	}

	if err != nil {
		fmt.Println(commander.Fatal("git clone returned an error - error -> ", err.Error()))
		return
	}

	t := strings.Trim(startDate.Format(settings.DefaultDateFormat)+"--"+endDate.Format(settings.DefaultDateFormat), "")

	var metricResultDtoList []FourKeyMetricResultDto

	if len(repositories) < 1 {
		fmt.Println(commander.Fatal("repository not found. please use -> $four-key add command or modify your configuration file \n"))
		fmt.Println(commander.Good(fmt.Sprintf("configuration file path: %s ", commander.GetFourKeyPath())))
		_ = commander.Open(commander.GetFourKeyPath())
		return
	}

	for _, repo := range repositories {
		metricsRequest := MetricsRequest{
			StartDate:         startDate,
			EndDate:           endDate,
			ReleaseTagPattern: repo.Configurations.ReleaseTagPattern,
			FixPatterns:       repo.Configurations.FixCommitPatterns,
		}

		metricsDto, err := CalculateMetrics(repo.Repository, metricsRequest)

		if err != nil {
			fmt.Println(commander.Fatal(err.Error(), " Project -> ", repo.Configurations.Name()))
			continue
		}

		metricsDto.RepoName = repo.Configurations.Name()
		metricsDto.TeamName = repo.Configurations.TeamName
		metricResultDtoList = append(metricResultDtoList, metricsDto)
	}

	generateMetricFiles(metricResultDtoList, t)
}

func generateMetricFiles(metricResultDtoList []FourKeyMetricResultDto, reportTimeAsString string) {
	outputSource := path.Join(s.Output, settings.DefaultGeneratedFileOutputDirName)
	separatedTeamItems.teams = map[string][]ChartItems{}
	err := CheckDirectory(outputSource)

	if err != nil {
		_ = CreateDirectory(s.Output, settings.DefaultGeneratedFileOutputDirName)
	}

	err = CheckDirectory(outputSource, settings.AllTeamsDefaultDirName)

	if err != nil {
		_ = CreateDirectory(outputSource, settings.AllTeamsDefaultDirName)
	}

	err = CheckDirectory(outputSource, settings.TeamBasedDefaultDirName)

	if err != nil {
		_ = CreateDirectory(outputSource, settings.TeamBasedDefaultDirName)
	}

	for i, metric := range metricResultDtoList {
		dirName := metric.RepoName + "_" + reportTimeAsString

		if metric.TeamName == "" {
			metric.TeamName = settings.DefaultTeamName
		}

		err = CheckDirectory(outputSource, metric.TeamName)

		if err != nil {
			_ = CreateDirectory(outputSource, metric.TeamName)
		}

		dirName, err := generateDirectory(path.Join(outputSource, metric.TeamName), dirName)

		if err != nil {
			fmt.Println(commander.Fatal(err.Error()))
			return
		}

		metrics := createChartItems(metric.MetricTags)
		metrics.teamName = metric.TeamName
		allChartItems = append(allChartItems, metrics)
		separatedTeamItems.teams[metrics.teamName] = append(separatedTeamItems.teams[metrics.teamName], metrics)
		err = generateOutput(dirName, metrics, metric, outputSource, false)

		if i == len(metricResultDtoList)-1 {
			var teamBasedMetrics []ChartItems
			var allTeamsMetrics ChartItems
			allTeamsMetrics.teamName = "AllTeamsResult"

			for name, chartItems := range separatedTeamItems.teams {
				teamChart := ChartItems{
					teamName: name,
				}
				mergeChartItems(chartItems, &teamChart)
				teamBasedMetrics = append(teamBasedMetrics, teamChart)
			}

			mergeChartItems(allChartItems, &allTeamsMetrics)

			for _, team := range teamBasedMetrics {
				metricDto := FourKeyMetricResultDto{
					RepoName:                 settings.TeamBasedDefaultDirName,
					TeamName:                 settings.TeamBasedDefaultDirName,
					DateRangeStart:           metric.DateRangeStart,
					DateRangeEnd:             metric.DateRangeEnd,
					CreationDate:             metric.CreationDate,
					DeploymentFrequencyCount: 0,
				}

				dirName, err := generateDirectory(path.Join(outputSource, settings.TeamBasedDefaultDirName), team.teamName+"_"+reportTimeAsString)

				if err != nil {
					fmt.Println(commander.Fatal(err.Error()))
					return
				}

				err = generateOutput(dirName, team, metricDto, outputSource, false)
			}

			metricDto := FourKeyMetricResultDto{
				RepoName:                 settings.AllTeamsDefaultDirName,
				TeamName:                 settings.AllTeamsDefaultDirName,
				DateRangeStart:           metric.DateRangeStart,
				DateRangeEnd:             metric.DateRangeEnd,
				CreationDate:             metric.CreationDate,
				DeploymentFrequencyCount: 0,
			}

			dirName, err := generateDirectory(path.Join(outputSource, settings.AllTeamsDefaultDirName), allTeamsMetrics.teamName+"_"+reportTimeAsString)

			if err != nil {
				fmt.Println(commander.Fatal(err.Error()))
				return
			}

			err = generateOutput(dirName, allTeamsMetrics, metricDto, outputSource, true)

			if err != nil {
				fmt.Println(commander.Warn("an error occurred while opening results folder", " you can see in -> ", path.Join(commander.GetFourKeyPath(), outputSource)))
			}
		} else {
			if err != nil {
				fmt.Println(commander.Warn("an error occurred while opening results folder", " you can see in -> ", path.Join(commander.GetFourKeyPath(), outputSource)))
			}
		}
	}
}

func createChartItems(metrics []TagMetricDto) ChartItems {
	var chartItems ChartItems
	for _, t := range metrics {
		chartItems.meanTimes = append(chartItems.meanTimes, ChartItem{
			Date:  t.TagDate.Format(settings.DefaultDateFormat),
			Value: t.MeanTimeRestoreAverage,
		})
		chartItems.leadTimes = append(chartItems.leadTimes, ChartItem{
			Date:  t.TagDate.Format(settings.DefaultDateFormat),
			Value: t.LeadTime,
		})
		chartItems.failPercentages = append(chartItems.failPercentages, ChartItem{
			Date:  t.TagDate.Format(settings.DefaultDateFormat),
			Value: t.ChangeFailPercentage,
		})
		chartItems.deploymentFrequencies = append(chartItems.deploymentFrequencies, ChartItem{
			Date:  t.TagDate.Format(settings.DefaultDateFormat),
			Value: deploymentInitialValue,
			Name:  t.TagName,
		})
	}

	return chartItems
}

func generateOutput(dir string, items ChartItems, results FourKeyMetricResultDto, outputSource string, open bool) error {
	var h *os.File
	h, err := os.Create(path.Join(outputSource, results.TeamName, dir, "index.html"))

	if err != nil {
		return err
	}

	html, err := createHtml(results, items)

	if err != nil {
		return err
	}

	_, err = h.WriteString(html)

	if err != nil {
		err = h.Close()

		if err != nil {
			return err
		}

		return err
	}

	if open {
		err := commander.Open(outputSource)

		if err != nil {
			fmt.Println(commander.Warn("an error occurred while opening results folder", " you can see in -> ", path.Join(commander.GetFourKeyPath(), outputSource)))
		}
	}

	fmt.Println(commander.Good("metrics file generated", " for -> ", results.TeamName, ":", results.RepoName, "in -> ", path.Join(outputSource, results.TeamName)))

	return nil
}

func generateDirectory(sourceDir, dir string) (string, error) {
	err := CreateDirectory(sourceDir, dir)

	if err != nil {
		counter := 0
		for {
			counter++

			if counter > 1000 {
				return "", errors.New("an infinite loop occurred while trying to create metric file")
			}

			d := dir + "_" + strconv.Itoa(counter)

			err = CreateDirectory(sourceDir, d)

			if err != nil {
				continue
			}

			return d, nil
		}
	}

	return dir, nil
}

func createHtml(dto FourKeyMetricResultDto, items ChartItems) (string, error) {

	mtJson, err := json.Marshal(items.meanTimes)
	ltJson, err := json.Marshal(items.leadTimes)
	fpJson, err := json.Marshal(items.failPercentages)
	dfJson, err := json.Marshal(items.deploymentFrequencies)

	if err != nil {
		fmt.Println(commander.Fatal("an error occurred while serializing"))
		return "", err
	}

	htmlTemplate := template.GetHtml()
	htmlTemplate = strings.Replace(htmlTemplate, "{repositoryName}", dto.RepoName, 1)
	htmlTemplate = strings.Replace(htmlTemplate, "{teamName}", dto.TeamName, 1)
	htmlTemplate = strings.Replace(htmlTemplate, "{startDate}", dto.DateRangeStart.Format(settings.DefaultDateFormat), 1)
	htmlTemplate = strings.Replace(htmlTemplate, "{endDate}", dto.DateRangeEnd.Format(settings.DefaultDateFormat), 1)
	htmlTemplate = strings.Replace(htmlTemplate, "{mtData}", string(mtJson), 1)
	htmlTemplate = strings.Replace(htmlTemplate, "{ltData}", string(ltJson), 1)
	htmlTemplate = strings.Replace(htmlTemplate, "{fpData}", string(fpJson), 1)
	htmlTemplate = strings.Replace(htmlTemplate, "{dfData}", string(dfJson), 1)

	return htmlTemplate, err
}

func mergeChartItems(source []ChartItems, target *ChartItems) {
	for _, m := range source {
		for _, item := range m.deploymentFrequencies {
			target.deploymentFrequencies = append(target.deploymentFrequencies, item)
		}
		for _, item := range m.failPercentages {
			target.failPercentages = append(target.failPercentages, item)
		}
		for _, item := range m.meanTimes {
			target.meanTimes = append(target.meanTimes, item)
		}
		for _, item := range m.leadTimes {
			target.leadTimes = append(target.leadTimes, item)
		}
	}
}
