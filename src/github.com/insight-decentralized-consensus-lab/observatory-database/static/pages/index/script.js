let updateTables = () => {

    // Recent Blocks Table
    $.get("/v1/json/recentblocks?n=5&with_inv=true", function(data, status) {
        let blocks = JSON.parse(data).reverse();

        $("#RecentBlocksTable").empty();

        for (let i = 0; i < blocks.length; i ++) {

            // prop time
            blocks[i].inv.sort((a,b) => {return a.network_time - b.network_time});

            let min_inv = 0;
            let max_inv = 0;
            for (let j = 0; j < blocks[i].inv.length; j ++) {
                if (min_inv == 0 && max_inv == 0) {
                    min_inv = blocks[i].inv[j].network_time;
                    max_inv = blocks[i].inv[j].network_time;

                    continue;
                }


                if (blocks[i].inv[j].network_time < min_inv)
                    min_inv = blocks[i].inv[j].network_time;
                if (blocks[i].inv[j].network_time > max_inv)
                    max_inv = blocks[i].inv[j].network_time;
            }

            let row_HTML = "<tr>";
            row_HTML += "<th scope=\"row\"> <a href=\"/block?hash=" + blocks[i].hash + "\" class=\"text-secondary\">" + blocks[i].height + "</a></th>";
            row_HTML += "<td>" + $.timeago(new Date(blocks[i].network_time)) + "</td>";
            row_HTML += "<td>" + (max_inv - min_inv) + " ms</td>";
            row_HTML += "</tr>";

            $("#RecentBlocksTable").append(row_HTML);
        }
    });

    // Recent Forks Table
    $.get("/v1/json/recentforks", function(data, status) {
        let forks = JSON.parse(data).reverse();

        $("#RecentForksTable").empty();

        for (let i = 0; i < forks.length; i ++) {

            let row_HTML = "<tr>";
            row_HTML += "<th scope=\"row\"><a href=\"/fork?min_height=" + (forks[i].height - 1) + "&max_height=" + (forks[i].height + 2) + "\" class=\"text-secondary\">" + forks[i].height + "</a></th>";
            row_HTML += "<td>" + $.timeago(new Date(forks[i].blocks[0].network_time)) + "</td>";
            row_HTML += "</tr>";

            $("#RecentForksTable").append(row_HTML);
        }

    });

    // Block Size Chart
    let ctx = document.getElementById('myChart').getContext('2d');
    $.get("/v1/json/recentblocks?n=100&with_inv=false", function(data, status){
        let blocks = JSON.parse(data);

        let sizes = [];
        for (let i = 0; i < blocks.length; i ++) {
            sizes.push(blocks[i].block_size);
        }

        let num_bars = 10;
        let bars = [0,0,0,0,0,0,0,0,0,0];
        let max_size = Math.max(...sizes);
        for (let i = 0; i < sizes.length; i ++) {
            bars[Math.min(Math.floor(bars.length * sizes[i]/max_size), num_bars - 1)] ++;
        }

        let config = {
            type: 'bar',
            data: {
                labels: ['', '', '', '', '', '', '', '', '', ''],
                datasets: [{
                    label: 'Block Sizes of Last 100',
                    data: bars,
                    backgroundColor: [
                        'rgba(255, 99, 132, 0.2)',
                        'rgba(54, 162, 235, 0.2)',
                        'rgba(255, 206, 86, 0.2)',
                        'rgba(75, 192, 192, 0.2)',
                        'rgba(153, 102, 255, 0.2)',
                        'rgba(255, 159, 64, 0.2)',
                        'rgba(255, 206, 86, 0.2)',
                        'rgba(75, 192, 192, 0.2)',
                        'rgba(153, 102, 255, 0.2)',
                        'rgba(255, 159, 64, 0.2)'

                    ],
                    borderColor: [
                        'rgba(255, 99, 132, 1)',
                        'rgba(54, 162, 235, 1)',
                        'rgba(255, 206, 86, 1)',
                        'rgba(75, 192, 192, 1)',
                        'rgba(153, 102, 255, 1)',
                        'rgba(255, 159, 64, 1)',
                        'rgba(255, 206, 86, 1)',
                        'rgba(75, 192, 192, 1)',
                        'rgba(153, 102, 255, 1)',
                        'rgba(255, 159, 64, 1)'
                    ],
                    borderWidth: 1
                }]
            },
            options: {
                scales: {
                    yAxes: [{
                        ticks: {
                            beginAtZero: true
                        }
                    }]
                }
            }
        };

        let myChart = new Chart(ctx, config);

    });

    console.log("Tables Updated");
}

$(document).ready(function() {

    $("time.timeago").timeago();

    updateTables();

    // window.setInterval(updateTables, 5000);
});


