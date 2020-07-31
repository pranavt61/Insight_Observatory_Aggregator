
$(document).ready(function() {
    $("time.timeago").timeago();

    let searchParams = new URLSearchParams(window.location.search);
    if (searchParams.has("min_height") && searchParams.has("max_height")) {
        min_height = searchParams.get("min_height");
        max_height = searchParams.get("max_height");

        // set title
        $("#TitleHeader").append("Fork " + min_height);

        // make request
        $.get("/v1/json/getblocksbyheightrange?min_height=" + min_height + "&max_height=" + max_height, function(data, status) {
            let fork = JSON.parse(data);
            fork.sort((a, b) => {return a.network_time - b.network_time});
            console.log(fork);

            let block_list = "";

            let main_chain = "";
            for (let i = fork.length - 1; i >= 0; i --) {

                let time = new Date(fork[i].network_time);

                if (main_chain == "" || fork[i].hash == main_chain) {
                    // main chain
                    block_list += `
                        <tr class="table-success">
                            <th scope="row">` + fork[i].height + `</th>

                            <td><a class="text-secondary" href="/block?hash=` + fork[i].hash + `">` + fork[i].hash + `</a></td>
                            <td>` + time.toString() + `</td>
                        </tr>`
                    main_chain = fork[i].prev_hash;
                } else {
                    // fork
                    block_list += `
                        <tr class="table-danger">
                            <th scope="row">` + fork[i].height + `</th>
                            <td><a class="text-secondary" href="/block?hash=` + fork[i].hash + `">` + fork[i].hash + `</a></td>
                            <td>` + time.toString() + `</td>
                        </tr>`
                }
            }

            $("#BlockTable").append(block_list);
        });

    } else {
        //ERR
        $("#TitleHeader").append("404");
    }
});


