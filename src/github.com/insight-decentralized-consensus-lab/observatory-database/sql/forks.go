package sql

import (
	"log"

	"zcash-obs-db/util"
)

func SQLSelectCurrentFork() bool {
	DBMutex.Lock()
	ret, err := DBConnection.Query(
		`SELECT
			COUNT(height)
		FROM blocks
		WHERE height = (SELECT MAX(height) FROM blocks);`,
	)
	DBMutex.Unlock()
	if err != nil {
		log.Printf("SQL Statement Prepare Error: %s\n", err.Error())
		return false
	}

	num_blocks := 0
	for ret.Next() {
		ret.Scan(&num_blocks)
	}

	if num_blocks > 1 {
		return true
	}

	return false
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

func SQLSelectRangeForks(nMinHeight int, nMaxHeight int) []util.ForkMessage {
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
			WHERE height BETWEEN ? AND ?
			GROUP BY height
			HAVING n > 1) as b2
		ON b1.height = b2.height;`,
		nMinHeight, nMaxHeight,
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
