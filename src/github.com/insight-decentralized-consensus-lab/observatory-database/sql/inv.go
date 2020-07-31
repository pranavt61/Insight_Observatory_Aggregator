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

func SQLSelectAllInvHalfRange() []uint64 {
	DBMutex.Lock()
	ret, err := DBConnection.Query(
		`SELECT
			MIN(network_time) as smallest,
			AVG(network_time) as average
		FROM inv 
		GROUP BY inv.hash;`,
	)
	DBMutex.Unlock()
	if err != nil {
		log.Printf("SQL Statement Prepare Error: %s\n", err.Error())
		return nil
	}

	times := make([]uint64, 0)
	var smallest_buffer uint64 = 0
	var average_buffer float64 = 0.0
	for ret.Next() {
		ret.Scan(
			&smallest_buffer,
			&average_buffer,
		)

		times = append(times, (uint64(average_buffer) - smallest_buffer))
	}

	return times
}
