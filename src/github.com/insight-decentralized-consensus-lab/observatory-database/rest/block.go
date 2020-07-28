package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"zcash-obs-db/sql"
	"zcash-obs-db/util"
)

func HandleRecentBlocks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HEYYOOOOOO")
	var blocks []util.BlockMessage = sql.SQLSelectRecentBlocks(5)
	fmt.Println(blocks)

	raw_json, err := json.Marshal(blocks)
	if err != nil {
		fmt.Fprintf(w, "500")
		return
	}

	fmt.Fprintf(w, string(raw_json))
}
