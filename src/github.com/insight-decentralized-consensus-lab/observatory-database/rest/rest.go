package rest

import (
	"zcash-obs-db/sql"
	"zcash-obs-db/util"

	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var JSONServer http.Server

// Index HTML File
func HandleGetFileIndex(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadFile("./static/pages/index/page.html")
	if err != nil {
		fmt.Fprintf(w, "ERROR: reading index")
		return
	}

	fmt.Fprintf(w, string(body))
}

// Block HTML File
func HandleGetFileBlock(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadFile("./static/pages/block/page.html")
	if err != nil {
		fmt.Fprintf(w, "ERROR: reading block")
		return
	}

	fmt.Fprintf(w, string(body))
}

// Fork HTML File
func HandleGetFileFork(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadFile("./static/pages/fork/page.html")
	if err != nil {
		fmt.Fprintf(w, "ERROR: reading block")
		return
	}

	fmt.Fprintf(w, string(body))
}

// Blocks

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

	var blocks []util.BlockMessage = sql.SQLSelectRangeBlocks(min_height, max_height)

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

	var blocks []util.BlockMessage = sql.SQLSelectHashBlocks(hash)

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

	var inv []util.InvMessage = sql.SQLSelectRecentInv(n)

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

	var inv []util.InvMessage = sql.SQLSelectHashInv(hash)

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

	var forks []util.ForkMessage = sql.SQLSelectRecentForks(n)

	var forks_json []byte
	forks_json, err = json.Marshal(forks)
	if err != nil {
		fmt.Fprintf(w, "ERROR: invalid forks")
		return
	}

	fmt.Fprintf(w, string(forks_json))
}

func HandleGetRecentForksChart(w http.ResponseWriter, r *http.Request) {
	var err error

	// parse query
	var mQueryValues map[string][]string = r.URL.Query()

	var max_time uint64 = 0
	if len(mQueryValues["max_time"]) != 0 {
		max_time, err = strconv.ParseUint(mQueryValues["max_time"][0], 10, 64)
		if err != nil {
			fmt.Fprintf(w, "ERROR: invalid 'max_time' argument")
			return
		}
	}

	var min_time uint64 = 0
	if len(mQueryValues["min_time"]) != 0 {
		min_time, err = strconv.ParseUint(mQueryValues["min_time"][0], 10, 64)
		if err != nil {
			fmt.Fprintf(w, "ERROR: invalid 'min_time' argument")
			return
		}
	}

	var forks []util.ForkMessage = sql.SQLSelectRecentForksChart(min_time, max_time)

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

	var forks []util.ForkMessage = sql.SQLSelectRangeForks(min_height, max_height)

	var forks_json []byte
	forks_json, err = json.Marshal(forks)
	if err != nil {
		fmt.Fprintf(w, "ERROR: invalid forks")
	}

	fmt.Fprintf(w, string(forks_json))
}

func StartJSONServer() {

	// TODO flags
	JSONServer = http.Server{Addr: ":8080"}

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Println(filepath.Join(path + "/static/"))

	fileserver := http.FileServer(http.Dir(filepath.Join(path + "/static")))

	http.Handle("/static/", http.StripPrefix("/static", fileserver))

	//-- Pages --//
	http.HandleFunc("/", HandleGetFileIndex)
	http.HandleFunc("/block", HandleGetFileBlock)
	http.HandleFunc("/fork", HandleGetFileFork)

	//-- JSON --//
	// Blocks
	http.HandleFunc("/recentblocks", HandleRecentBlocks)
	http.HandleFunc("/rangeblocks", HandleGetRangeBlocks)
	http.HandleFunc("/hashblocks", HandleGetHashBlocks)

	// Inv
	http.HandleFunc("/recentinv", HandleGetRecentInv)
	http.HandleFunc("/hashinv", HandleGetHashInv)

	// Forks
	http.HandleFunc("/recentforks", HandleGetRecentForks)
	http.HandleFunc("/rangeforks", HandleGetRangeForks)
	http.HandleFunc("/v1/rest/recentforkschart", HandleGetRecentForksChart)

	if err := JSONServer.ListenAndServe(); err != nil {
		fmt.Printf("JSON server Shutdown: %s\n", err.Error())
	}
}

func StopJSONServer() {

	if err := JSONServer.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		fmt.Printf("HTTP server Shutdown: %s\n", err.Error())
	}
}
