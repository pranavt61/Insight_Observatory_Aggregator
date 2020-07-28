package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"zcash-obs-db/sql"
	"zcash-obs-db/util"
)

func HandleRecentForks(w http.ResponseWriter, r *http.Request) {
	var forks []util.ForkMessage = sql.SQLSelectRecentForks(5)

	raw_json, err := json.Marshal(forks)
	if err != nil {
		fmt.Fprintf(w, "500")
		return
	}

	fmt.Fprintf(w, string(raw_json))
}
