package template

func GetHtml() string {
	return `

<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>four-key Metrics</title>

    <link href="https://fonts.googleapis.com/css?family=Rubik:400,500&display=swap" rel="stylesheet">
    <script src="https://cdn.jsdelivr.net/npm/chart.js@2.9.3/dist/Chart.min.js"></script>

    <script type="text/javascript">
        document.addEventListener("DOMContentLoaded", function () {
            const COLORS = ["rgba(242,130,50,1)", "rgba(102,225,191,1)", "rgba(235,69,47,1)", "rgba(121,123,170,1)"];
            const LABEL_TYPES = ["weekly", "monthly", "average"];
            const CHART_NAMES = { "meanTimeChart": "meanTimeChart", "leadTimeChart": "leadTimeChart", "failPercentagesChart": "failPercentagesChart", "deploymentFrequencyChart": "deploymentFrequencyChart" };
            const ALL_MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
            const chartButtons = Array.from(document.querySelectorAll('.btn-chart'));

            const MEAN_TIME_CHART = "meanTimeChart"
            let chartData = {};
            let charts = {};
            
            chartData[CHART_NAMES["meanTimeChart"]] = {mtData};
            chartData[CHART_NAMES["leadTimeChart"]] = {ltData};
            chartData[CHART_NAMES["failPercentagesChart"]] = {fpData};
            chartData[CHART_NAMES["deploymentFrequencyChart"]] = {dfData};

            charts[CHART_NAMES["meanTimeChart"]] = "";
            charts[CHART_NAMES["leadTimeChart"]] = "";
            charts[CHART_NAMES["failPercentagesChart"]] = "";
            charts[CHART_NAMES["deploymentFrequencyChart"]] = "";

            function getWeekAndYear(d) {
                d = new Date(Date.UTC(d.getFullYear(), d.getMonth(), d.getDate()));
                const firstDayOfYear = new Date(d.getUTCFullYear(), 0, 1);
                let year = d.getUTCFullYear();
                let weekNo = Math.ceil((((d - firstDayOfYear) / 86400000) + firstDayOfYear.getDay()) / 7);

                if (weekNo > 52) {
                    weekNo = 1;
                    year = d.getUTCFullYear() + 1;
                }

                return {
                    weekNo,
                    year
                };
            }

            function prepareData(chartId) {
                const data = chartData[chartId];
                let newData = Object.assign([], data);
                let months = [];
                let years = [];
                let weeks = [];
                let days = [];
                let monthly = [];
                let weekly = [];
                let totalValue;

                newData.sort(function (a, b) {
                    return new Date(a.date) - new Date(b.date);
                });

                newData.forEach(item => {
                    const splittedDate = item.date.split("-");
                    const date = new Date(splittedDate[0] + "-" + splittedDate[1] + "-" + splittedDate[2]);
                    const month = ALL_MONTHS[date.getMonth()];
                    const day = date.getDate();
                    const year = date.getFullYear();
                    const weekAndYear = getWeekAndYear(date);

                    if (!item.month) {
                        item.month = month;
                    }
                    if (!item.day) {
                        item.day = day;
                    }
                    if (!item.year) {
                        item.year = year;
                    }
                    if (!item.week) {
                        item.week = weekAndYear.weekNo;
                    }
                    if (!days.includes(day)) {
                        days.push(day);
                    }
                    if (!months.includes(month)) {
                        months.push(month);
                    }
                    if (!years.includes(year)) {
                        years.push(year);
                    }
                    if (!weeks.includes(weekAndYear.weekNo)) {
                        weeks.push({ number: weekAndYear.weekNo, year: weekAndYear.year });
                    }
                });

                years.forEach(year => {
                    if (!monthly.includes(year)) {
                        monthly[year] = [];
                    }

                    if (!weekly.includes(year)) {
                        weekly[year] = [];
                        for (let i = 0; i < 52; i++) {
                            weekly[year].push({
                                label: i + 1,
                                totalValue: 0
                            })
                        }
                    }

                    weeks.forEach((week) => {
                        totalValue = 0;
                        let totalWeeks = weeks.filter(item =>
                            item.number == week.number && year == week.year
                        )

                        if (chartId === CHART_NAMES[MEAN_TIME_CHART]) {
                            totalValue = newData.filter(item =>
                                item.week == week.number && item.year == year
                            ).reduce((total, item) => {
                                return total + parseFloat(item.value);
                            }, totalValue);

                            weekly[year][week.number - 1].totalValue = totalValue;


                            let list = newData.filter(item =>
                                item.week == week.number && item.year == year && item.value > 0
                            )

                            weekly[year][week.number - 1].failedReleaseCount = list.length

                        } else if (chartId === CHART_NAMES["deploymentFrequencyChart"]) {
                            weekly[year][week.number - 1].totalValue = totalWeeks.length;
                        } else {
                            totalValue = newData.filter(item =>
                                item.week == week.number && item.year == year
                            ).reduce((total, item) => {
                                return total + parseFloat(item.value);
                            }, totalValue);

                            weekly[year][week.number - 1].totalValue = totalValue;
                        }

                        if (chartId === CHART_NAMES[MEAN_TIME_CHART]) {
                            console.log(weekly[year][week.number - 1])
                        }
                    });

                    ALL_MONTHS.forEach(mon => {
                        totalValue = 0;

                        totalValue = newData.filter(item => item.month == mon && item.year == year).reduce((total, item) => {
                            return total + parseFloat(item.value);
                        }, totalValue)

                        if (chartId !== MEAN_TIME_CHART) {
                            monthly[year].push({
                                label: mon,
                                year,
                                totalValue: totalValue
                            })
                        } else {
                            let failedReleases = newData.filter(item => item.month == mon && item.year == year && item.value > 0);
                            monthly[year].push({
                            label: mon,
                            year,
                            totalValue: totalValue,
                            failedReleaseCount: failedReleases.length
                        })
                        }
                    });
                });

                return {
                    data: newData,
                    monthly,
                    weekly
                }
            }

            function secondsToString(seconds) {
                const dayNumber = Math.floor((seconds % 31536000) / 86400);
                const hourNumber = Math.floor(((seconds % 31536000) % 86400) / 3600);
                const minuteNumber = Math.floor((((seconds % 31536000) % 86400) % 3600) / 60);
                const totalHourNumber = (dayNumber * 24) + hourNumber;

                return {
                    timeText: totalHourNumber + "h " + minuteNumber + "m",
                    timeTextNumeric: totalHourNumber + "." + minuteNumber,
                    totalHourNumber
                };

            }

            function getDataSets(data) {
                let index = -1;
                return Object.entries(data).map(([key, values]) => {
                    index++;
                    return {
                        label: key,
                        data: values.map(item => item.totalValue), fill: false,
                        backgroundColor: [
                            COLORS[index]
                        ],
                        borderColor: [
                            COLORS[index]
                        ],
                        borderWidth: 5
                    }
                })
            }

            function getAverageOfPreparedData(data, activeLabelType, chartId, deploymentFrequencyData) {
                for (const [key, items] of Object.entries(data[activeLabelType])) {
                    items.map(item => {
                        for (const [dfKey, dfItems] of Object.entries(deploymentFrequencyData[activeLabelType])) {

                            const releaseData = dfItems.find(dfItem => dfItem.label === item.label && dfItem.year === item.year);
                            let divider = chartId === MEAN_TIME_CHART ? item.failedReleaseCount : releaseData.totalValue
                            const averageValue = divider > 0 ? calculateAverage(item.totalValue, divider) : item.totalValue;

                            if (chartId === CHART_NAMES["meanTimeChart"] || chartId === CHART_NAMES["leadTimeChart"] && averageValue > 0) {
                                const timeLabel = secondsToString(averageValue);
                                item.totalValue = timeLabel.timeTextNumeric;
                            } else {
                                item.totalValue = averageValue;
                            }
                        }
                    })
                }

                return data;
            }

            function calculateAverage(sum, count) {
                console.log("sum -> ", sum, "count -> ", count, "avg -> ", sum / count)
                return sum / count;
            }

            function getChartData(chartId, type) {
                const activeLabelType = LABEL_TYPES.find(labelType => labelType === type);
                const deploymentFrequencyData = prepareData("deploymentFrequencyChart")
                const data = chartId === "deploymentFrequencyChart" ? deploymentFrequencyData : getAverageOfPreparedData(prepareData(chartId), activeLabelType, chartId, deploymentFrequencyData);
                return {
                    labels: Object.values(data[activeLabelType])[0].map(item => item.label),
                    datasets: getDataSets(data[activeLabelType])
                }
            }

            function createChart(chartId) {
                const ctx = document.getElementById(chartId).getContext("2d");
                charts[chartId] = new Chart(ctx, {
                    type: "line",
                    responsive: true,
                    data: getChartData(chartId, LABEL_TYPES[1]),
                    options: {
                        tooltips: {
                            mode: "point",
                            intersect: false,
                            callbacks: {
                                label: function (tooltipItem, data) {
                                    if (chartId === CHART_NAMES["deploymentFrequencyChart"]) {
                                        return data.datasets[tooltipItem.datasetIndex].label + ": " + tooltipItem.value + " Releases";
                                    } else if (chartId === CHART_NAMES["failPercentagesChart"]) {
                                        return data.datasets[tooltipItem.datasetIndex].label + ": " + "%" + parseFloat(tooltipItem.value).toFixed(2).replace(/\.?0+$/, "");
                                    } else {
                                        const label = data.datasets[tooltipItem.datasetIndex].data[tooltipItem.index].toString() || "";
                                        if (label.length > 1) {
                                            const labelArr = label.split(".");
                                            return labelArr[0] + "h " + labelArr[1].slice(0, 2) + "m";
                                        } else {
                                            return label;
                                        }
                                    }
                                }
                            }
                        }
                    }
                });
            }

            function updateChart(chartId, labelType) {
                charts[chartId].config.data = getChartData(chartId, labelType);
                charts[chartId].update();
            }

            createChart("deploymentFrequencyChart");
            createChart("leadTimeChart");
            createChart("meanTimeChart");
            createChart("failPercentagesChart");

            chartButtons.forEach(button => {
                button.addEventListener("click", function (event) {
                    const chartId = event.target.getAttribute("data-chart-id");
                    const buttonType = event.target.getAttribute("data-button-type");
                    chartButtons.forEach(btn => {
                        if (btn.classList.contains("active") && button.getAttribute("data-chart-id") === btn.getAttribute("data-chart-id")) {
                            btn.classList.remove("active");
                        }
                    });

                    if (!button.classList.contains("active")) {
                        button.classList.add("active");
                    }

                    updateChart(chartId, buttonType);
                })
            });
        })
    </script>

    <style>
        * {
            font-family: "Rubik", sans-serif;
        }

        body {
            margin: 0;
            text-rendering: optimizeLegibility !important;
            -webkit-font-smoothing: antialiased !important;
            background-color: #f4f4f4;
        }

        div {
            width: 100%;
            max-width: 100%;
        }

        .container {
            height: auto;
            width: 1000px;
            max-width: 100%;
            margin: 0 auto;
            padding: 50px;
            background-color: #ffffff;
        }

        .container .title {
            margin-top: 20px;
            margin-bottom: 10px;
            font-size: 3rem;
            display: flex;
            justify-content: center;
            font-weight: 500;
        }

        .container .subtitle {
            display: flex;
            justify-content: center;
            font-size: 1rem;
            font-weight: 400;
            margin: 0;
        }

        .container .subtitle span {
            margin: 0 10px;
        }

        .container .chart-title {
            font-size: 2rem;
            margin-left: 10px;
            font-weight: 400;
        }

        canvas {
            width: 100%;
            margin: 50px 0;
        }

        .button-wrapper {
            margin: auto;
            width: max-content;
        }

        .btn-chart {
            -moz-appearance: none;
            -webkit-appearance: none;
            -ms-progress-appearance: unset;
            border: 1px solid #767676;
            background-color: transparent;
            color: #767676;
            font-size: 12px;
            padding: 7px 14px;
            margin: 0 5px;
            outline: none;
            cursor: pointer;
            user-select: none;
        }

        .btn-chart:hover,
        .btn-chart.active {
            border: 1px solid #2185d0;
            color: #2185d0;
        }

        .btn-chart:disabled,
        .btn-chart:disabled:hover {
            border: none;
            background-color: #767676;
            color: #cdcdcd;
            cursor: not-allowed;
        }

        @media (max-width: 1080px) {
            .container {
                width: 100%;
                max-width: calc(100% - 40px);
                padding: 50px 20px;
            }
        }
    </style>
</head>

<body>
    <div class="container">
        <h4 class="title">four-key Metrics</h4>
        <h6 class="subtitle">{repositoryName} <span>|</span> {teamName} <span>|</span> {startDate} - {endDate}</h6>

        <div class="deployment-frequency-chart-wrapper">
            <h6 class="chart-title">Deployment Frequencies</h6>
            <canvas id="deploymentFrequencyChart"></canvas>
            <div class="button-wrapper">
                <button data-chart-id="deploymentFrequencyChart" data-button-type="weekly"
                    class="btn-chart">Weekly</button>
                <button data-chart-id="deploymentFrequencyChart" data-button-type="monthly"
                    class="btn-chart active">Monthly</button>
                <button data-chart-id="deploymentFrequencyChart" data-button-type="average" class="btn-chart"
                    disabled>Average</button>
            </div>
        </div>

        <div class="lead-time-chart-wrapper">
            <h6 class="chart-title">Lead Times</h6>
            <canvas id="leadTimeChart"></canvas>
            <div class="button-wrapper">
                <button data-chart-id="leadTimeChart" data-button-type="weekly" class="btn-chart">Weekly</button>
                <button data-chart-id="leadTimeChart" data-button-type="monthly"
                    class="btn-chart active">Monthly</button>
                <button data-chart-id="leadTimeChart" data-button-type="average" class="btn-chart"
                    disabled>Average</button>
            </div>
        </div>

        <div class="mean-time-chart-wrapper">
            <h6 class="chart-title">Mean Times</h6>
            <canvas id="meanTimeChart" width="400" height="400"></canvas>
            <div class="button-wrapper">
                <button data-chart-id="meanTimeChart" data-button-type="weekly" class="btn-chart">Weekly</button>
                <button data-chart-id="meanTimeChart" data-button-type="monthly"
                    class="btn-chart active">Monthly</button>
                <button data-chart-id="meanTimeChart" data-button-type="average" class="btn-chart"
                    disabled>Average</button>
            </div>
        </div>

        <div class="fail-percengates-chart-wrapper">
            <h6 class="chart-title">Fail Percentages</h6>
            <canvas id="failPercentagesChart"></canvas>
            <div class="button-wrapper">
                <button data-chart-id="failPercentagesChart" data-button-type="weekly" class="btn-chart">Weekly</button>
                <button data-chart-id="failPercentagesChart" data-button-type="monthly"
                    class="btn-chart active">Monthly</button>
                <button data-chart-id="failPercentagesChart" data-button-type="average" class="btn-chart"
                    disabled>Average</button>
            </div>
        </div>
    </div>
</body>
</html>

	`
}
