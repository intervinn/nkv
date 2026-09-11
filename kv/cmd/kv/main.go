package main

import (
	"log"

	"github.com/intervinn/nkv/kv"
)

func main() {
	srv := kv.NewServer(kv.NewMap())
	if err := srv.Start(":8080"); err != nil {
		log.Fatalln("failed to start:", err)
	}

	log.Println("starting kv...")
	srv.WaitForShutdown()
}
