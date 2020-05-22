package template

func GetHtml() string {
	return `

<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>four-key Metrics</title>

    <link href="https://fonts.googleapis.com/css?family=Rubik:400,500&display=swap" rel="stylesheet">
    <script src="https://www.amcharts.com/lib/4/core.js"></script>
    <script src="https://www.amcharts.com/lib/4/charts.js"></script>
    <script src="https://www.amcharts.com/lib/4/themes/animated.js"></script>

    <script type="text/javascript">
        am4core.ready(function () {
            var mtData = {mtData};
            initMeanTime(mtData, groupBy(mtData, false));

            var ltData = {ltData};
            initLeadTime(ltData, groupBy(ltData, false));

            var fpData = {fpData};
            initFailPercentage(fpData, groupBy(fpData, false));

            var dfData = {dfData};
            initDeploymentFrequencies(dfData, groupBy(dfData, true));
        });


        function initMeanTime(data, groupped) {
            for (var i = 0; i < groupped.length; i++) {
                var monthOfYear = Object.keys(groupped[i]).map(function (key) {
                    return {
                        month: key,
                        ...groupped[i][key],
                    };
                });

                monthOfYear = monthOfYear.filter(function (key) {
                    return !isNaN(parseInt(key.month));
                });
                monthOfYear.forEach(function (entry) {
                    var date = new Date(entry[entry.month][0].date);
                    data.push({
                        date: date.getFullYear() + "-" + (date.getMonth() + 1) + "-15",
                        monthValue: entry[entry.month].reduce((a, b) => a + b.value, 0),
                    });
                });
            }

            var monthData = data.filter(function (entry) {
                return !isNaN(parseInt(entry.monthValue))
            }).sort(function (a, b) {
                return new Date(a.date) - new Date(b.date);
            });
            data.push({
                date: monthData[0].date,
                averageValue: monthData[0].monthValue,
            });
            data.push({
                date: monthData[monthData.length - 1].date,
                averageValue: monthData.reduce((a, b) => a + b.monthValue, 0) / monthData.length,
            });

            am4core.useTheme(am4themes_animated);

            var chart = am4core.create("meanTimeDiv", am4charts.XYChart);

            chart.colors.step = 2;

            chart.data = data.sort(function (a, b) {
                return new Date(a.date) - new Date(b.date);
            });

            var dateAxis = chart.xAxes.push(new am4charts.DateAxis());
            dateAxis.renderer.minGridDistance = 50;

            function createAxisAndSeries(field, name, opposite, selectedBullet, hidden) {
                var valueAxis = chart.yAxes.push(new am4charts.ValueAxis());
                if (chart.yAxes.indexOf(valueAxis) != 0) {
                    valueAxis.syncWithAxis = chart.yAxes.getIndex(0);
                }

                var series = chart.series.push(new am4charts.LineSeries());
                series.dataFields.valueY = field;
                series.dataFields.dateX = "date";
                series.strokeWidth = 3;
                series.yAxis = valueAxis;
                series.name = name;
                series.tooltipText = "[bold]{valueY}[/] hours";
                series.showOnInit = true;
                series.hidden = hidden;

                var interfaceColors = new am4core.InterfaceColorSet();

                var bullet = {};
                switch (selectedBullet) {
                    case "triangle":
                        bullet = series.bullets.push(new am4charts.Bullet());
                        bullet.width = 12;
                        bullet.height = 12;
                        bullet.horizontalCenter = "middle";
                        bullet.verticalCenter = "middle";

                        var triangle = bullet.createChild(am4core.Triangle);
                        triangle.stroke = interfaceColors.getFor("background");
                        triangle.strokeWidth = 3;
                        triangle.direction = "top";
                        triangle.width = 12;
                        triangle.height = 12;
                        break;
                    case "rectangle":
                        bullet = series.bullets.push(new am4charts.Bullet());
                        bullet.width = 10;
                        bullet.height = 10;
                        bullet.horizontalCenter = "middle";
                        bullet.verticalCenter = "middle";

                        var rectangle = bullet.createChild(am4core.Rectangle);
                        rectangle.stroke = interfaceColors.getFor("background");
                        rectangle.strokeWidth = 3;
                        rectangle.width = 10;
                        rectangle.height = 10;
                        break;
                    default:
                        bullet = series.bullets.push(new am4charts.CircleBullet());
                        bullet.circle.stroke = interfaceColors.getFor("background");
                        bullet.circle.strokeWidth = 3;
                        break;
                }

                valueAxis.renderer.line.strokeOpacity = 1;
                valueAxis.renderer.line.strokeWidth = 3;
                valueAxis.renderer.line.stroke = series.stroke;
                valueAxis.renderer.labels.template.fill = series.stroke;
                valueAxis.renderer.opposite = opposite;
            }

            //initMeanTime
            createAxisAndSeries("value", "Daily", false, "circle", false);
            createAxisAndSeries("monthValue", "Monthly", false, "triangle", true);
            createAxisAndSeries("averageValue", "Average", false, "rectangle", false);

            chart.legend = new am4charts.Legend();
            chart.cursor = new am4charts.XYCursor();
        }

        function initLeadTime(data, groupped) {
            for (var i = 0; i < groupped.length; i++) {
                var monthOfYear = Object.keys(groupped[i]).map(function (key) {
                    return {
                        month: key,
                        ...groupped[i][key],
                    };
                });

                monthOfYear = monthOfYear.filter(function (key) {
                    return !isNaN(parseInt(key.month));
                });
                monthOfYear.forEach(function (entry) {
                    var date = new Date(entry[entry.month][0].date);
                    data.push({
                        date: date.getFullYear() + "-" + (date.getMonth() + 1) + "-15",
                        monthValue: entry[entry.month].reduce((a, b) => a + b.value, 0),
                    });
                });
            }

            var monthData = data.filter(function (entry) {
                return !isNaN(parseInt(entry.monthValue))
            }).sort(function (a, b) {
                return new Date(a.date) - new Date(b.date);
            });
            data.push({
                date: monthData[0].date,
                averageValue: monthData[0].monthValue,
            });
            data.push({
                date: monthData[monthData.length - 1].date,
                averageValue: monthData.reduce((a, b) => a + b.monthValue, 0) / monthData.length,
            });

            am4core.useTheme(am4themes_animated);

            var chart = am4core.create("leadTimeDiv", am4charts.XYChart);

            chart.colors.step = 2;
            chart.data = data.sort(function (a, b) {
                return new Date(a.date) - new Date(b.date);
            });

            var dateAxis = chart.xAxes.push(new am4charts.DateAxis());
            dateAxis.renderer.minGridDistance = 50;

            function createAxisAndSeries(field, name, opposite, selectedBullet, hidden) {
                var valueAxis = chart.yAxes.push(new am4charts.ValueAxis());
                if (chart.yAxes.indexOf(valueAxis) != 0) {
                    valueAxis.syncWithAxis = chart.yAxes.getIndex(0);
                }

                var series = chart.series.push(new am4charts.LineSeries());
                series.dataFields.valueY = field;
                series.dataFields.dateX = "date";
                series.strokeWidth = 3;
                series.yAxis = valueAxis;
                series.name = name;
                series.tooltipText = "[bold]{valueY}[/] hours";
                series.showOnInit = true;
                series.hidden = hidden;

                var interfaceColors = new am4core.InterfaceColorSet();

                var bullet = {};
                switch (selectedBullet) {
                    case "triangle":
                        bullet = series.bullets.push(new am4charts.Bullet());
                        bullet.width = 12;
                        bullet.height = 12;
                        bullet.horizontalCenter = "middle";
                        bullet.verticalCenter = "middle";

                        var triangle = bullet.createChild(am4core.Triangle);
                        triangle.stroke = interfaceColors.getFor("background");
                        triangle.strokeWidth = 3;
                        triangle.direction = "top";
                        triangle.width = 12;
                        triangle.height = 12;
                        break;
                    case "rectangle":
                        bullet = series.bullets.push(new am4charts.Bullet());
                        bullet.width = 10;
                        bullet.height = 10;
                        bullet.horizontalCenter = "middle";
                        bullet.verticalCenter = "middle";

                        var rectangle = bullet.createChild(am4core.Rectangle);
                        rectangle.stroke = interfaceColors.getFor("background");
                        rectangle.strokeWidth = 3;
                        rectangle.width = 10;
                        rectangle.height = 10;
                        break;
                    default:
                        bullet = series.bullets.push(new am4charts.CircleBullet());
                        bullet.circle.stroke = interfaceColors.getFor("background");
                        bullet.circle.strokeWidth = 3;
                        break;
                }

                valueAxis.renderer.line.strokeOpacity = 1;
                valueAxis.renderer.line.strokeWidth = 3;
                valueAxis.renderer.line.stroke = series.stroke;
                valueAxis.renderer.labels.template.fill = series.stroke;
                valueAxis.renderer.opposite = opposite;
            }

            //initLeadTime
            createAxisAndSeries("value", "Daily", false, "circle", false);
            createAxisAndSeries("monthValue", "Monthly", false, "triangle", true);
            createAxisAndSeries("averageValue", "Average", false, "rectangle", false);

            chart.legend = new am4charts.Legend();
            chart.cursor = new am4charts.XYCursor();
        }

        function initFailPercentage(data, groupped) {
            for (var i = 0; i < groupped.length; i++) {
                var monthOfYear = Object.keys(groupped[i]).map(function (key) {
                    return {
                        month: key,
                        ...groupped[i][key],
                    };
                });

                monthOfYear = monthOfYear.filter(function (key) {
                    return !isNaN(parseInt(key.month));
                });
                monthOfYear.forEach(function (entry) {
                    var date = new Date(entry[entry.month][0].date);
                    data.push({
                        date: date.getFullYear() + "-" + (date.getMonth() + 1) + "-15",
                        monthValue: entry[entry.month].reduce((a, b) => a + b.value, 0),
                    });
                });
            }

            var monthData = data.filter(function (entry) {
                return !isNaN(parseInt(entry.monthValue))
            }).sort(function (a, b) {
                return new Date(a.date) - new Date(b.date);
            });
            data.push({
                date: monthData[0].date,
                averageValue: monthData[0].monthValue,
            });
            data.push({
                date: monthData[monthData.length - 1].date,
                averageValue: monthData.reduce((a, b) => a + b.monthValue, 0) / monthData.length,
            });

            am4core.useTheme(am4themes_animated);

            var chart = am4core.create("failPercentageDiv", am4charts.XYChart);

            chart.colors.step = 2;
            chart.data = data.sort(function (a, b) {
                return new Date(a.date) - new Date(b.date);
            });

            var dateAxis = chart.xAxes.push(new am4charts.DateAxis());
            dateAxis.renderer.minGridDistance = 50;

            function createAxisAndSeries(field, name, opposite, selectedBullet, hidden) {
                var valueAxis = chart.yAxes.push(new am4charts.ValueAxis());
                if (chart.yAxes.indexOf(valueAxis) != 0) {
                    valueAxis.syncWithAxis = chart.yAxes.getIndex(0);
                }

                var series = chart.series.push(new am4charts.LineSeries());
                series.dataFields.valueY = field;
                series.dataFields.dateX = "date";
                series.strokeWidth = 3;
                series.yAxis = valueAxis;
                series.name = name;
                series.tooltipText = "%[bold]{valueY}[/]";
                series.showOnInit = true;
                series.hidden = hidden;

                var interfaceColors = new am4core.InterfaceColorSet();

                var bullet = {};
                switch (selectedBullet) {
                    case "triangle":
                        bullet = series.bullets.push(new am4charts.Bullet());
                        bullet.width = 12;
                        bullet.height = 12;
                        bullet.horizontalCenter = "middle";
                        bullet.verticalCenter = "middle";

                        var triangle = bullet.createChild(am4core.Triangle);
                        triangle.stroke = interfaceColors.getFor("background");
                        triangle.strokeWidth = 3;
                        triangle.direction = "top";
                        triangle.width = 12;
                        triangle.height = 12;
                        break;
                    case "rectangle":
                        bullet = series.bullets.push(new am4charts.Bullet());
                        bullet.width = 10;
                        bullet.height = 10;
                        bullet.horizontalCenter = "middle";
                        bullet.verticalCenter = "middle";

                        var rectangle = bullet.createChild(am4core.Rectangle);
                        rectangle.stroke = interfaceColors.getFor("background");
                        rectangle.strokeWidth = 3;
                        rectangle.width = 10;
                        rectangle.height = 10;
                        break;
                    default:
                        bullet = series.bullets.push(new am4charts.CircleBullet());
                        bullet.circle.stroke = interfaceColors.getFor("background");
                        bullet.circle.strokeWidth = 3;
                        break;
                }

                valueAxis.renderer.line.strokeOpacity = 1;
                valueAxis.renderer.line.strokeWidth = 3;
                valueAxis.renderer.line.stroke = series.stroke;
                valueAxis.renderer.labels.template.fill = series.stroke;
                valueAxis.renderer.opposite = opposite;
            }
            //FailPercentage
            createAxisAndSeries("value", "Daily", false, "circle", false);
            createAxisAndSeries("monthValue", "Monthly", false, "triangle", true);
            createAxisAndSeries("averageValue", "Average", false, "rectangle", false);

            chart.legend = new am4charts.Legend();
            chart.cursor = new am4charts.XYCursor();
        }

        function initDeploymentFrequencies(data, groupped) {
            for (var i = 0; i < groupped.length; i++) {
                var monthOfYear = Object.keys(groupped[i]).map(function (key) {
                    return {
                        month: key,
                        ...groupped[i][key],
                    };
                });

                monthOfYear = monthOfYear.filter(function (key) {
                    return !isNaN(parseInt(key.month));
                });
                monthOfYear.forEach(function (entry) {
                    var date = new Date(entry[entry.month][0].date);
                    data.push({
                        date: date.getFullYear() + "-" + (date.getMonth() + 1) + "-15",
                        monthValue: entry[entry.month].reduce((a, b) => a + b.value, 0),
                    });
                });
            }

            var monthData = data.filter(function (entry) {
                return !isNaN(parseInt(entry.monthValue))
            }).sort(function (a, b) {
                return new Date(a.date) - new Date(b.date);
            });
            data.push({
                date: monthData[0].date,
                averageValue: monthData[0].monthValue,
            });
            data.push({
                date: monthData[monthData.length - 1].date,
                averageValue: monthData.reduce((a, b) => a + b.monthValue, 0) / monthData.length,
            });

            am4core.useTheme(am4themes_animated);

            var chart = am4core.create("deploymentFrequencyDiv", am4charts.XYChart);

            chart.colors.step = 2;
            chart.data = data.sort(function (a, b) {
                return new Date(a.date) - new Date(b.date);
            });

            var dateAxis = chart.xAxes.push(new am4charts.DateAxis());
            dateAxis.renderer.minGridDistance = 50;

            function createAxisAndSeries(field, name, opposite, selectedBullet, hidden) {
                var valueAxis = chart.yAxes.push(new am4charts.ValueAxis());
                if (chart.yAxes.indexOf(valueAxis) != 0) {
                    valueAxis.syncWithAxis = chart.yAxes.getIndex(0);
                }

                var series = chart.series.push(new am4charts.LineSeries());
                series.dataFields.valueY = field;
                series.dataFields.dateX = "date";
                series.strokeWidth = 3;
                series.yAxis = valueAxis;
                series.name = name;
                series.tooltipText = "{name}: [bold]{valueY}[/]";
                series.showOnInit = true;
                series.hidden = hidden;

                var interfaceColors = new am4core.InterfaceColorSet();

                var bullet = {};
                switch (selectedBullet) {
                    case "triangle":
                        bullet = series.bullets.push(new am4charts.Bullet());
                        bullet.width = 12;
                        bullet.height = 12;
                        bullet.horizontalCenter = "middle";
                        bullet.verticalCenter = "middle";

                        var triangle = bullet.createChild(am4core.Triangle);
                        triangle.stroke = interfaceColors.getFor("background");
                        triangle.strokeWidth = 3;
                        triangle.direction = "top";
                        triangle.width = 12;
                        triangle.height = 12;
                        break;
                    case "rectangle":
                        bullet = series.bullets.push(new am4charts.Bullet());
                        bullet.width = 10;
                        bullet.height = 10;
                        bullet.horizontalCenter = "middle";
                        bullet.verticalCenter = "middle";

                        var rectangle = bullet.createChild(am4core.Rectangle);
                        rectangle.stroke = interfaceColors.getFor("background");
                        rectangle.strokeWidth = 3;
                        rectangle.width = 10;
                        rectangle.height = 10;
                        break;
                    default:
                        bullet = series.bullets.push(new am4charts.CircleBullet());
                        bullet.circle.stroke = interfaceColors.getFor("background");
                        bullet.circle.strokeWidth = 3;
                        break;
                }

                valueAxis.renderer.line.strokeOpacity = 1;
                valueAxis.renderer.line.strokeWidth = 3;
                valueAxis.renderer.line.stroke = series.stroke;
                valueAxis.renderer.labels.template.fill = series.stroke;
                valueAxis.renderer.opposite = opposite;
            }

            //initDeploymentFrequencies
            createAxisAndSeries("value", "Daily", false, "circle", false);
            createAxisAndSeries("monthValue", "Monthly", false, "triangle", true);
            createAxisAndSeries("averageValue", "Average", false, "rectangle", false);

            chart.legend = new am4charts.Legend();
            chart.cursor = new am4charts.XYCursor();
        }

        function groupBy(data, isDf) {
            let res = {};

            let fn = (o, year, month, array) => {
                o[year][month] = {
                    [month]: data.filter(({date: d}) => (year + "-" + month) === d.slice(0, 7))
                };
            };

            for (var i = 0; i < data.length; i++) {
                if (isDf && data.length > i + 1 && data[i].date === data[i + 1].date) {
                    data[i].value += data[i + 1].value;
                }

                let [year, month] = data[i].date.match(/\d+/g);
                if (!res[year]) res[year] = {};
                fn(res, year, month, data);

                if (isDf && data.length > i + 1 && data[i].date === data[i + 1].date) {
                    data.splice(i + 1, 1);
                    i--;
                }
            }

            return Object.keys(res).map(function (key) {
                return {
                    year: parseInt(key),
                    ...res[key]
                };
            });
        }
    </script>

    <style>
        body {
            font-family: 'Rubik', sans-serif;
        }

        div {
            width: 100%;
            height: 300px;
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
    </style>
</head>
<body>
<div class="container">
    <h4 class="title">four-key Metrics</h4>
    <h6 class="subtitle">{repositoryName} <span>|</span> {teamName} <span>|</span> {startDate} - {endDate}</h6>

    <h6 class="chart-title">Deployment Frequencies</h6>
    <div id="deploymentFrequencyDiv"></div>

    <h6 class="chart-title">Lead Times</h6>
    <div id="leadTimeDiv"></div>

    <h6 class="chart-title">Mean Times</h6>
    <div id="meanTimeDiv"></div>

    <h6 class="chart-title">Fail Percentages</h6>
    <div id="failPercentageDiv"></div>
</div>
</body>
</html>



	`
}
