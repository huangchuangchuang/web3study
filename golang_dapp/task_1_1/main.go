package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"golang_dapp/read"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthCtrl struct {
	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey
}

// NewEthTest 创建一个新的EthTest实例
func NewEthCtrl(clientAddr string, privKeyAddr string) *EthCtrl {
	client, err := ethclient.Dial(clientAddr)
	if err != nil {
		log.Fatal(err)
	}

	// 加载私钥
	privateKey, err := crypto.HexToECDSA(privKeyAddr) //3.55
	if err != nil {
		log.Fatal(err)
	}
	return &EthCtrl{
		client:     client,
		privateKey: privateKey,
	}
}

func (ec *EthCtrl) GetPublicAddress() (common.Address, error) {
	publicKey := ec.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	// 从公钥生成以太坊地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address, nil
}

func (ec *EthCtrl) AddDelay() {
	// 添加延迟; 连续的 API 调用可能超过了 Infura 免费账户的速率限制（通常是每秒几个请求）
	time.Sleep(1000 * time.Millisecond)
}

func (ec *EthCtrl) transfer(value *big.Int, toPubAddr string) {
	// value := big.NewInt(amount) // in wei (0.01 eth)
	gasLimit := uint64(21000) // in units
	gasPrice, err := ec.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("gasLimit:", gasLimit)
	fmt.Println("gasPrice:", gasPrice)
	//  ETH 发送给谁（公钥地址）
	toPubAddress := common.HexToAddress(toPubAddr) // 1.59
	var data []byte

	fromAddress, err := ec.GetPublicAddress()
	if err != nil {
		log.Fatal(err)
	}

	ec.AddDelay() // 添加延迟
	nonce, err := ec.client.PendingNonceAt(context.Background(), fromAddress)
	// nonce, err := ec.GetNonce()
	if err != nil {
		log.Fatal(err)
	}

	// 生成我们的未签名以太坊事务
	tx := types.NewTransaction(nonce, toPubAddress, value, gasLimit, gasPrice, data)

	ec.AddDelay() // 添加延迟；
	chainID, err := ec.client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 发件人的私钥对事务进行签名
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), ec.privateKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("signedTx:", signedTx)

	// 将已签名的事务广播到整个网络
	ec.AddDelay() // 添加延迟；
	err = ec.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())

}

func main() {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// 拼接配置文件的完整路径
	configPath := filepath.Join(wd, "config/task_1_1.json")
	fmt.Println("configPath", configPath)

	config, err := read.ReadTransferConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("config.InfuraUrl:", config.InfuraUrl)
	fmt.Println("config.UserPrivateKey:", config.UserPrivateKey)
	fmt.Println("config.ToAddress:", config.ToAddress)

	eth_ctrl := NewEthCtrl(
		config.InfuraUrl,
		config.UserPrivateKey, // 私钥；account5
	)

	eth_ctrl.transfer(big.NewInt(10000000000000000), config.ToAddress)

}
