package main

type BlockMessage struct {
	Height       uint    `json:"height"`
	Hash         string  `json:"hash"`
	Prev_hash    string  `json:"prev_hash"`
	Coinbase_tx  string  `json:"coinbase_tx"`
	Num_tx       uint    `json:"num_tx"`
	Difficulty   float64 `json:"difficulty"`
	Block_size   uint    `json:"block_size"`
	Miner_time   uint64  `json:"miner_time"`
	Network_time uint64  `json:"network_time"`
}

type InvMessage struct {
	Hash         string `json:"hash"`
	Peer_ip      string `json:"peer_ip"`
	Network_time uint64 `json:"network_time"`
	Session_id   int64  `json:"session_id"`
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

type ForkMessage struct {
	Height     uint           `json:"height"`
	Num_blocks uint           `json:"num_blocks"`
	Blocks     []BlockMessage `json:"blocks"`
}
