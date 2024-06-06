package beansack

import (
	"github.com/soumitsalman/coffeemaker/sdk/beansack/nlp"
	"github.com/soumitsalman/coffeemaker/sdk/beansack/store"
)

const (
	BEANSACK    = "beansack"
	BEANS       = "beans"
	NOISES      = "noises"
	KEYWORDS    = "keywords"
	NEWSNUGGETS = "concepts"
)

var (
	beanstore   *store.Store[Bean]
	nuggetstore *store.Store[BeanNugget]
	noisestore  *store.Store[MediaNoise]
	embedder    *nlp.EmbeddingsDriver
	pb_client   *nlp.ParrotboxClient
)

const (
	// _SEARCH_EMB = "search_embeddings"
	_CLASSIFICATION_EMB = "category_embeddings"
	_SUMMARY            = "summary"
)

type BeanSackError string

func (err BeanSackError) Error() string {
	return string(err)
}

func InitializeBeanSack(db_conn_str, emb_url string, emb_ctx int, pb_auth_token string) error {
	beanstore = store.New(db_conn_str, BEANSACK, BEANS,
		// store.WithMinSearchScore[Bean](0.55), // TODO: change this to 0.8 in future
		// store.WithSearchTopN[Bean](10),
		store.WithDataIDAndEqualsFunction(getBeanId, Equals),
	)
	noisestore = store.New[MediaNoise](db_conn_str, BEANSACK, NOISES)
	nuggetstore = store.New[BeanNugget](db_conn_str, BEANSACK, NEWSNUGGETS)

	if beanstore == nil || nuggetstore == nil {
		return BeanSackError("Initialization Failed. db_conn_str Not working.")
	}

	pb_client = nlp.NewParrotboxClient(pb_auth_token)
	embedder = nlp.NewLlamaFileDriver(emb_url, emb_ctx)

	return nil
}
