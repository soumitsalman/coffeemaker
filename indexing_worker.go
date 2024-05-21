package main

import (
	"log"

	"github.com/robfig/cron"
	"github.com/soumitsalman/beansack/sdk"
	reddit "github.com/soumitsalman/go-reddit/collector"
	news "github.com/soumitsalman/newscollector/collector"
)

// func newBeansHandler(ctx *gin.Context) {
// 	var beans []sdk.Bean
// 	if ctx.BindJSON(&beans) != nil {
// 		ctx.String(http.StatusBadRequest, _ERROR_MESSAGE)
// 	} else {
// 		sdk.AddBeans(beans)
// 		ctx.String(http.StatusOK, _SUCCESS_MESSAGE)
// 	}
// }

// func serverAuthHandler(ctx *gin.Context) {
// 	// log.Println(ctx.GetHeader("X-API-Key"), getInternalAuthToken())
// 	if ctx.GetHeader("X-API-Key") == getInternalAuthToken() {
// 		ctx.Next()
// 	} else {
// 		ctx.AbortWithStatus(http.StatusUnauthorized)
// 	}
// }

// // func rectifyHandler(ctx *gin.Context) {
// // 	go sdk.Rectify()
// // 	ctx.String(http.StatusOK, _SUCCESS_MESSAGE)
// // }

// func newIndexer() *gin.Engine {
// 	router := gin.Default()

// 	// SERVICE TO SERVICE AUTH
// 	auth_group := router.Group("/")
// 	auth_group.Use(initializeRateLimiter(), serverAuthHandler)
// 	// PUT /beans
// 	auth_group.PUT("/beans", newBeansHandler)
// 	auth_group.POST("/rectify", rectifyHandler)

// 	return router
// }

const _SITEMAPS_PATH = "./sitemaps.csv"

func StartIndexer() {
	c := cron.New()

	save_to_beansack := func(beans []sdk.Bean) { sdk.AddBeans(beans) }

	// initialize collectors
	nc := news.NewCollector(_SITEMAPS_PATH, save_to_beansack)
	rc := reddit.NewCollector(reddit.NewCollectorConfig(save_to_beansack))

	// run collection
	c.AddFunc(getCollectionSchedule(), func() {
		log.Println("[INDEXER] Running news collector")
		nc.Collect()
		log.Println("[INDEXER] Running reddit collector")
		rc.Collect()
	})

	// run clean up
	c.AddFunc(getCleanupSchedule(), func() {
		log.Println("[INDEXER] Running Cleanup")
		sdk.Cleanup(30)
	})

	c.Start()
}
