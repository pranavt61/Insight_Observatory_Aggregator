package main

type BlockMessage struct {
	height       uint
	hash         string
	prev_hash    string
	coinbase_tx  string
	num_tx       uint
	difficulty   float64
	block_size   uint
	miner_time   uint64
	network_time uint64
}

type InvMessage struct {
	hash         string
	peer_ip      string
	network_time uint64
	session_id   int64
}

type PeerConnMessage struct {
	peer_ip      string
	version      uint
	subversion   string
	start_height uint
	services     uint64
	peer_time    uint64
	network_time uint64
	session_id   int64
}

type PeerDisMessage struct {
	peer_ip      string
	network_time uint64
	session_id   int64
}
