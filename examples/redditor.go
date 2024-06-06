package examples

import (
	"github.com/soumitsalman/coffeemaker/sdk/redditor"
)

func RedditAndStoreLocally() {
	config := redditor.NewCollectorConfig(localFileStore)
	redditor.NewCollector(config).Collect()
}
