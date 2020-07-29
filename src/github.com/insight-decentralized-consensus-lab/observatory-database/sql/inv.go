package sql

import (
	"log"

	"zcash-obs-db/util"
)

func SQLInsertInv(inv util.InvMessage) {

	DBMutex.Lock()
	stmt, err := DBConnection.Prepare(
		`INSERT INTO inv
			(
				hash,
				peer_ip,
				network_time,
				session_id
			) VALUES(?,?,?,?);`,
	)
	DBMutex.Unlock()
	if err != nil {
		log.Printf("SQL Statement Prepare Error: %s\n", err.Error())
		return
	}

	_, err = stmt.Exec(
		inv.Hash,
		inv.Peer_ip,
		inv.Network_time,
		inv.Session_id,
	)
	if err != nil {
		log.Printf("SQL Statement Exec Error: %s\n", err.Error())
		return
	}
	stmt.Close()
	return
}
