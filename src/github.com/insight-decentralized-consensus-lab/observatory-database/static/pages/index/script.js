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

    console.log("Tables Updated");
};

let updateCharts = () => {

    // Block Prop Time Chart
    let BlockPropTimeCtx = document.getElementById('ChartPropTime').getContext('2d');
    $.get("/v1/json/allinvhalfrange", function(data, status){
        let times = JSON.parse(data);

        let num_bars = 10;
        let bars = [0,0,0,0,0,0,0,0,0,0];
        let max_time = Math.max(...times);
        for (let i = 0; i < times.length; i ++) {
            bars[Math.min(Math.floor(bars.length * times[i]/max_time), num_bars - 1)] ++;
        }

        let labels = [];
        for (let i = 0; i < bars.length; i ++) {
            let l = Math.floor(i * (max_time) / num_bars) / 1000;
            labels.push(l.toString());
        }

        let config = {
            type: 'bar',
            data: {
                labels: labels,
                datasets: [{
                    label: 'Observed Block Propigation Time',
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

        let myChart = new Chart(BlockPropTimeCtx, config);
    });

    // Block Size Chart
    let BlockSizeCtx = document.getElementById('ChartBlockSize').getContext('2d');
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

        let myChart = new Chart(BlockSizeCtx, config);
    });

    let ForkFreqCtx = document.getElementById('ChartForkFreq').getContext('2d');
    $.get("/v1/json/currentheight", function(data, status){
        let max_height = parseInt(data);

        let min_height = 889208

        $.get("/v1/json/rangeforks?min_height=" + min_height + "&max_height=" + max_height, function(data,status){
            let forks = JSON.parse(data);

            let heights = [];
            for (let i = 0; i < forks.length; i ++) {
                heights.push(forks[i].height);
            }

            let num_bars = 10;
            let bars = [0,0,0,0,0,0,0,0,0,0];
            for (let i = 0; i < heights.length; i ++) {
                bars[Math.min(Math.floor(num_bars * ((heights[i] - min_height) / (max_height - min_height))), num_bars - 1)] ++;
            }

            let labels = [];
            for (let i = 0; i < bars.length; i ++) {
                let l = min_height + Math.floor(i * (max_height - min_height) / num_bars);
                labels.push(l.toString());
            }

            let config = {
                type: 'bar',
                data: {
                    labels: labels,
                    datasets: [{
                        label: 'Fork Freq by Height',
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

            let myChart = new Chart(ForkFreqCtx, config);

        });
    });
};

$(document).ready(function() {

    $("time.timeago").timeago();

    updateTables();
    updateCharts();

    window.setInterval(updateTables, 5000);
});


