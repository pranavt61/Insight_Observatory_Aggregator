package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"zcash-obs-db/sql"
)

func HandleAllInvHalfRange(w http.ResponseWriter, r *http.Request) {
	inv_times := sql.SQLSelectAllInvHalfRange()

	raw_json, err := json.Marshal(inv_times)
	if err != nil {
		fmt.Fprintf(w, "500")
		return
	}

	fmt.Fprintf(w, string(raw_json))
}
