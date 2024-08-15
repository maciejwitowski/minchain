package lib

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

type TransactionsInput interface {
	InputChannel(ctx context.Context) <-chan string
}

type UserInput struct {
	reader *bufio.Reader
}

func NewUserInput() TransactionsInput {
	return &UserInput{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (ui *UserInput) InputChannel(ctx context.Context) <-chan string {
	messages := make(chan string)

	go func() {
		defer close(messages)

		for {
			select {
			case <-ctx.Done():
				log.Println("context cancelled")
			default:
				fmt.Print("> ")
				message, err := ui.reader.ReadString('\n')
				if err != nil {
					fmt.Println("Error reading the message from stdin:", err)
					break
				}
				message = strings.TrimSuffix(message, "\n")
				if message == "" {
					continue
				}
				messages <- message
			}
		}
	}()

	return messages
}
