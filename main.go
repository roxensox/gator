package main

import (
	"fmt"
	"github.com/roxensox/gator/internal/config"
)

func main() {
	cfg := config.NewConfig()
	cfg.Read()
	cfg.SetUser("ryan")
	cfg.Read()
	fmt.Println("Contents:")
	fmt.Printf(" - %s\n", *cfg.DB_URL)
	fmt.Printf(" - %s\n", *cfg.CurrentUser)
}
