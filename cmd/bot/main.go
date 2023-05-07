package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	client "github.com/sgoldenf/vkbot/internal/vk_client"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("WARNING: No .env file found")
	}
}

func main() {
	vkClient := client.New(os.Getenv("ACCESS_TOKEN"), os.Getenv("GROUP_ID"))
	for {
		vkClient.Poll()
	}
}
