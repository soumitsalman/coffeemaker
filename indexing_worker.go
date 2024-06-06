package main

import (
	"log"

	"github.com/robfig/cron"
	sack "github.com/soumitsalman/coffeemaker/sdk/beansack"
	news "github.com/soumitsalman/coffeemaker/sdk/newscollector"
	reddit "github.com/soumitsalman/coffeemaker/sdk/redditor"
)

func StartIndexer() {
	c := cron.New()

	// initialize collectors
	nc := news.NewCollector(getSitemaps(), sack.AddBeans)
	rc := reddit.NewCollector(reddit.NewCollectorConfig(sack.AddBeans))

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
		sack.Rectify()
		<-coll_session
	})

	// TODO: remove this. Now we are running rectification after collection anyway
	// // run rectification
	// c.AddFunc(getRectifySchedule(), func() {
	// 	log.Println("[INDEXER] Running Rectification")

	// })

	// run clean up
	c.AddFunc(getCleanupSchedule(), func() {
		log.Println("[INDEXER] Running Cleanup")
		sack.Cleanup(30)
	})

	c.Start()
}
