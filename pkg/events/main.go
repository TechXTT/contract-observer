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

type LockEvent struct {
	Amount      *big.Int
	TargetChain *big.Int
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

func DecodeEvent(contractABI abi.ABI, vLog types.Log) (LockEvent, error) {
	event := contractABI.Events["Lock"]

	lockEventMap := make(map[string]interface{})
	err := contractABI.UnpackIntoMap(lockEventMap, event.Name, vLog.Data)
	if err != nil {
		return LockEvent{}, err
	}

	amount, ok := lockEventMap["amount"].(*big.Int)
	if !ok {
		return LockEvent{}, err
	}

	targetChain, ok := lockEventMap["targetChain"].(*big.Int)
	if !ok {
		return LockEvent{}, err
	}

	lockEvent := LockEvent{
		Amount:      amount,
		TargetChain: targetChain,
	}

	return lockEvent, nil

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
			data, err := DecodeEvent(contractABI, vLog)
			if err != nil {
				log.Fatal(err)
				continue
			}

			fmt.Println("----------------------------------------")
			txHash := vLog.TxHash

			assetID64 := vLog.Topics[1].Big().Uint64()
			assetID := uint16(assetID64)

			sourceChain, err := client.NetworkID(context.Background())
			if err != nil {
				log.Fatal(err)
				continue
			}

			user := common.HexToAddress(vLog.Topics[2].Hex())

			fmt.Println("txHash: ", txHash.Hex())
			fmt.Println("executor: ", user.Hex())
			fmt.Println("amount:", data.Amount)
			fmt.Println("assetID: ", assetID)
			fmt.Println("sourceChain: ", sourceChain.String())

			err = boundContract.Call(&callOpts, nil, "vote", txHash, user, data.Amount, assetID, sourceChain)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Vote successful!")
		}
	}
}
