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
		RunIndexer()
	case "DUAL":
		RunIndexer()
		RunCDN()
	}
	select {} // wait for jobs
}
