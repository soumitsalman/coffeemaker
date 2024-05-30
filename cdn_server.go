package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/soumitsalman/beansack/sdk"
	"golang.org/x/time/rate"
)

const (
	_ERROR_MESSAGE   = "YO! do you even code?! Input format is fucked. Read this: https://github.com/soumitsalman/coffemaker."
	_SUCCESS_MESSAGE = "I gotchu."
)

type queryParams struct {
	Window int      `form:"window"`
	TopN   int      `form:"topn"`
	Kinds  []string `form:"kind"`
}

type bodyParams struct {
	Nuggets    []string    `json:"nuggets,omitempty"`
	Categories []string    `json:"categories,omitempty"`
	Embeddings [][]float32 `json:"embeddings,omitempty"`
	Context    string      `json:"context,omitempty"`
	URLs       []string    `json:"urls,omitempty"`
}

func extractParams(ctx *gin.Context) (*sdk.SearchOptions, []string) {
	options := sdk.NewSearchOptions()

	var query_params queryParams
	// if query params are mal-formed return error
	if ctx.ShouldBindQuery(&query_params) != nil {
		ctx.String(http.StatusBadRequest, _ERROR_MESSAGE)
		return nil, nil
	}
	if len(query_params.Kinds) > 0 {
		options.WithKind(query_params.Kinds)
	}
	if query_params.Window > 0 {
		options.WithTimeWindow(query_params.Window)
	}
	if query_params.TopN > 0 {
		options.WithTopN(query_params.TopN)
	}

	var body_params bodyParams
	// if body params are provided, assign them or else proceed without them
	if ctx.ShouldBindJSON(&body_params) == nil {
		options.SearchTexts = body_params.Categories
		options.SearchEmbeddings = body_params.Embeddings
		options.Context = body_params.Context
		options.WithURLs(body_params.URLs)
	}
	return options, body_params.Nuggets
}

func searchBeansHandler(ctx *gin.Context) {
	options, nuggets := extractParams(ctx)
	if options == nil {
		return
	}

	var res []sdk.Bean
	if len(nuggets) > 0 {
		res = sdk.NuggetSearch(nuggets, options)
	} else {
		res = sdk.FuzzySearch(options)
	}
	sendBeans(res, ctx)
}

func trendingBeansHandler(ctx *gin.Context) {
	options, _ := extractParams(ctx)
	if options == nil {
		return
	}
	ctx.JSON(http.StatusOK, sdk.TrendingBeans(options))
}

func retrieveBeansHandler(ctx *gin.Context) {
	options, _ := extractParams(ctx)
	if options == nil {
		return
	}
	ctx.JSON(http.StatusOK, sdk.Retrieve(options))
}

func trendingNuggetsHandler(ctx *gin.Context) {
	options, _ := extractParams(ctx)
	if options == nil {
		return
	}
	ctx.JSON(http.StatusOK, sdk.TrendingNuggets(options))
}

func initializeRateLimiter() gin.HandlerFunc {
	limiter := rate.NewLimiter(_RATE_LIMIT, _RATE_TPS)
	return func(ctx *gin.Context) {
		if limiter.Allow() {
			ctx.Next()
		} else {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
		}
	}
}

func sendBeans(res []sdk.Bean, ctx *gin.Context) {
	if len(res) > 0 {
		ctx.JSON(http.StatusOK, res)
	} else {
		ctx.Status(http.StatusNoContent)
	}
}

func newCDNServer() *gin.Engine {
	router := gin.Default()

	// NO NEED FOR AUTH: this is open to public
	open_group := router.Group("/")
	open_group.Use(initializeRateLimiter())
	// GET /beans
	open_group.GET("/beans", retrieveBeansHandler)
	// GET /beans/trending?window=1
	open_group.GET("/beans/trending", trendingBeansHandler)
	// GET /beans/search?window=1
	open_group.GET("/beans/search", searchBeansHandler)
	// GET /nuggets/trending?window=1
	open_group.GET("/nuggets/trending", trendingNuggetsHandler)

	return router
}

func RunCDN() {
	newCDNServer().Run()
}
