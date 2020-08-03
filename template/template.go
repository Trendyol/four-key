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
            const LABEL_TYPES = ["weekly", "monthly", "average"]
            const ALL_MONTHS_LONG = ["January", "February", "March", "April", "May", "June",
                "July", "August", "September", "October", "November", "December"];
            const ALL_MONTHS_SHORT = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"]
            const chartButtons = Array.from(document.querySelectorAll('.btn-chart'))

            let chartData = {
                meanTimeChart: {mtData},
                leadTimeChart: {ltData},
                failPercengatesChart: {fpData},
                deploymentFrequencyChart: {dfData}
            }
            let charts = {
                meanTimeChart: "",
                leadTimeChart: "",
                failPercengatesChart: "",
                deploymentFrequencyChart: ""
            }

            function prepareData(data) {
                let newData = Object.assign([], data)
                let months = []
                let years = []
                let days = []
                let monthly = []

                newData.sort(function (a, b) {
                    return new Date(a.date) - new Date(b.date);
                });

                newData.forEach(item => {
                    const splittedDate = item.date.split("-")
                    const date = new Date(splittedDate[0] + "-" + splittedDate[1] + "-" + splittedDate[2])
                    const month = ALL_MONTHS_SHORT[date.getMonth()]
                    const day = date.getDate()
                    const year = date.getFullYear()
                    if (!item.month) {
                        item.month = month
                    }
                    if (!item.day) {
                        item.day = day
                    }
                    if (!item.year) {
                        item.year = year
                    }
                    if (!days.includes(day)) {
                        days.push(day)
                    }
                    if (!months.includes(month)) {
                        months.push(month)
                    }
                    if (!years.includes(year)) {
                        years.push(year)
                    }
                });

                years.forEach(year => {
                    if (!monthly.includes(year)) {
                        monthly[year] = []
                    }

                    let totalValue = 0
                    totalValue = newData.filter(item => item.year == year).reduce(function getSum(total, item) {
                        return total + parseInt(item.value);
                    }, totalValue)

                    ALL_MONTHS_SHORT.forEach(mon => {
                        totalValue = 0
                        totalValue = newData.filter(item => item.month == mon && item.year == year).reduce(function getSum(total, item) {
                            return total + parseInt(item.value);
                        }, totalValue)

                        monthly[year].push({
                            label: mon,
                            totalValue
                        })

                    });

                });

                return {
                    data: newData,
                    monthly
                }
            }

            function randomNum() {
                return Math.floor(Math.random() * 256);
            }

            function randomRGB() {
                var red = randomNum();
                var green = randomNum();
                var blue = randomNum();
                return red + ", " + green + ", 235";
            }

            function getDataSets(data) {
                const a = Object.entries(data).map(([key, values]) => {
                    return {
                        label: key,
                        data: values.map(item => item.totalValue), fill: false,
                        backgroundColor: [
                            "rgba(" + randomRGB() + ", 0.2)"
                        ],
                        borderColor: [
                            "rgba(" + randomRGB() + ", 1)"
                        ],
                        borderWidth: 5
                    }
                })

                return a
            }

            function getChartData(chartId, type) {
                const data = prepareData(chartData[chartId])
                const activeLabelType = LABEL_TYPES.find(labelType => labelType == type)

                return {
                    labels: activeLabelType == "monthly" ? ALL_MONTHS_SHORT : Object.keys(data[activeLabelType]),
                    datasets: getDataSets(data[activeLabelType])
                }
            }

            function createChart(chartId) {
                const ctx = document.getElementById(chartId).getContext('2d');
                charts[chartId] = new Chart(ctx, {
                    type: 'line',
                    data: getChartData(chartId, LABEL_TYPES[1]),
                    options: {
                        scales: {
                            yAxes: [{
                                ticks: {
                                    beginAtZero: true
                                }
                            }]
                        }
                    }
                });
            }

            function updateChart(chartId, labelType) {
                charts[chartId].config.data = getChartData(chartId, labelType)
                charts[chartId].update()
            }

            createChart("deploymentFrequencyChart")

            chartButtons.forEach(button => {
                button.addEventListener("click", function (event) {
                    const chartId = event.target.getAttribute("data-chart-id")
                    const buttonType = event.target.getAttribute('data-button-type')
                    updateChart(chartId, buttonType)
                })
            });
        })
        </script>

        <style>
        body {
            font-family: 'Rubik', sans-serif;
        }

        div {
            width: 100%;
            max-width: 100%;
        }

        .container {
            height: 100%;
            width: 100%;
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
            margin-block-end: 0.5em;
            font-size: 2rem;
            margin-left: 10px;
            font-weight: 400;
        }

        canvas {
            margin: 50px 0;
        }

        .button-wrapper {
            margin: auto;
            width: max-content;
        }

        .btn-chart {
            appearance: none;
            border: 1px solid #767676;
            background-color: transparent;
            color: #767676;
            padding: 7px 14px;
            margin: 0 5px;
            outline: none;
            cursor: pointer;
        }

        .btn-chart:hover {
            border: 1px solid #2185d0;
            color: #2185d0;
        }
    </style>
      </head>
      <body>
      <div class="container">
      <h4 class="title">four-key Metrics</h4>
      <h6 class="subtitle">allTeams <span>|</span> allTeams <span>|</span> 2019-01-01 - 2021-01-01</h6>

      <div class="deployment-frequency-chart-wrapper">
          <h6 class="chart-title">Deployment Frequencies</h6>
          <canvas id="deploymentFrequencyChart"></canvas>
          <div class="button-wrapper">
              <button data-chart-id="deploymentFrequencyChart" data-button-type="weekly" class="btn-chart">Weekly</button>
              <button data-chart-id="deploymentFrequencyChart" data-button-type="monthly" class="btn-chart">Monthly</button>
              <button data-chart-id="deploymentFrequencyChart" data-button-type="average" class="btn-chart">Average</button>
          </div>
      </div>

      <h6 class="chart-title">Lead Times</h6>
      <canvas id="leadTimeChart"></canvas>

      <h6 class="chart-title">Mean Times</h6>
      <canvas id="meanTimeChart" width="400" height="400"></canvas>
      <!--   <div id="meanTimeDiv"></div> -->

      <h6 class="chart-title">Fail Percentages</h6>
      <canvas id="failPercengatesChart"></canvas>
  </div>
</body>
</html>

	`
}
