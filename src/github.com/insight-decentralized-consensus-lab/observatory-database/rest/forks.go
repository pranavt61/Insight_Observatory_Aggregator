package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"zcash-obs-db/sql"
	"zcash-obs-db/util"
)

func HandleRecentForks(w http.ResponseWriter, r *http.Request) {
	// parse query
	var mQueryValues map[string][]string = r.URL.Query()

	var err error
	nNumForks := 5
	if len(mQueryValues["n"]) != 0 {
		nNumForks, err = strconv.Atoi(mQueryValues["n"][0])
		if err != nil {
			fmt.Fprintf(w, "422")
			return
		}
	}

	var forks []util.ForkMessage = sql.SQLSelectRecentForks(nNumForks)

	raw_json, err := json.Marshal(forks)
	if err != nil {
		fmt.Fprintf(w, "500")
		return
	}

	fmt.Fprintf(w, string(raw_json))
}

func HandleRangeForks(w http.ResponseWriter, r *http.Request) {
	// parse query
	var mQueryValues map[string][]string = r.URL.Query()

	var err error
	nMinHeight := 0
	if len(mQueryValues["min_height"]) != 0 {
		nMinHeight, err = strconv.Atoi(mQueryValues["min_height"][0])
		if err != nil {
			fmt.Fprintf(w, "422")
			return
		}
	} else {
		fmt.Fprintf(w, "422")
	}

	nMaxHeight := 0
	if len(mQueryValues["max_height"]) != 0 {
		nMaxHeight, err = strconv.Atoi(mQueryValues["max_height"][0])
		if err != nil {
			fmt.Fprintf(w, "422")
			return
		}
	} else {
		fmt.Fprintf(w, "422")
	}

	var forks []util.ForkMessage = sql.SQLSelectRangeForks(nMinHeight, nMaxHeight)

	raw_json, err := json.Marshal(forks)
	if err != nil {
		fmt.Fprintf(w, "500")
		return
	}

	fmt.Fprintf(w, string(raw_json))
}
