package examples

import (
	"log"
	"os"
	"time"

	nc "github.com/soumitsalman/coffeemaker/sdk/newscollector"
)

func ScrapeAndStoreLocally() {
	start_time := time.Now()
	// initialize to save locally
	collector := nc.NewCollector(os.Getenv("SITEMAPS_FILE"), localFileStore)
	collector.Collect()
	log.Println("Collection took", time.Since(start_time))
}
