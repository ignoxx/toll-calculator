package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ignoxx/toll-calculator/types"
)

const (
	wsAddr       = "ws://127.0.0.1:30000/ws"
	sendInterval = time.Second
)

func main() {
	c, _, err := websocket.DefaultDialer.Dial(wsAddr, nil)
	if err != nil {
		log.Fatal("failed to connect to WS:", err)
	}
	defer c.Close()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	t := time.NewTicker(sendInterval)
	for {
		select {
		case <-signalChan:
			log.Println("interrupt signal received, shutting down client")
			return
		case <-t.C:
			obu := genObuData()
			_ = c.WriteJSON(obu)
		}
	}
}

func genObuData() types.ObuData {
	lat, long := genGeoData()

	return types.ObuData{
		ObuID: rand.Intn(999999999),
		Lat:   lat,
		Long:  long,
	}
}

/*
- Latitude: pick between -90 and 90
- Longitude: pick between -180 and 180
*/
func genGeoData() (float64, float64) {
	lat := rand.Float64()*180 - 90
	long := rand.Float64()*360 - 180
	return lat, long
}
