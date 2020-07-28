package main

import (
	"zcash-obs-db/rest"
	"zcash-obs-db/shutdown"
	"zcash-obs-db/sql"
	"zcash-obs-db/util"
	"zcash-obs-db/websockets"

	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {

	// Define shutdown action
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			// sig is a ^C, handle it
			log.Printf("KILL")

			shutdown.RequestShutdown()
		}
	}()

	// Open Database Connection
	sql.OpenDBConnection()
	log.Printf("DB connection open")

	// Read obs.json
	ObsData, err := ioutil.ReadFile("./obs.json")
	if err != nil {
		log.Fatal("FATAL ERROR: Reading obs.json - %s", err.Error())
	}

	// Parse JSON
	var vJSON []map[string]string
	err = json.Unmarshal(ObsData, &vJSON)
	if err != nil {
		log.Fatal("FATAL ERROR: Parsing ObsData - %s", err.Error())
	}

	// Start Websocket connections
	vSessions := make([]util.OBSSession, len(vJSON))
	for i := 0; i < len(vJSON); i++ {
		vSessions[i].Name = vJSON[i]["name"]
		vSessions[i].Url = vJSON[i]["url"]
		vSessions[i].Origin = vJSON[i]["origin"]

		go websockets.HandleWebsocket(vSessions[i])
	}

	// Start JSON Server
	go rest.StartJSONServer()

	for {
		if shutdown.GetShutdownStatus() {

			// Stop JSON server
			rest.StopJSONServer()

			var countdown int = 3

			log.Printf("Shutting Down in %d seconds...", countdown)
			for ; countdown > 0; countdown-- {
				time.Sleep(1 * time.Second)
				log.Printf("%d", countdown)
			}
			break
		}
	}
}
