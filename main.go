package main

import (
	"log"

	"github.com/joho/godotenv"
	sack "github.com/soumitsalman/coffeemaker/sdk/beansack"
)

func main() {
	godotenv.Load()

	if err := sack.InitializeBeanSack(getDBConnectionString(), getEmbedderUrl(), getEmbedderCtx(), getLLMServiceAPIKey()); err != nil {
		log.Fatalln("Initialization not working", err)
	}

	switch getInstanceMode() {
	case "CDN":
		log.Println("Running in CDN Mode.")
		RunCDN()
	case "INDEXER":
		log.Println("Running in Indexer Mode.")
		StartIndexer()
		// this is not a blocking call so wait for jobs
		select {}
	case "DUAL":
		log.Println("Running in Dual Mode.")
		StartIndexer()
		// this is blocking call so this should always be sequenced after the indexer starts
		// no need to use select {} after indexer since the thread will be blocked on RunCDN
		RunCDN()
	case "DEBUG":
		RunDebug()
	default:
		log.Println("WTF is this?!")
	}
	log.Println("[coffeemaker] shutting down")

}
