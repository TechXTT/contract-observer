package events

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type VoteData struct {
	TxHash      common.Hash
	Executor    common.Address
	Amount      *big.Int
	AssetID     uint16
	SourceChain *big.Int
}

func SubscribeToEvents(c *ethclient.Client, contractAddress common.Address, logs chan<- types.Log, contractABI abi.ABI) (ethereum.Subscription, error) {

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{contractABI.Events["Lock"].ID}},
	}

	subscription, err := c.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func DecodeEvent(client *ethclient.Client, contractABI abi.ABI, vLog types.Log) (VoteData, error) {
	event := contractABI.Events["Lock"]

	lockEventMap := make(map[string]interface{})
	err := contractABI.UnpackIntoMap(lockEventMap, event.Name, vLog.Data)
	if err != nil {
		return VoteData{}, err
	}

	amount, ok := lockEventMap["amount"].(*big.Int)
	if !ok {
		return VoteData{}, err
	}

	txHash := vLog.TxHash

	assetID64 := vLog.Topics[1].Big().Uint64()
	assetID := uint16(assetID64)

	sourceChain, err := client.NetworkID(context.Background())
	if err != nil {
		return VoteData{}, err
	}

	user := common.HexToAddress(vLog.Topics[2].Hex())

	voteData := VoteData{
		TxHash:      txHash,
		Executor:    user,
		Amount:      amount,
		AssetID:     assetID,
		SourceChain: sourceChain,
	}

	return voteData, nil

}

func RunSubscription(client *ethclient.Client, contractAddress common.Address, contractAbi string) {
	logs := make(chan types.Log)

	fileBytes, err := os.ReadFile("./pkg/events/abi.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fileContents := string(fileBytes)

	contractABI, err := abi.JSON(strings.NewReader(fileContents))
	if err != nil {
		log.Fatal(err)
	}

	subscription, err := SubscribeToEvents(client, contractAddress, logs, contractABI)
	if err != nil {
		log.Fatal(err)
	}
	defer subscription.Unsubscribe()

	boundContract := bind.NewBoundContract(contractAddress, contractABI, client, client, client)

	callOpts := bind.CallOpts{
		Pending: true,
		From:    common.HexToAddress(os.Getenv("OBSERVER_ADDRESS")),
	}

	fmt.Println("Observer address: ", callOpts.From.Hex())

	fmt.Println("Subscription started:")

	for {
		select {
		case err := <-subscription.Err():
			log.Fatal(err)
		case vLog := <-logs:
			data, err := DecodeEvent(client, contractABI, vLog)
			if err != nil {
				log.Fatal(err)
				continue
			}

			fmt.Println("----------------------------------------")

			fmt.Println("txHash: ", data.TxHash.Hex())
			fmt.Println("executor: ", data.Executor.Hex())
			fmt.Println("amount:", data.Amount)
			fmt.Println("assetID: ", data.AssetID)
			fmt.Println("sourceChain: ", data.SourceChain)

			err = boundContract.Call(&callOpts, nil, "vote", data.TxHash, data.Executor, data.Amount, data.AssetID, data.SourceChain)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Vote successful!")
		}
	}
}
