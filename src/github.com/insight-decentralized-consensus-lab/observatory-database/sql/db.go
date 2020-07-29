package sql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

var DBConnection *sql.DB
var DBMutex sync.Mutex

func OpenDBConnection() {
	var err error

	DBConnection, err = sql.Open("mysql", "OBS_USER:pass@tcp(localhost:3306)/OBS_Cluster")
	if err != nil {
		panic(err)
	}
}

func CloseDBConnection() {
	DBMutex.Lock()
	DBConnection.Close()
	DBMutex.Unlock()
}
