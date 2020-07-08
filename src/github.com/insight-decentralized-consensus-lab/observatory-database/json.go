package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Blocks
func HandleGetRecentBlocks(w http.ResponseWriter, r *http.Request) {
	var err error

	// parse query
	var mQueryValues map[string][]string = r.URL.Query()

	var n int = 5
	if len(mQueryValues["n"]) != 0 {
		n, err = strconv.Atoi(mQueryValues["n"][0])
		if err != nil {
			fmt.Fprintf(w, "ERROR: invalid argument")
			return
		}
	}

	var blocks []BlockMessage = SQLSelectRecentBlocks(n)

	var blocks_json []byte
	blocks_json, err = json.Marshal(blocks)
	if err != nil {
		fmt.Fprintf(w, "ERROR: invalid blocks")
		return
	}

	fmt.Fprintf(w, string(blocks_json))
}

func HandleGetRangeBlocks(w http.ResponseWriter, r *http.Request) {
	var err error

	// parse query
	var mQueryValues map[string][]string = r.URL.Query()

	var max_height int = 0
	if len(mQueryValues["max_height"]) != 0 {
		max_height, err = strconv.Atoi(mQueryValues["max_height"][0])
		if err != nil {
			fmt.Fprintf(w, "ERROR: invalid 'max_height' argument")
			return
		}
	}

	var min_height int = 0
	if len(mQueryValues["min_height"]) != 0 {
		min_height, err = strconv.Atoi(mQueryValues["min_height"][0])
		if err != nil {
			fmt.Fprintf(w, "ERROR: invalid 'min_height' argument")
			return
		}
	}

	fmt.Printf("MIN: %d - MAX: %d\n", min_height, max_height)

	var blocks []BlockMessage = SQLSelectRangeBlocks(min_height, max_height)

	fmt.Println(blocks)

	var blocks_json []byte
	blocks_json, err = json.Marshal(blocks)
	if err != nil {
		fmt.Fprintf(w, "ERROR: invalid blocks")
		return
	}

	fmt.Fprintf(w, string(blocks_json))
}

// Forks
func HandleGetRecentForks(w http.ResponseWriter, r *http.Request) {
	var err error

	// parse query
	var mQueryValues map[string][]string = r.URL.Query()

	var n int = 5
	if len(mQueryValues["n"]) != 0 {
		n, err = strconv.Atoi(mQueryValues["n"][0])
		if err != nil {
			fmt.Fprintf(w, "ERROR: invalid 'n' argument")
			return
		}
	}

	var forks []ForkMessage = SQLSelectRecentForks(n)

	var forks_json []byte
	forks_json, err = json.Marshal(forks)
	if err != nil {
		fmt.Fprintf(w, "ERROR: invalid blocks")
		return
	}

	fmt.Fprintf(w, string(forks_json))
}

func HandleGetRangeForks(w http.ResponseWriter, r *http.Request) {
	var err error

	// parse query
	var mQueryValues map[string][]string = r.URL.Query()

	var max_height int = 0
	if len(mQueryValues["max_height"]) != 0 {
		max_height, err = strconv.Atoi(mQueryValues["max_height"][0])
		if err != nil {
			fmt.Fprintf(w, "ERROR: invalid 'max_height' argument")
			return
		}
	}

	var min_height int = 0
	if len(mQueryValues["min_height"]) != 0 {
		min_height, err = strconv.Atoi(mQueryValues["min_height"][0])
		if err != nil {
			fmt.Fprintf(w, "ERROR: invalid 'min_height' argument")
			return
		}
	}

	var forks []ForkMessage = SQLSelectRangeForks(min_height, max_height)

	var forks_json []byte
	forks_json, err = json.Marshal(forks)
	if err != nil {
		fmt.Fprintf(w, "ERROR: invalid forks")
	}

	fmt.Fprintf(w, string(forks_json))
}

func StartJSONService() {

	http.HandleFunc("/recentblocks", HandleGetRecentBlocks)
	http.HandleFunc("/rangeblocks", HandleGetRangeBlocks)

	http.HandleFunc("/recentforks", HandleGetRecentForks)
	http.HandleFunc("/rangeforks", HandleGetRangeForks)

	http.ListenAndServe(":8080", nil)
}
