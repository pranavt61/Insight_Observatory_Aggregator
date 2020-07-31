
$(document).ready(function() {

    let searchParams = new URLSearchParams(window.location.search);
    if (searchParams.has("hash")) {
        block_hash = searchParams.get("hash");
        
        $.get("/v1/json/getblockbyhash?hash=" + block_hash, function(data, status){
            var block = JSON.parse(data);

            if (block.height == 0) {
                // ERR
                $("#TitleHeader").append("404")
                return;
            }

            $("#TitleHeader").append("Block " + block.height)

            var block_table = ""
            var formatTableKey = (string) => {
                // replace "_" with " "
                string = string.replace("_", " ");

                string = string.replace("tx", "TX");
                // capitilize
                return string.charAt(0).toUpperCase() + string.slice(1);
            };

            for (var key in block) {
                if (block.hasOwnProperty(key)) {
                    if (key == "height") {
                        continue;
                    }
                    block_table += `<div class="row">
                                    <div id="TableHeader" class="col-2">
                                    ` + formatTableKey(key) + `
                                    </div>
                                    <div class="col-10">
                                    ` + block[key] + `
                                    </div>
                                </div>`

              }
            }
            $("#BlockTable").append(block_table)
        });

    } else {
        // ERR
        
        $("#TitleHeader").append("404")
    }



});


