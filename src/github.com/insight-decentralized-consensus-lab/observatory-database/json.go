package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

var JSONServer http.Server

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

	var blocks []BlockMessage = SQLSelectRangeBlocks(min_height, max_height)

	var blocks_json []byte
	blocks_json, err = json.Marshal(blocks)
	if err != nil {
		fmt.Fprintf(w, "ERROR: invalid blocks")
		return
	}

	fmt.Fprintf(w, string(blocks_json))
}

func HandleGetHashBlocks(w http.ResponseWriter, r *http.Request) {
	var err error

	// parse query
	var mQueryValues map[string][]string = r.URL.Query()

	var hash string = ""
	if len(mQueryValues["hash"]) != 0 {
		hash = mQueryValues["hash"][0]
	}

	var blocks []BlockMessage = SQLSelectHashBlocks(hash)

	var blocks_json []byte
	blocks_json, err = json.Marshal(blocks)
	if err != nil {
		fmt.Fprintf(w, "ERROR: invalid blocks")
		return
	}

	fmt.Fprintf(w, string(blocks_json))
}

// Inv
func HandleGetRecentInv(w http.ResponseWriter, r *http.Request) {
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

	var inv []InvMessage = SQLSelectRecentInv(n)

	var inv_json []byte
	inv_json, err = json.Marshal(inv)
	if err != nil {
		fmt.Fprintf(w, "ERROR: invalid blocks")
		return
	}

	fmt.Fprintf(w, string(inv_json))
}

func HandleGetHashInv(w http.ResponseWriter, r *http.Request) {
	var err error

	// parse query
	var mQueryValues map[string][]string = r.URL.Query()

	var hash string = ""
	if len(mQueryValues["hash"]) != 0 {
		hash = mQueryValues["hash"][0]
	}

	var inv []InvMessage = SQLSelectHashInv(hash)

	var inv_json []byte
	inv_json, err = json.Marshal(inv)
	if err != nil {
		fmt.Fprintf(w, "ERROR: invalid blocks")
		return
	}

	fmt.Fprintf(w, string(inv_json))
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

func StartJSONServer() {

	JSONServer = http.Server{Addr: ":8080"}

	http.HandleFunc("/recentblocks", HandleGetRecentBlocks)
	http.HandleFunc("/rangeblocks", HandleGetRangeBlocks)
	http.HandleFunc("/hashblocks", HandleGetHashBlocks)

	http.HandleFunc("/recentinv", HandleGetRecentInv)
	http.HandleFunc("/hashinv", HandleGetHashInv)

	http.HandleFunc("/recentforks", HandleGetRecentForks)
	http.HandleFunc("/rangeforks", HandleGetRangeForks)

	go func() {
		if err := JSONServer.ListenAndServe(); err != nil {
			fmt.Printf("JSON server Shutdown: %s\n", err.Error())
		}
	}()
}

func StopJSONServer() {

	if err := JSONServer.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		fmt.Printf("HTTP server Shutdown: %s\n", err.Error())
	}
}
