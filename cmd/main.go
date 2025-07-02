package main

import (
	"fmt"

	"TelegramBot/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
}
