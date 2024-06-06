package newscollector

import (
	"encoding/csv"
	"log"
	"os"

	ds "github.com/soumitsalman/coffeemaker/sdk/beansack"
	datautils "github.com/soumitsalman/data-utils"
)

type NewsSiteCollector struct {
	site_loaders []*WebLoader
	store_func   func([]ds.Bean)
}

func NewCollector(sitemaps string, store_func func([]ds.Bean)) NewsSiteCollector {
	return NewsSiteCollector{
		site_loaders: createSiteLoaders(sitemaps),
		store_func:   store_func,
	}
}

func (collector NewsSiteCollector) Collect() {
	for _, loader := range collector.site_loaders {
		docs := loader.LoadSite()
		log.Println(len(docs), "new beans found from", loader.Config.Sitemap)
		// storeNewBeans(docs)
		collector.store_func(toBeans(docs))
	}
}

func readSitemapsCSV(sitemaps string) [][]string {
	file, _ := os.Open(sitemaps)
	defer file.Close()
	items, _ := csv.NewReader(file).ReadAll()
	// ignore the header
	return items[1:]
}

func createSiteLoaders(sitemaps string) []*WebLoader {
	site_loaders := datautils.Transform(readSitemapsCSV(sitemaps), func(item *[]string) *WebLoader {
		return NewDefaultNewsSitemapLoader(2, (*item)[0])
	})
	return append(site_loaders,
		// this is a specialied loader
		NewYCHackerNewsSiteLoader(),
	)
}

func toBeans(docs []*Article) []ds.Bean {
	beans := make([]ds.Bean, len(docs))
	for i, doc := range docs {
		beans[i].Url = doc.URL
		beans[i].Source = doc.Source
		beans[i].Title = doc.Title
		beans[i].Kind = ds.ARTICLE
		beans[i].Text = doc.Text
		beans[i].Author = doc.Author
		beans[i].Created = doc.PublishDate
		beans[i].Keywords = doc.Keywords
		if doc.Comments > 0 || doc.Likes > 0 {
			beans[i].MediaNoise = &ds.MediaNoise{
				BeanUrl:       doc.URL,
				Source:        doc.Source,
				Comments:      doc.Comments,
				ThumbsupCount: doc.Likes,
			}
		}
	}
	return beans
}
