package common

import (
	"flag"
)

type Config struct {
	PubSubTopic   string
	ListeningPort int
}

func InitConfig() Config {
	topic := flag.String("t", "", "pubsub topic")
	port := flag.Int("l", 0, "wait for incoming connections")
	flag.Parse()
	return Config{
		PubSubTopic:   *topic,
		ListeningPort: *port,
	}
}
