package sql

import (
	"log"

	"zcash-obs-db/util"
)

func SQLInsertPeerConn(peer util.PeerConnMessage) {

	DBMutex.Lock()
	stmt, err := DBConnection.Prepare(
		`INSERT INTO peer_conn
			(
				peer_ip,
				version,
				subversion,
				start_height,
				services,
				peer_time,
				network_time,
				session_id
			) VALUES(?,?,?,?,?,?,?,?);`,
	)
	DBMutex.Unlock()
	if err != nil {
		log.Printf("SQL Statement Prepare Error: %s\n", err.Error())
		return
	}

	_, err = stmt.Exec(
		peer.Peer_ip,
		peer.Version,
		peer.Subversion,
		peer.Start_height,
		peer.Services,
		peer.Peer_time,
		peer.Network_time,
		peer.Session_id,
	)
	if err != nil {
		log.Printf("SQL Statement Exec Error: %s\n", err.Error())
		return
	}

	stmt.Close()
	return
}

func SQLInsertPeerDis(peer util.PeerDisMessage) {

	DBMutex.Lock()
	stmt, err := DBConnection.Prepare(
		`UPDATE peer_conn set disconnect_time = ? where session_id = ? AND peer_ip = ? AND disconnect_time is null;`,
	)
	DBMutex.Unlock()
	if err != nil {
		log.Printf("SQL Statement Prepare Error: %s\n", err.Error())
		return
	}

	_, err = stmt.Exec(
		peer.Network_time,
		peer.Session_id,
		peer.Peer_ip,
	)
	if err != nil {
		log.Printf("SQL Statement Exec Error: %s\n", err.Error())
		return
	}

	stmt.Close()
	return
}
