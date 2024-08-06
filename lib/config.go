package lib

import (
	"crypto/ecdsa"
	"flag"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

type Config struct {
	PubSubTopic     string
	ListeningPort   int
	PrivateKey      *ecdsa.PrivateKey
	IsBlockProducer bool
}

func InitConfig() (*Config, error) {
	topic := flag.String("topic", "", "pubsub topic")
	port := flag.Int("port", 0, "wait for incoming connections")
	isBlockProducer := flag.Bool("block-producer", false, "whether this node is a block producer")
	flag.Parse()

	privateKey, err := ethcrypto.LoadECDSA(".pk")
	if err != nil {
		return nil, err
	}

	return &Config{
		PubSubTopic:     *topic,
		ListeningPort:   *port,
		IsBlockProducer: *isBlockProducer,
		PrivateKey:      privateKey,
	}, nil
}
