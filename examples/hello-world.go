package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Maniacal/go-xai-sdk/client"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	apiKey := os.Getenv("XAI_API_KEY") // Get from https://console.x.ai
	fmt.Println(apiKey)
	endpoint := "api.x.ai:443"

	c, err := client.New(apiKey, endpoint)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer c.Close()

	resp, err := c.Models.ListLanguageModels(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("list models failed: %v", err)
	}

	fmt.Println("Available Language Models:")
	for _, m := range resp.Models {
		fmt.Printf("- Name: %s, Version: %s, Created: %v\n", m.Name, m.Version, m.Created.AsTime())
	}
}
