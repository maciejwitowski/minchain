package lib

import (
	"crypto/ecdsa"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ListeningPort   int
	PrivateKey      *ecdsa.PrivateKey
	IsBlockProducer bool
	BlockTime       time.Duration
	Inputs          []string
}

const (
	INPUT_STDIN = "stdin"
	INPUT_API   = "api"
)

func InitConfig() Config {
	portStr := os.Getenv("P2P_PORT")
	if portStr == "" {
		log.Fatal("unknown p2p port")
	}
	port, _ := strconv.Atoi(portStr)

	isBlockProducerStr := os.Getenv("IS_BLOCK_PRODUCER")
	isBlockProducer := false
	if isBlockProducerStr == "true" {
		isBlockProducer = true
	}

	// Expects comma-separated inputs: cli, api
	var inputs []string
	inputsStr := os.Getenv("INPUTS")
	if inputsStr != "" {
		inputs = strings.Split(inputsStr, ",")
	}

	privateKey, err := ethcrypto.LoadECDSA(".pk")
	if err != nil {
		log.Fatal(err)
	}

	return Config{
		ListeningPort:   port,
		IsBlockProducer: isBlockProducer,
		PrivateKey:      privateKey,
		BlockTime:       5 * time.Second,
		Inputs:          inputs,
	}
}
