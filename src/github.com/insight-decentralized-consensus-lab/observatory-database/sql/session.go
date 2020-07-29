package sql

import (
	"log"
	"time"

	"zcash-obs-db/util"
)

func SQLInsertSession(session util.OBSSession) int64 {
	// SQL entry
	DBMutex.Lock()
	stmt, err := DBConnection.Prepare(
		`INSERT INTO obs_sessions
		(
			ip,
			name,
			start_time
		) VALUES(?,?,?);`,
	)
	DBMutex.Unlock()
	if err != nil {
		log.Fatal("SQL Statement Prepare Error: %s\n", err.Error())
	}

	sql_res, err := stmt.Exec(
		session.Url,
		session.Name,
		time.Now().UnixNano()/int64(time.Millisecond),
	)
	if err != nil {
		log.Fatal("SQL Statement Exec Error: %s\n", err.Error())
	}
	stmt.Close()

	session_id, err := sql_res.LastInsertId()
	if err != nil {
		log.Fatal("SQL LastInsertID Error: %s\n", err.Error())
	}

	return session_id
}

func SQLUpdateDisconnectSession(session_id int64) {
	log.Printf("DISCONNECT SQL")
	// SQL entry
	DBMutex.Lock()
	stmt, err := DBConnection.Prepare(
		"UPDATE obs_sessions SET end_time=? WHERE session_id=?;",
	)
	DBMutex.Unlock()
	if err != nil {
		log.Fatal("SQL Statement Prepare Error: %s\n", err.Error())
	}

	_, err = stmt.Exec(
		time.Now().UnixNano()/int64(time.Millisecond),
		session_id,
	)
	if err != nil {
		log.Printf("SQL Statement Exec Error: %s\n", err.Error())
	}
	stmt.Close()
}
