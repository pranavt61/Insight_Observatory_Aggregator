package sql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
	"time"
	"zcash-obs-db/util"
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

//--- Network INSERT Commands --//

func SQLInsertBlock(block util.BlockMessage) {

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

//--- JSON Server SELECT Commands --//
func SQLSelectRecentBlocks(n int) []util.BlockMessage {
	DBMutex.Lock()
	ret, err := DBConnection.Query(
		`SELECT
			*
		FROM (
			SELECT *
			FROM blocks
			ORDER BY blocks.height DESC
			LIMIT ?
		) AS recent_blocks
		LEFT JOIN inv ON recent_blocks.hash = inv.hash;`,
		n,
	)
	DBMutex.Unlock()
	if err != nil {
		log.Printf("SQL Statement Prepare Error: %s\n", err.Error())
		return nil
	}

	blocks := make([]util.BlockMessage, 0)
	block_buffer := struct {
		Blocks_height       uint
		Blocks_hash         string
		Blocks_prev_hash    string
		Blocks_coinbase_tx  string
		Blocks_num_tx       uint
		Blocks_difficulty   float64
		Blocks_block_size   uint
		Blocks_miner_time   uint64
		Blocks_network_time uint64
		Inv_hash            string
		Inv_peer_ip         string
		Inv_network_time    uint64
		Inv_session_id      int64
	}{}
	seen_blocks := make(map[string]int)
	i := 0
	for ret.Next() {

		// Parse SQL
		ret.Scan(
			&(block_buffer.Blocks_height),
			&(block_buffer.Blocks_hash),
			&(block_buffer.Blocks_prev_hash),
			&(block_buffer.Blocks_coinbase_tx),
			&(block_buffer.Blocks_num_tx),
			&(block_buffer.Blocks_difficulty),
			&(block_buffer.Blocks_block_size),
			&(block_buffer.Blocks_miner_time),
			&(block_buffer.Blocks_network_time),
			&(block_buffer.Inv_hash),
			&(block_buffer.Inv_peer_ip),
			&(block_buffer.Inv_network_time),
			&(block_buffer.Inv_session_id),
		)

		if index, ok := seen_blocks[block_buffer.Blocks_hash]; ok {
			// block seen
			inv := util.InvMessage{
				block_buffer.Inv_hash,
				block_buffer.Inv_peer_ip,
				block_buffer.Inv_network_time,
				block_buffer.Inv_session_id,
			}

			blocks[index].Inv = append(blocks[index].Inv, inv)
		} else {
			block := util.BlockMessage{
				block_buffer.Blocks_height,
				block_buffer.Blocks_hash,
				block_buffer.Blocks_prev_hash,
				block_buffer.Blocks_coinbase_tx,
				block_buffer.Blocks_num_tx,
				block_buffer.Blocks_difficulty,
				block_buffer.Blocks_block_size,
				block_buffer.Blocks_miner_time,
				block_buffer.Blocks_network_time,
				make([]util.InvMessage, 0),
			}
			blocks = append(blocks, block)

			inv := util.InvMessage{
				block_buffer.Inv_hash,
				block_buffer.Inv_peer_ip,
				block_buffer.Inv_network_time,
				block_buffer.Inv_session_id,
			}
			blocks[i].Inv = append(blocks[i].Inv, inv)

			seen_blocks[block.Hash] = i
			i++
		}

	}

	return blocks
}

func SQLSelectRecentBlocksTable(n int) []util.BlockMessage {
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

	var blocks []util.BlockMessage = make([]util.BlockMessage, 0)
	var i int = 0
	for ret.Next() {
		blocks = append(blocks, util.BlockMessage{})

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

		DBMutex.Lock()
		inv_ret, err := DBConnection.Query(
			`SELECT * from inv
			WHERE hash=?;`,
			blocks[i].Hash,
		)
		DBMutex.Unlock()
		if err != nil {
			log.Printf("SQL Statement Prepare Error: %s\n", err.Error())
			return nil
		}

		j := 0
		for inv_ret.Next() {
			blocks[i].Inv = append(blocks[i].Inv, util.InvMessage{})

			inv_ret.Scan(
				&blocks[i].Inv[j].Hash,
				&blocks[i].Inv[j].Peer_ip,
				&blocks[i].Inv[j].Network_time,
				&blocks[i].Inv[j].Session_id,
			)

			j++
		}

		i++
	}

	return blocks
}

func SQLSelectRangeBlocks(min int, max int) []util.BlockMessage {
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

	var blocks []util.BlockMessage = make([]util.BlockMessage, 0)
	var i int = 0
	for ret.Next() {
		blocks = append(blocks, util.BlockMessage{})

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

func SQLSelectHashBlocks(hash string) []util.BlockMessage {

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

	var blocks []util.BlockMessage = make([]util.BlockMessage, 0)
	var i int = 0
	for ret.Next() {
		blocks = append(blocks, util.BlockMessage{})

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

func SQLSelectRecentInv(n int) []util.InvMessage {
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

	var inv []util.InvMessage = make([]util.InvMessage, 0)
	var i int = 0
	for ret.Next() {
		inv = append(inv, util.InvMessage{})

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

func SQLSelectHashInv(hash string) []util.InvMessage {
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

	var inv []util.InvMessage = make([]util.InvMessage, 0)
	var i int = 0
	for ret.Next() {
		inv = append(inv, util.InvMessage{})

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

func SQLSelectRecentForks(n int) []util.ForkMessage {
	DBMutex.Lock()
	ret, err := DBConnection.Query(
		`SELECT
			*
		FROM blocks as b1
		INNER JOIN (
			SELECT
				height,
				COUNT(height) n
			FROM blocks
			GROUP BY height
			HAVING n > 1
			ORDER BY blocks.height DESC
			LIMIT ?) as b2
		ON b1.height = b2.height;`,
		n,
	)
	DBMutex.Unlock()
	if err != nil {
		log.Printf("SQL Statement Prepare Error: %s\n", err.Error())
		return nil
	}

	forks := make([]util.ForkMessage, 0)
	fork_buffer := struct {
		Block_height       uint
		Block_hash         string
		Block_prev_hash    string
		Block_coinbase_tx  string
		Block_num_tx       uint
		Block_difficulty   float64
		Block_size         uint
		Block_miner_time   uint64
		Block_network_time uint64
		Fork_height        uint
		Fork_size          int
	}{}
	seen_forks := make(map[uint]int)
	i := 0
	for ret.Next() {

		ret.Scan(
			&(fork_buffer.Block_height),
			&(fork_buffer.Block_hash),
			&(fork_buffer.Block_prev_hash),
			&(fork_buffer.Block_coinbase_tx),
			&(fork_buffer.Block_num_tx),
			&(fork_buffer.Block_difficulty),
			&(fork_buffer.Block_size),
			&(fork_buffer.Block_miner_time),
			&(fork_buffer.Block_network_time),
			&(fork_buffer.Fork_height),
			&(fork_buffer.Fork_size),
		)

		log.Println(fork_buffer)

		if index, ok := seen_forks[fork_buffer.Block_height]; ok {
			// fork seen
			block := util.BlockMessage{
				fork_buffer.Block_height,
				fork_buffer.Block_hash,
				fork_buffer.Block_prev_hash,
				fork_buffer.Block_coinbase_tx,
				fork_buffer.Block_num_tx,
				fork_buffer.Block_difficulty,
				fork_buffer.Block_size,
				fork_buffer.Block_miner_time,
				fork_buffer.Block_network_time,
				make([]util.InvMessage, 0),
			}

			forks[index].Blocks = append(forks[index].Blocks, block)
			forks[index].Num_blocks++
		} else {
			// new fork
			fork := util.ForkMessage{
				fork_buffer.Block_height,
				1,
				make([]util.BlockMessage, 0),
			}
			forks = append(forks, fork)

			block := util.BlockMessage{
				fork_buffer.Block_height,
				fork_buffer.Block_hash,
				fork_buffer.Block_prev_hash,
				fork_buffer.Block_coinbase_tx,
				fork_buffer.Block_num_tx,
				fork_buffer.Block_difficulty,
				fork_buffer.Block_size,
				fork_buffer.Block_miner_time,
				fork_buffer.Block_network_time,
				make([]util.InvMessage, 0),
			}

			forks[i].Blocks = append(forks[i].Blocks, block)

			seen_forks[fork.Height] = i

			i++
		}
	}

	return forks
}

func SQLSelectRecentForksChart(min_time uint64, max_time uint64) []util.ForkMessage {
	DBMutex.Lock()
	forks_ret, err := DBConnection.Query(
		`SELECT
			height,
			COUNT(height) n
		FROM blocks
		GROUP BY height
		HAVING n > 1;`,
	)
	DBMutex.Unlock()
	if err != nil {
		log.Printf("SQL Statement Prepare Error: %s\n", err.Error())
		return nil
	}

	var forks []util.ForkMessage = make([]util.ForkMessage, 0)
	var i int = 0
	for forks_ret.Next() {
		forks = append(forks, util.ForkMessage{})

		forks_ret.Scan(
			&forks[i].Height,
			&forks[i].Num_blocks,
		)

		// Query Block data
		forks[i].Blocks = make([]util.BlockMessage, forks[i].Num_blocks)

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

func SQLSelectRangeForks(min int, max int) []util.ForkMessage {

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

	var forks []util.ForkMessage = make([]util.ForkMessage, 0)
	var i int = 0
	for ret.Next() {
		forks = append(forks, util.ForkMessage{})

		ret.Scan(
			&forks[i].Height,
			&forks[i].Num_blocks,
		)

		i++
	}

	return forks
}
