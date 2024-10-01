package main

import (
	"fmt"

	c "github.com/dawitel/Ashok-reverse-proxy-test/internal/config"
	h "github.com/dawitel/Ashok-reverse-proxy-test/internal/handler"
)

func main() {
	cfg, err := c.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	srv := h.NewServer(h.ProxyHandler, cfg)
	srv.Start()
}
