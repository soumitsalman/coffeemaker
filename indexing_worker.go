package main

import (
	"log"

	"github.com/robfig/cron"
	"github.com/soumitsalman/beansack/sdk"
	reddit "github.com/soumitsalman/go-reddit/collector"
	news "github.com/soumitsalman/newscollector/collector"
)

const _SITEMAPS_PATH = "./sitemaps.csv"

func StartIndexer() {
	c := cron.New()

	// initialize collectors
	nc := news.NewCollector(_SITEMAPS_PATH, sdk.AddBeans)
	rc := reddit.NewCollector(reddit.NewCollectorConfig(sdk.AddBeans))

	// the channel is to keep collection synchronization
	// if a collection session is already in progress, then the next collection instruction will wait until this session is finished
	coll_session := make(chan bool, 1)
	// run collection
	c.AddFunc(getCollectionSchedule(), func() {
		// start collection session
		coll_session <- true
		log.Println("[INDEXER] Running news collector")
		nc.Collect()
		log.Println("[INDEXER] Running reddit collector")
		rc.Collect()
		// finish collection session so that the next session can continue
		<-coll_session
	})

	// run rectification
	c.AddFunc(getRectifySchedule(), func() {
		log.Println("[INDEXER] Running Rectification")
		sdk.Rectify()
	})

	// run clean up
	c.AddFunc(getCleanupSchedule(), func() {
		log.Println("[INDEXER] Running Cleanup")
		sdk.Cleanup(30)
	})

	c.Start()
}
