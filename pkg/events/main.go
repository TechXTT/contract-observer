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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type LockEvent struct {
	AssetID     uint16
	Token       common.Address
	Amount      *big.Int
	User        common.Address
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

	fmt.Println("lockEventMap: ", lockEventMap)

	//parse map to struct
	assetID, ok := lockEventMap["assetID"].(uint16)
	if !ok {
		assetID = 0
		fmt.Println("Error parsing assetID")
	}

	token, ok := lockEventMap["token"].(common.Address)
	if !ok {
		token = common.HexToAddress("0x0")
		fmt.Println("Error parsing token")
	}

	amount, ok := lockEventMap["amount"].(*big.Int)
	if !ok {
		amount = big.NewInt(0)
	}

	user, ok := lockEventMap["user"].(common.Address)
	if !ok {
		user = common.HexToAddress("0x0")
	}

	targetChain, ok := lockEventMap["targetChain"].(*big.Int)
	if !ok {
		targetChain = big.NewInt(0)
	}

	lockEvent := LockEvent{
		AssetID:     assetID,
		Token:       token,
		Amount:      amount,
		User:        user,
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

	// boundContract := bind.NewBoundContract(contractAddress, contractABI, client, client, client)

	// callOpts := bind.CallOpts{
	// 	Pending: true,
	// 	From:    common.HexToAddress("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"),
	// }

	// err = boundContract.Call(&callOpts, nil, "lock", assetID, amount, targetChain)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	for {
		select {
		case err := <-subscription.Err():
			log.Fatal(err)
		case vLog := <-logs:
			data, err := DecodeEvent(contractABI, vLog)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("----------------------------------------")
			txHash := vLog.TxHash

			sourceChain := new(big.Int)
			sourceChain.SetString(vLog.Topics[1].Hex(), 0)

			fmt.Println("txHash: ", txHash.Hex())
			fmt.Println("executor: ", data.User.Hex())
			fmt.Println("amount:", data.Amount)
			fmt.Println("assetID: ", data.AssetID)
			fmt.Println("sourceChain: ", sourceChain)

			// err = boundContract.Call(&callOpts, nil, "vote", txHash, data.User, data.Amount, data.AssetID, sourceChain)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			fmt.Println("Vote successful!")
		}
	}
}
