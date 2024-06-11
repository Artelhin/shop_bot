package main

import (
	"fmt"
	"shop_bot/config"
)

func main() {
	fmt.Println("startup...")
	defer fmt.Println("app closed")

	cfg, err := config.FromFile("./config.yaml")
	if err != nil {
		fmt.Printf("can't cofigure: %s", err)
		return
	}
	fmt.Println(cfg)
}
