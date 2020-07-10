package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"
)

type OBSSession struct {
	name   string
	url    string
	origin string
}

func main() {

	// Define shutdown action
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			// sig is a ^C, handle it
			log.Printf("KILL")

			RequestShutdown()
		}
	}()

	// Open Database Connection
	OpenDBConnection()
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
	vSessions := make([]OBSSession, len(vJSON))
	for i := 0; i < len(vJSON); i++ {
		vSessions[i].name = vJSON[i]["name"]
		vSessions[i].url = vJSON[i]["url"]
		vSessions[i].origin = vJSON[i]["origin"]

		go HandleWebsocket(vSessions[i])
	}

	// Start JSON Server
	StartJSONServer()

	for {
		if GetShutdownStatus() {

			// Stop JSON server
			StopJSONServer()

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
