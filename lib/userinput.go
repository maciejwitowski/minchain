package lib

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
)

func UserInput(ctx context.Context) <-chan string {
	messages := make(chan string)
	reader := bufio.NewReader(os.Stdin)

	go func() {
		defer close(messages)

		for {
			select {
			case <-ctx.Done():
				log.Println("context cancelled")
			default:
				fmt.Print("> ")
				message, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Error reading the message:", err)
				}
				messages <- message
			}
		}
	}()

	return messages
}
