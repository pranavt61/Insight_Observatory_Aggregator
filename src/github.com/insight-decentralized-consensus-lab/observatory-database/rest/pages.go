package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Index HTML File
func HandleFileIndex(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadFile("./static/pages/index/page.html")
	if err != nil {
		fmt.Fprintf(w, "ERROR: reading index")
		return
	}

	fmt.Fprintf(w, string(body))
}

// Block HTML File
func HandleFileBlock(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadFile("./static/pages/block/page.html")
	if err != nil {
		fmt.Fprintf(w, "ERROR: reading block")
		return
	}

	fmt.Fprintf(w, string(body))
}

// Fork HTML File
func HandleFileFork(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadFile("./static/pages/fork/page.html")
	if err != nil {
		fmt.Fprintf(w, "ERROR: reading block")
		return
	}

	fmt.Fprintf(w, string(body))
}
