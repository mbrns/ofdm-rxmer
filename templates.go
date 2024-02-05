package main

const table = `
<!DOCTYPE html>
<html lang="en">
<head>
	<style>
	table.dataTable thead .sorting:after,
	table.dataTable thead .sorting:before,
	table.dataTable thead .sorting_asc:after,
	table.dataTable thead .sorting_asc:before,
	table.dataTable thead .sorting_asc_disabled:after,
	table.dataTable thead .sorting_asc_disabled:before,
	table.dataTable thead .sorting_desc:after,
	table.dataTable thead .sorting_desc:before,
	table.dataTable thead .sorting_desc_disabled:after,
	table.dataTable thead .sorting_desc_disabled:before {
	bottom: .5em;
	}
	</style>
    <meta charset="UTF-8">
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">
	<link rel="stylesheet" type="text/css" href="https://cdn.datatables.net/v/bs4/dt-1.10.20/datatables.min.css"/>

    <title>rxMER Results</title>
</head>
<body>

<div class="container">
    <table id="dtBasicExample" class="table">
        <thead>
            <tr>
                <th scope="col">Frequency</th>
                <th scope="col">MER</th>
            </tr>
        </thead>
        <tbody>
            {{ range $key, $value := . }}
            <tr>
                <th scope="row">{{$key}}</th>
                <td>{{$value}} dB</td>
            </tr>
            {{ end }}
        </tbody>
    </table>
</div>


<script src="https://code.jquery.com/jquery-3.4.1.slim.min.js" integrity="sha384-J6qa4849blE2+poT4WnyKhv5vZF5SrPo0iEjwBvKU7imGFAV0wwj1yYfoRSJoZ+n" crossorigin="anonymous"></script>
<script src="https://cdn.jsdelivr.net/npm/popper.js@1.16.0/dist/umd/popper.min.js" integrity="sha384-Q6E9RHvbIyZFJoft+2mJbHaEWldlvI9IOYy5n3zV9zzTtmI3UksdQRVvoxMfooAo" crossorigin="anonymous"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/js/bootstrap.min.js" integrity="sha384-wfSDF2E50Y2D1uUdj0O3uMBJnjuUD4Ih7YwaYd1iqfktj0Uod8GCExl3Og8ifwB6" crossorigin="anonymous"></script>
<script type="text/javascript" src="https://cdn.datatables.net/v/bs4/dt-1.10.20/datatables.min.js"></script>


<script>
$(document).ready(function () {
$('#dtBasicExample').DataTable({
	"paging": false
});
$('.dataTables_length').addClass('bs-select');
});
</script>

</body>

</html>
`

const chart = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">
    <title>rxMER Results</title>
</head>
<body>

<div id="chartContainer" style="height: 300px; width: 100%;"></div>

<script src="https://code.jquery.com/jquery-3.4.1.slim.min.js" integrity="sha384-J6qa4849blE2+poT4WnyKhv5vZF5SrPo0iEjwBvKU7imGFAV0wwj1yYfoRSJoZ+n" crossorigin="anonymous"></script>
<script src="https://cdn.jsdelivr.net/npm/popper.js@1.16.0/dist/umd/popper.min.js" integrity="sha384-Q6E9RHvbIyZFJoft+2mJbHaEWldlvI9IOYy5n3zV9zzTtmI3UksdQRVvoxMfooAo" crossorigin="anonymous"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/js/bootstrap.min.js" integrity="sha384-wfSDF2E50Y2D1uUdj0O3uMBJnjuUD4Ih7YwaYd1iqfktj0Uod8GCExl3Og8ifwB6" crossorigin="anonymous"></script>
<script type="text/javascript" src="https://canvasjs.com/assets/script/canvasjs.min.js"></script>


<script type="text/javascript">
    window.onload = function () {

        var chart = new CanvasJS.Chart("chartContainer",
            {
                zoomEnabled: true,
                title:{
                    text: "rxMER per subcarrier"
                },
                axisY:{
                    includeZero: false
                },
                data: data,

            });

        chart.render();
    };

    var limit = 4096;    //increase number of dataPoints by increasing this

    var y = 0;
    var data = [];
    var dataSeries = { type: "line" };
    var dataPoints = [];
    dataPoints.push(
        {{ range $key, $value := . }}
        {
            x: {{$key}},
            y: {{$value}}
        },
        {{end}}
    );

    dataSeries.dataPoints = dataPoints;
    data.push(dataSeries);

</script>


</body>

</html>
`
