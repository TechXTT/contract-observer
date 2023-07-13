package websocket

import (
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/joho/godotenv/autoload"
)

func InitWsClient() *ethclient.Client {
	ws_url := os.Getenv("WS_URL")
	client, err := ethclient.Dial(ws_url)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
