package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TechXTT/contract-observer/pkg/events"
	"github.com/TechXTT/contract-observer/pkg/websocket"
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	wsclient := websocket.InitWsClient()

	fmt.Println("network id: ", func() string {
		network, err := wsclient.NetworkID(context.Background())
		if err != nil {
			panic(err)
		}
		return network.String()
	}())

	// separate thread to send requests to the network to keep the connection alive
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-context.Background().Done():
				return
			case <-ticker.C:
				_, err := wsclient.BlockByNumber(context.Background(), nil)
				if err != nil {
					log.Println("Error sending keep-alive purposed request: ", err)
				}
			}
		}
	}()

	contractAddress := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))
	events.RunSubscription(wsclient, contractAddress, "./pkg/events/abi.json")

	for range time.Tick(time.Second) {
	}

}
