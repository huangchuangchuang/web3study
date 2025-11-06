package read

import (
	"encoding/json"
	"fmt"
	"os"
)

type ContractConfig struct {
	InfuraUrl          string `json:"infura_url"`
	ContractAddress    string `json:"contract_address"`
	DeployerPrivateKey string `json:"deployer_private_key"`
}

func ReadContractConfig(filename string) (ContractConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return ContractConfig{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config ContractConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return ContractConfig{}, fmt.Errorf("failed to decode config file: %w", err)
	}
	return config, nil

}

type TransferConfig struct {
	InfuraUrl      string `json:"infura_url"`
	UserPrivateKey string `json:"user_private_key"`
	ToAddress      string `json:"to_address"`
}

func ReadTransferConfig(filename string) (TransferConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return TransferConfig{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config TransferConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return TransferConfig{}, fmt.Errorf("failed to decode config file: %w", err)
	}
	return config, nil

}
