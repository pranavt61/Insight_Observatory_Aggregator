package util

type BlockMessage struct {
	Height       uint         `json:"height"`
	Hash         string       `json:"hash"`
	Prev_hash    string       `json:"prev_hash"`
	Coinbase_tx  string       `json:"coinbase_tx"`
	Num_tx       uint         `json:"num_tx"`
	Difficulty   float64      `json:"difficulty"`
	Block_size   uint         `json:"block_size"`
	Miner_time   uint64       `json:"miner_time"`
	Network_time uint64       `json:"network_time"`
	Inv          []InvMessage `json:"inv"`
}

type InvMessage struct {
	Hash         string `json:"hash"`
	Peer_ip      string `json:"peer_ip"`
	Network_time uint64 `json:"network_time"`
	Session_id   int64  `json:"session_id"`
}

type PeerConnMessage struct {
	Peer_ip      string
	Version      uint
	Subversion   string
	Start_height uint
	Services     uint64
	Peer_time    uint64
	Network_time uint64
	Session_id   int64
}

type PeerDisMessage struct {
	Peer_ip      string
	Network_time uint64
	Session_id   int64
}

type ForkMessage struct {
	Height     uint           `json:"height"`
	Num_blocks uint           `json:"num_blocks"`
	Blocks     []BlockMessage `json:"blocks"`
}
