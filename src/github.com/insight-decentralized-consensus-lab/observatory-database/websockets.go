package main

import (
	"bytes"
	"encoding/json"
	"golang.org/x/net/websocket"
	"log"
	"math"
	"time"
)

// Helper
func clearBuffer(buf []byte, n int) {
	for i := 0; i < n; i++ {
		buf[i] = 0
	}
}

// Helper
// Hangs until connection established
func connectWebsocket(session OBSSession) (*websocket.Conn, int64) {
	url := session.url
	origin := session.origin

	nConnectionAttempt := 0

	for {
		log.Printf("Connection Websocket to %s", origin)

		// Connection attemps start at 1 second, and doubles every try
		// up to 30 seconds
		time.Sleep(time.Duration(math.Min(math.Exp2(float64(nConnectionAttempt)), 30)) * time.Second)

		ws, err := websocket.Dial(url, "", origin)

		if err == nil {
			log.Printf("Websocket open with %s", origin)

			// SQL entry
			session_id := SQLInsertSession(session)

			return ws, session_id
		}

		log.Printf("Failed to connect %s Websocket...Retrying", origin)
		nConnectionAttempt++
	}
}

// Run as Go Routine
func HandleWebsocket(session OBSSession) {

	// initial connection
	ws, session_id := connectWebsocket(session)

	// Read websocket messages
	var wsBuffer = make([]byte, 512)   // store websocket data
	var JSONBuffer = make([]byte, 512) // store indv. json objs
	var json_i = 0                     // read buffer index
	var unfinishedBrackets = 0         // track "{...}"
	for {
		// Check for Shutdown
		if GetShutdownStatus() {
			SQLUpdateDisconnectSession(session_id)

			break
		}

		nBytesRead, err := ws.Read(wsBuffer)
		if err != nil {
			SQLUpdateDisconnectSession(session_id)

			// reconnect
			ws, session_id = connectWebsocket(session)
		}

		for wsb_i := 0; wsb_i < nBytesRead; wsb_i++ {
			var c = wsBuffer[wsb_i]

			// start and end characters
			if c == '[' || c == ']' ||
				(unfinishedBrackets == 0 && c == ',') {
				continue
			}

			if c == '{' {
				unfinishedBrackets += 1
			}
			if c == '}' {
				unfinishedBrackets -= 1
			}

			JSONBuffer[json_i] = c
			json_i += 1

			if unfinishedBrackets == 0 {
				log.Printf("%s BUFFER: %s\n", session.name, string(JSONBuffer))

				var result map[string]interface{}
				err = json.Unmarshal(bytes.Trim(JSONBuffer, "\x00"), &result)
				if err != nil {
					log.Printf("JSON Parsing Error: %s\n", err.Error())

					clearBuffer(JSONBuffer, 512)
					json_i = 0
					continue
				}

				if result["type"] == "block" {
					msg := BlockMessage{
						uint(result["data"].(map[string]interface{})["height"].(float64)),
						result["data"].(map[string]interface{})["hash"].(string),
						result["data"].(map[string]interface{})["prev_hash"].(string),
						result["data"].(map[string]interface{})["coinbase_tx"].(string),
						uint(result["data"].(map[string]interface{})["num_tx"].(float64)),
						result["data"].(map[string]interface{})["difficulty"].(float64),
						uint(result["data"].(map[string]interface{})["block_size"].(float64)),
						uint64(result["data"].(map[string]interface{})["miner_time"].(float64) * 1000),
						uint64(result["data"].(map[string]interface{})["network_time"].(float64) * 1000),
					}

					SQLInsertBlock(msg)
				} else if result["type"] == "inv" {
					msg := InvMessage{
						result["data"].(map[string]interface{})["hash"].(string),
						result["data"].(map[string]interface{})["peer_ip"].(string),
						uint64(result["data"].(map[string]interface{})["network_time"].(float64) * 1000),
						session_id,
					}

					SQLInsertInv(msg)
				} else if result["type"] == "peer_conn" {
					msg := PeerConnMessage{
						result["data"].(map[string]interface{})["peer_ip"].(string),
						uint(result["data"].(map[string]interface{})["version"].(float64)),
						result["data"].(map[string]interface{})["subversion"].(string),
						uint(result["data"].(map[string]interface{})["start_height"].(float64)),
						uint64(result["data"].(map[string]interface{})["services"].(float64)),
						uint64(result["data"].(map[string]interface{})["peer_time"].(float64) * 1000),
						uint64(result["data"].(map[string]interface{})["network_time"].(float64) * 1000),
						session_id,
					}

					SQLInsertPeerConn(msg)
				} else if result["type"] == "peer_dis" {
					msg := PeerDisMessage{
						result["data"].(map[string]interface{})["peer_ip"].(string),
						uint64(result["data"].(map[string]interface{})["network_time"].(float64) * 1000),
						session_id,
					}

					SQLInsertPeerDis(msg)
				}

				clearBuffer(JSONBuffer, 512)
				json_i = 0
			}
		}
	}
}
