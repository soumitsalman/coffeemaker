package nlp

import (
	ctx "context"
	"log"
	"os"
	"time"

	"github.com/avast/retry-go"
	hfemb "github.com/tmc/langchaingo/embeddings/huggingface"
	hfllm "github.com/tmc/langchaingo/llms/huggingface"
	"github.com/tmc/langchaingo/textsplitter"
)

const (
	_RETRY_DELAY = 20 * time.Second
)

const (
	_EMBEDDINGS_MODEL = "nomic-ai/nomic-embed-text-v1"
	// _KEYWORDS_MODEL   = "ilsilfverskiold/tech-keywords-extractor"
	_SUMMARY_MODEL   = "google/flan-t5-base"
	_TEXT_CHUNK_SIZE = 8192
)

type HuggingfaceDriver struct {
	// small_embedder *hfemb.Huggingface
	embedder       *hfemb.Huggingface
	text_splitter  textsplitter.TokenSplitter
	keywords_model *hfllm.LLM
	summary_moodel *hfllm.LLM
}

func NewHuggingfaceDriver() *HuggingfaceDriver {
	emb_llm, _ := hfllm.New(hfllm.WithToken(getHuggingfaceToken()))
	embedder, err := hfemb.NewHuggingface(hfemb.WithClient(*emb_llm), hfemb.WithModel(_EMBEDDINGS_MODEL))
	if err != nil {
		log.Printf("[NewHuggingfaceDriver] Failed Loading %s. %v\n", _EMBEDDINGS_MODEL, err)
		return nil
	}
	// keywords_model, err := hfllm.New(hfllm.WithToken(getHuggingfaceToken()), hfllm.WithModel(_KEYWORDS_MODEL))
	// if err != nil {
	// 	log.Printf("[NewHuggingfaceDriver] Failed Loading %s. %v\n", _KEYWORDS_MODEL, err)
	// 	return nil
	// }
	summary_model, err := hfllm.New(hfllm.WithToken(getHuggingfaceToken()), hfllm.WithModel(_SUMMARY_MODEL))
	if err != nil {
		log.Printf("[NewHuggingfaceDriver] Failed Loading %s. %v\n", _SUMMARY_MODEL, err)
		return nil
	}
	return &HuggingfaceDriver{
		text_splitter: textsplitter.NewTokenSplitter(textsplitter.WithChunkSize(_TEXT_CHUNK_SIZE)),
		embedder:      embedder,
		// keywords_model: keywords_model,
		summary_moodel: summary_model,
	}
}

func (driver *HuggingfaceDriver) CreateTextEmbeddings(text string) ([]float32, error) {
	var res []float32
	retry.Do(func() error {
		vecs, err := driver.embedder.EmbedDocuments(ctx.Background(), []string{text})
		if err != nil {
			log.Printf("[Huggingface Driver | %s]: error generating embeddings.%v\n", _EMBEDDINGS_MODEL, err)
			return err
		}
		res = vecs[0]
		return nil
	}, retry.Delay(_RETRY_DELAY))
	return res, nil
}

func (driver *HuggingfaceDriver) CreateBatchTextEmbeddings(texts []string) ([][]float32, error) {
	var res [][]float32
	retry.Do(func() error {
		vecs, err := driver.embedder.EmbedDocuments(ctx.Background(), texts)
		if err != nil {
			log.Printf("[Huggingface Driver | %s]: error generating embeddings.%v\n", _EMBEDDINGS_MODEL, err)
			return err
		}
		res = vecs
		return nil
	}, retry.Delay(_RETRY_DELAY))
	return res, nil
}

func getHuggingfaceToken() string {
	return os.Getenv("HUGGINGFACE_API_TOKEN")
}
