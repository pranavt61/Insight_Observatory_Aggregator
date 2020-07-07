package main

import (
	"fmt"
	"net/http"
)

func HandleGetRecentBlocks(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "BLOCKS\n")
}

func StartJSONService() {
	http.HandleFunc("/recentblocks", HandleGetRecentBlocks)

	http.ListenAndServe(":8080", nil)
}
