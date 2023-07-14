package main

import (
	"context"
	"fmt"
	"os"

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

	contractAddress := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))
	events.RunSubscription(wsclient, contractAddress, "./pkg/events/abi.json")
}
