package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"zcash-obs-db/sql"
	"zcash-obs-db/util"
)

func HandleCurrentHeight(w http.ResponseWriter, r *http.Request) {
	height := sql.SQLSelectCurrentHeight()

	fmt.Fprintf(w, strconv.Itoa(height))
}

func HandleGetBlockByHash(w http.ResponseWriter, r *http.Request) {
	// parse query
	var mQueryValues map[string][]string = r.URL.Query()

	var err error
	hash := "NIL"
	if len(mQueryValues["hash"]) != 0 {
		hash = mQueryValues["hash"][0]
	}

	var block util.BlockMessage = sql.SQLSelectBlockByHash(hash)

	raw_json, err := json.Marshal(block)
	if err != nil {
		fmt.Fprintf(w, "500")
		return
	}

	fmt.Fprintf(w, string(raw_json))
}

func HandleGetBlocksByHeightRange(w http.ResponseWriter, r *http.Request) {
	// parse query
	var mQueryValues map[string][]string = r.URL.Query()

	var err error
	min_height := 0
	if len(mQueryValues["min_height"]) != 0 {
		min_height, err = strconv.Atoi(mQueryValues["min_height"][0])
		if err != nil {
			fmt.Fprintf(w, "422")
			return
		}
	}

	max_height := 0
	if len(mQueryValues["max_height"]) != 0 {
		max_height, err = strconv.Atoi(mQueryValues["max_height"][0])
		if err != nil {
			fmt.Fprintf(w, "422")
			return
		}
	}

	blocks := sql.SQLSelectBlocksByHeightRange(min_height, max_height)

	raw_json, err := json.Marshal(blocks)
	if err != nil {
		fmt.Fprintf(w, "500")
		return
	}

	fmt.Fprintf(w, string(raw_json))
}

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
