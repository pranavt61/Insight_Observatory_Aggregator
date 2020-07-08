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

func SQLInsertBlock(msg BlockMessage) {

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
		msg.Height,
		msg.Hash,
		msg.Prev_hash,
		msg.Coinbase_tx,
		msg.Num_tx,
		msg.Difficulty,
		msg.Block_size,
		msg.Miner_time,
		msg.Network_time,
	)
	if err != nil {
		log.Printf("SQL Statement Exec Error: %s\n", err.Error())
		return
	}

	stmt.Close()

	return
}

func SQLInsertInv(msg InvMessage) {

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
		msg.hash,
		msg.peer_ip,
		msg.network_time,
		msg.session_id,
	)
	if err != nil {
		log.Printf("SQL Statement Exec Error: %s\n", err.Error())
		return
	}
	stmt.Close()
	return
}

func SQLInsertPeerConn(msg PeerConnMessage) {

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
		msg.peer_ip,
		msg.version,
		msg.subversion,
		msg.start_height,
		msg.services,
		msg.peer_time,
		msg.network_time,
		msg.session_id,
	)
	if err != nil {
		log.Printf("SQL Statement Exec Error: %s\n", err.Error())
		return
	}

	stmt.Close()
	return
}

func SQLInsertPeerDis(msg PeerDisMessage) {

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
		msg.network_time,
		msg.session_id,
		msg.peer_ip,
	)
	if err != nil {
		log.Printf("SQL Statement Exec Error: %s\n", err.Error())
		return
	}

	stmt.Close()
	return
}

//--- Network SELECT Commands --//

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

func SQLSelectRecentForks(n int) []ForkMessage {
	DBMutex.Lock()
	ret, err := DBConnection.Query(
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

func SQLSelectInv() {

}
