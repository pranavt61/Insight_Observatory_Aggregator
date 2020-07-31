package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

var JSONServer http.Server

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
	http.HandleFunc("/", HandleFileIndex)
	http.HandleFunc("/block", HandleFileBlock)
	http.HandleFunc("/fork", HandleFileFork)

	//-- JSON --//
	// Blocks
	http.HandleFunc("/v1/json/currentheight", HandleCurrentHeight)
	http.HandleFunc("/v1/json/getblockbyhash", HandleGetBlockByHash)
	http.HandleFunc("/v1/json/getblocksbyheightrange", HandleGetBlocksByHeightRange)
	http.HandleFunc("/v1/json/recentblocks", HandleRecentBlocks)

	// Forks
	http.HandleFunc("/v1/json/currentfork", HandleCurrentFork)
	http.HandleFunc("/v1/json/recentforks", HandleRecentForks)
	http.HandleFunc("/v1/json/rangeforks", HandleRangeForks)

	// Inv
	http.HandleFunc("/v1/json/allinvhalfrange", HandleAllInvHalfRange)

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
