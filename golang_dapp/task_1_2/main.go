package main

import (
	"context"
	"fmt"
	"golang_dapp/counter"
	"golang_dapp/read"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// WaitForTransaction 等待交易被确认
func WaitForTransaction(client *ethclient.Client, txHash common.Hash) error {
	for {
		_, isPending, err := client.TransactionByHash(context.Background(), txHash)
		if err != nil {
			fmt.Println("eee...")
			return err
		}

		if !isPending {
			// 交易已被打包，再等待几个区块确认
			receipt, err := client.TransactionReceipt(context.Background(), txHash)
			if err != nil {
				return err
			}

			if receipt.Status == types.ReceiptStatusSuccessful {
				fmt.Println("Transaction confirmed successfully")
				return nil
			} else {
				return fmt.Errorf("transaction failed")
			}
		}

		fmt.Println("Transaction is pending, waiting...")
		time.Sleep(5 * time.Second) // 等待5秒后再次检查
	}
}
func main() {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// 拼接配置文件的完整路径
	configPath := filepath.Join(wd, "config/task_1_2.json")
	fmt.Println("configPath", configPath)

	config, err := read.ReadContractConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("config.InfuraUrl:", config.InfuraUrl)
	fmt.Println("config.ContractAddress:", config.ContractAddress)
	fmt.Println("config.DeployerPrivateKey:", config.DeployerPrivateKey)

	client, err := ethclient.Dial(config.InfuraUrl)
	if err != nil {
		log.Fatal(err)
	}
	usedContract, err := counter.NewCounter(common.HexToAddress(config.ContractAddress), client)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(config.DeployerPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155111))
	if err != nil {
		log.Fatal(err)
	}

	callOpt := &bind.CallOpts{Context: context.Background()}

	pre_count, err := usedContract.GetCount(callOpt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("pre_count:", pre_count)

	value := big.NewInt(1)
	tx, err := usedContract.Increment(opt, value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tx hash:", tx.Hash().Hex())

	time.Sleep(time.Second)
	// 等待交易确认
	fmt.Println("Waiting for transaction to be confirmed...")
	if err := WaitForTransaction(client, tx.Hash()); err != nil {
		log.Fatal(err)
	}

	cur_count, err := usedContract.GetCount(callOpt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("cur_count:", cur_count)

}
