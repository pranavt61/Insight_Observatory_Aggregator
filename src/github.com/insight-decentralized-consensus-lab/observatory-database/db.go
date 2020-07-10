package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
	"time"
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

//-- Websocket Entries --//
func SQLInsertSession(session OBSSession) int64 {
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
		session.url,
		session.name,
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

//--- Network INSERT Commands --//

func SQLInsertBlock(block BlockMessage) {

	DBMutex.Lock()
	stmt, err := DBConnection.Prepare(
		`INSERT INTO blocks
			(
				height,
				hash,
				prev_hash,
				coinbase_tx,
				num_tx,
				difficulty,
				block_size,
				miner_time,
				network_time
			) VALUES(?,?,?,?,?,?,?,?,?);`,
	)
	DBMutex.Unlock()
	if err != nil {
		log.Printf("SQL Statement Prepare Error: %s\n", err.Error())
		return
	}

	_, err = stmt.Exec(
		block.Height,
		block.Hash,
		block.Prev_hash,
		block.Coinbase_tx,
		block.Num_tx,
		block.Difficulty,
		block.Block_size,
		block.Miner_time,
		block.Network_time,
	)
	if err != nil {
		log.Printf("SQL Statement Exec Error: %s\n", err.Error())
		return
	}

	stmt.Close()

	return
}

func SQLInsertInv(inv InvMessage) {

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

func SQLInsertPeerConn(peer PeerConnMessage) {

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
		peer.peer_ip,
		peer.version,
		peer.subversion,
		peer.start_height,
		peer.services,
		peer.peer_time,
		peer.network_time,
		peer.session_id,
	)
	if err != nil {
		log.Printf("SQL Statement Exec Error: %s\n", err.Error())
		return
	}

	stmt.Close()
	return
}

func SQLInsertPeerDis(peer PeerDisMessage) {

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
		peer.network_time,
		peer.session_id,
		peer.peer_ip,
	)
	if err != nil {
		log.Printf("SQL Statement Exec Error: %s\n", err.Error())
		return
	}

	stmt.Close()
	return
}

//--- JSON Server SELECT Commands --//

func SQLSelectRecentBlocks(n int) []BlockMessage {
	DBMutex.Lock()
	ret, err := DBConnection.Query(
		`SELECT * from blocks
		ORDER BY blocks.height DESC
		LIMIT ?;`,
		n,
	)
	DBMutex.Unlock()
	if err != nil {
		log.Printf("SQL Statement Prepare Error: %s\n", err.Error())
		return nil
	}

	var blocks []BlockMessage = make([]BlockMessage, 0)
	var i int = 0
	for ret.Next() {
		blocks = append(blocks, BlockMessage{})

		ret.Scan(
			&blocks[i].Height,
			&blocks[i].Hash,
			&blocks[i].Prev_hash,
			&blocks[i].Coinbase_tx,
			&blocks[i].Num_tx,
			&blocks[i].Difficulty,
			&blocks[i].Block_size,
			&blocks[i].Miner_time,
			&blocks[i].Network_time,
		)

		i++
	}

	return blocks
}

func SQLSelectRangeBlocks(min int, max int) []BlockMessage {
	DBMutex.Lock()
	ret, err := DBConnection.Query(
		`SELECT * from blocks
		WHERE
			height > ?
			AND
			height < ?;`,
		min-1, max,
	)
	DBMutex.Unlock()
	if err != nil {
		log.Printf("SQL Statement Prepare Error: %s\n", err.Error())
		return nil
	}

	var blocks []BlockMessage = make([]BlockMessage, 0)
	var i int = 0
	for ret.Next() {
		blocks = append(blocks, BlockMessage{})

		ret.Scan(
			&blocks[i].Height,
			&blocks[i].Hash,
			&blocks[i].Prev_hash,
			&blocks[i].Coinbase_tx,
			&blocks[i].Num_tx,
			&blocks[i].Difficulty,
			&blocks[i].Block_size,
			&blocks[i].Miner_time,
			&blocks[i].Network_time,
		)

		i++
	}

	return blocks
}

func SQLSelectHashBlocks(hash string) []BlockMessage {

	// count Rows
	DBMutex.Lock()
	ret, err := DBConnection.Query(
		`SELECT
			*
		FROM blocks
		WHERE hash = ?;`,
		hash,
	)
	DBMutex.Unlock()
	if err != nil {
		log.Printf("SQL Statement Prepare Error: %s\n", err.Error())
		return nil
	}

	var blocks []BlockMessage = make([]BlockMessage, 0)
	var i int = 0
	for ret.Next() {
		blocks = append(blocks, BlockMessage{})

		ret.Scan(
			&blocks[i].Height,
			&blocks[i].Hash,
			&blocks[i].Prev_hash,
			&blocks[i].Coinbase_tx,
			&blocks[i].Num_tx,
			&blocks[i].Difficulty,
			&blocks[i].Block_size,
			&blocks[i].Miner_time,
			&blocks[i].Network_time,
		)

		i++
	}

	return blocks
}

func SQLSelectRecentInv(n int) []InvMessage {
	DBMutex.Lock()
	ret, err := DBConnection.Query(
		`SELECT * from inv
		ORDER BY inv.network_time DESC
		LIMIT ?;`,
		n,
	)
	DBMutex.Unlock()
	if err != nil {
		log.Printf("SQL Statement Prepare Error: %s\n", err.Error())
		return nil
	}

	var inv []InvMessage = make([]InvMessage, 0)
	var i int = 0
	for ret.Next() {
		inv = append(inv, InvMessage{})

		ret.Scan(
			&inv[i].Hash,
			&inv[i].Peer_ip,
			&inv[i].Network_time,
			&inv[i].Session_id,
		)

		i++
	}

	return inv
}

func SQLSelectHashInv(hash string) []InvMessage {
	DBMutex.Lock()
	ret, err := DBConnection.Query(
		`SELECT * from inv
		WHERE hash = ?;`,
		hash,
	)
	DBMutex.Unlock()
	if err != nil {
		log.Printf("SQL Statement Prepare Error: %s\n", err.Error())
		return nil
	}

	var inv []InvMessage = make([]InvMessage, 0)
	var i int = 0
	for ret.Next() {
		inv = append(inv, InvMessage{})

		ret.Scan(
			&inv[i].Hash,
			&inv[i].Peer_ip,
			&inv[i].Network_time,
			&inv[i].Session_id,
		)

		i++
	}

	return inv
}

func SQLSelectRecentForks(n int) []ForkMessage {
	DBMutex.Lock()
	forks_ret, err := DBConnection.Query(
		`SELECT
			height,
			COUNT(height) n
		FROM blocks
		GROUP BY height
		HAVING n > 1
		LIMIT ?;`,
		n,
	)
	DBMutex.Unlock()
	if err != nil {
		log.Printf("SQL Statement Prepare Error: %s\n", err.Error())
		return nil
	}

	var forks []ForkMessage = make([]ForkMessage, 0)
	var i int = 0
	for forks_ret.Next() {
		forks = append(forks, ForkMessage{})

		forks_ret.Scan(
			&forks[i].Height,
			&forks[i].Num_blocks,
		)

		// Query Block data
		forks[i].Blocks = make([]BlockMessage, forks[i].Num_blocks)

		DBMutex.Lock()
		block_ret, err := DBConnection.Query(
			`SELECT
				*
			FROM blocks 
			WHERE height = ?;`,
			forks[i].Height,
		)
		DBMutex.Unlock()
		if err != nil {
			log.Printf("SQL Statement Prepare Error: %s\n", err.Error())
			return nil
		}

		var j int = 0
		for block_ret.Next() {
			block_ret.Scan(
				&forks[i].Blocks[j].Height,
				&forks[i].Blocks[j].Hash,
				&forks[i].Blocks[j].Prev_hash,
				&forks[i].Blocks[j].Coinbase_tx,
				&forks[i].Blocks[j].Num_tx,
				&forks[i].Blocks[j].Difficulty,
				&forks[i].Blocks[j].Block_size,
				&forks[i].Blocks[j].Miner_time,
				&forks[i].Blocks[j].Network_time,
			)

			j++
		}

		i++
	}

	return forks
}

func SQLSelectRangeForks(min int, max int) []ForkMessage {

	// count Rows
	DBMutex.Lock()
	ret, err := DBConnection.Query(
		`SELECT
			height,
			COUNT(height) n
		FROM blocks
		GROUP BY height
		HAVING n > 1 AND height > ? AND height < ?;`,
		min-1, max,
	)
	DBMutex.Unlock()
	if err != nil {
		log.Printf("SQL Statement Prepare Error: %s\n", err.Error())
		return nil
	}

	var forks []ForkMessage = make([]ForkMessage, 0)
	var i int = 0
	for ret.Next() {
		forks = append(forks, ForkMessage{})

		ret.Scan(
			&forks[i].Height,
			&forks[i].Num_blocks,
		)

		i++
	}

	return forks
}
