package sql

import (
	_ "github.com/go-sql-driver/mysql"
	"log"

	"zcash-obs-db/util"
)

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

func SQLSelectBlockByHash(hash string) util.BlockMessage {
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
		return util.BlockMessage{}
	}

	block := util.BlockMessage{}
	for ret.Next() {
		ret.Scan(
			&block.Height,
			&block.Hash,
			&block.Prev_hash,
			&block.Coinbase_tx,
			&block.Num_tx,
			&block.Difficulty,
			&block.Block_size,
			&block.Miner_time,
			&block.Network_time,
		)
	}

	return block
}

func SQLSelectCurrentHeight() int {
	DBMutex.Lock()
	ret, err := DBConnection.Query(
		`SELECT height
		FROM blocks
		ORDER BY blocks.height DESC
		LIMIT 1;`,
	)
	DBMutex.Unlock()
	if err != nil {
		log.Printf("SQL Statement Prepare Error: %s\n", err.Error())
		return -1
	}

	height := 0
	for ret.Next() {
		ret.Scan(&height)
	}

	return height
}

func SQLSelectBlocksByHeightRange(min_height int, max_height int) []util.BlockMessage {
	DBMutex.Lock()
	ret, err := DBConnection.Query(
		`SELECT *
		FROM blocks
		WHERE height BETWEEN ? AND ?;`,
		min_height, max_height,
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
	}{}
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
		)

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
	}

	return blocks
}

func SQLSelectRecentBlocks(n int) []util.BlockMessage {
	DBMutex.Lock()
	ret, err := DBConnection.Query(
		`SELECT *
		FROM blocks
		ORDER BY blocks.height DESC
		LIMIT ?;`,
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
	}{}
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
		)

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
	}

	return blocks
}

func SQLSelectRecentBlocksWithInv(n int) []util.BlockMessage {
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
