package lib

import (
	"crypto/ecdsa"
	"flag"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"log"
)

type Config struct {
	ListeningPort   int
	PrivateKey      *ecdsa.PrivateKey
	IsBlockProducer bool
}

func InitConfig() Config {
	port := flag.Int("port", 0, "wait for incoming connections")
	isBlockProducer := flag.Bool("block-producer", false, "whether this node is a block producer")
	flag.Parse()

	privateKey, err := ethcrypto.LoadECDSA(".pk")
	if err != nil {
		log.Fatal(err)
	}

	return Config{
		ListeningPort:   *port,
		IsBlockProducer: *isBlockProducer,
		PrivateKey:      privateKey,
	}
}
