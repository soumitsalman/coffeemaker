package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/soumitsalman/beansack/sdk"
)

func main() {
	godotenv.Load()

	if err := sdk.InitializeBeanSack(getDBConnectionString(), getEmbedderUrl(), getLLMServiceAPIKey()); err != nil {
		log.Fatalln("Initialization not working", err)
	}

	switch getInstanceMode() {
	case "CDN":
		RunCDN()
	case "DEBUG":
		RunDebug()
	case "INDEXER":
		StartIndexer()
		// this is not a blocking call so wait for jobs
		select {}
	case "DUAL":
		StartIndexer()
		// this is blocking call so this should always be sequenced after the indexer starts
		// no need to use select {} after indexer since the thread will be blocked on RunCDN
		RunCDN()
	}
	log.Println("[coffeemaker] shutting down")

}
