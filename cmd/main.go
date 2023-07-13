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

	wssclient := websocket.InitWsClient()

	fmt.Println("network id: ", func() string {
		network, err := wssclient.NetworkID(context.Background())
		if err != nil {
			panic(err)
		}
		return network.String()
	}())

	go func() {
		contractAddress := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))
		events.RunSubscription(wssclient, contractAddress, "./pkg/events/abi.json")
		fmt.Println("Subscription started")
	}()

	select {}

}
