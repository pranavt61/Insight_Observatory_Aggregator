package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"zcash-obs-db/sql"
	"zcash-obs-db/util"
)

func HandleRecentBlocks(w http.ResponseWriter, r *http.Request) {
	// parse query
	var mQueryValues map[string][]string = r.URL.Query()

	var err error
	nNumBlocks := 5
	if len(mQueryValues["n"]) != 0 {
		nNumBlocks, err = strconv.Atoi(mQueryValues["n"][0])
		if err != nil {
			fmt.Fprintf(w, "422")
			return
		}
	}

	withInv := false
	if len(mQueryValues["with_inv"]) != 0 {
		withInv, err = strconv.ParseBool(mQueryValues["with_inv"][0])
		if err != nil {
			fmt.Fprintf(w, "422")
			return
		}
	}

	var blocks []util.BlockMessage

	if withInv {
		blocks = sql.SQLSelectRecentBlocksWithInv(nNumBlocks)
	} else {
		blocks = sql.SQLSelectRecentBlocks(nNumBlocks)
	}

	raw_json, err := json.Marshal(blocks)
	if err != nil {
		fmt.Fprintf(w, "500")
		return
	}

	fmt.Fprintf(w, string(raw_json))
}
