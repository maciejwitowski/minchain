package lib

import (
	"crypto/ecdsa"
	"flag"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

type Config struct {
	PubSubTopic   string
	ListeningPort int
	PrivateKey    *ecdsa.PrivateKey
}

func InitConfig() (*Config, error) {
	topic := flag.String("t", "", "pubsub topic")
	port := flag.Int("l", 0, "wait for incoming connections")
	flag.Parse()

	privateKey, err := ethcrypto.LoadECDSA(".pk")
	if err != nil {
		return nil, err
	}

	return &Config{
		PubSubTopic:   *topic,
		ListeningPort: *port,
		PrivateKey:    privateKey,
	}, nil
}
