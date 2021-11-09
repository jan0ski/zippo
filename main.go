package main

import (
	zippo "github.com/jan0ski/zippo/pkg"
)

func main() {
	s := &zippo.Server{
		Config: &zippo.ServerConfig{
			ButaneTemplate: "./config.yaml",
		},
	}
	s.Run()
}
