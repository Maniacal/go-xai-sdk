package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Maniacal/go-xai-sdk/client"
	xai "github.com/Maniacal/go-xai-sdk/xai/api/v1"
)

func main() {
	apiKey := os.Getenv("XAI_API_KEY") // Get from https://console.x.ai
	endpoint := "api.x.ai:443"

	// Create a chat client
	c, err := client.New(apiKey, endpoint)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer c.Close()

	// Maintain conversation history
	var messages []*xai.Message
	messages = append(messages, &xai.Message{
		Role: xai.MessageRole_ROLE_SYSTEM,
		Content: []*xai.Content{
			{
				Content: &xai.Content_Text{
					Text: "You are a helpful AI assistant.",
				},
			},
		},
	})

	// Main loop to read user input and send messages
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Start chatting (type 'exit' to quit):")
	for {
		fmt.Print("You: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		// Append user message to history
		messages = append(messages, &xai.Message{
			Role: xai.MessageRole_ROLE_USER,
			Content: []*xai.Content{
				{
					Content: &xai.Content_Text{
						Text: input,
					},
				},
			},
		})

		// Create request with full history
		req := &xai.GetCompletionsRequest{
			Messages: messages,
			Model:    "grok-3-mini", // Specify the model to use
		}

		// Start streaming
		stream, err := c.Chat.GetCompletionChunk(context.Background(), req)
		if err != nil {
			fmt.Printf("Failed to start stream: %v\n", err)
			continue
		}

		// Receive and print response
		var assistantContent string
		fmt.Print("AI: ")
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Printf("\nError receiving message: %v\n", err)
				break
			}
			for _, choice := range resp.Choices {
				if choice.Delta != nil && choice.Delta.Content != "" {
					fmt.Print(choice.Delta.Content)
					assistantContent += choice.Delta.Content
				}
			}
		}
		fmt.Println()

		// Append assistant's response to history
		messages = append(messages, &xai.Message{
			Role: xai.MessageRole_ROLE_ASSISTANT,
			Content: []*xai.Content{
				{
					Content: &xai.Content_Text{
						Text: assistantContent,
					},
				},
			},
		})
	}
}
