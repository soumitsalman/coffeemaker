package nlp

import (
	"fmt"
	"log"

	datautils "github.com/soumitsalman/data-utils"
)

const (
	SEARCH_QUERY    = "search_query"
	SEARCH_DOCUMENT = "search_document"
	CLASSIFICATION  = "classification"
	SIMILARITY      = "clustering"
)

type inferenceInput struct {
	Inputs []string `json:"inputs"`
}

type EmbeddingServerError string

func (err EmbeddingServerError) Error() string {
	return string(err)
}

type EmbeddingsDriver struct {
	embed_url string
	ctx       int
	// splitter  textsplitter.TokenSplitter
}

func NewEmbeddingsDriver(url string, ctx int) *EmbeddingsDriver {
	return &EmbeddingsDriver{
		embed_url: url,
		ctx:       ctx,
	}
}

func (driver *EmbeddingsDriver) CreateBatchTextEmbeddings(texts []string, task_type string) [][]float32 {
	// if the count is over the window size split in half and try
	if CountTokens(texts) > driver.ctx {
		return append(
			driver.CreateBatchTextEmbeddings(texts[:len(texts)/2], task_type),
			driver.CreateBatchTextEmbeddings(texts[len(texts)/2:], task_type)...)
	}
	input_texts := datautils.Transform(texts, func(item *string) string { return driver.toEmbeddingInput(*item, task_type) })
	embs := driver.createEmbeddings(&inferenceInput{input_texts})
	// if the embeddings generation is failing insert duds
	if embs == nil {
		return make([][]float32, len(texts))
	}
	return embs
}

func (driver *EmbeddingsDriver) CreateTextEmbeddings(text string, task_type string) []float32 {
	output := driver.createEmbeddings(&inferenceInput{[]string{driver.toEmbeddingInput(text, task_type)}})
	if len(output) >= 1 {
		return output[0]
	}
	return nil
}

func (driver *EmbeddingsDriver) toEmbeddingInput(text, task_type string) string {
	if len(task_type) > 0 {
		text = fmt.Sprintf("%s: %s", task_type, text)
	}
	return text
}

func (driver *EmbeddingsDriver) createEmbeddings(input *inferenceInput) [][]float32 {
	return retryT(
		func() ([][]float32, error) {
			if embs, err := postHTTPRequest[[][]float32](driver.embed_url, "", input); err != nil {
				log.Printf("[EmbeddingsDriver] Embedding generation failed. %v\n", err)
				return nil, err
			} else if len(embs) != len(input.Inputs) {
				err_msg := fmt.Sprintf("[EmbeddingsDriver] Embedding generation failed. Expected number of embeddings %d. Generated number of embeddings: %d", len(input.Inputs), len(embs))
				log.Println(err_msg)
				return nil, EmbeddingServerError(err_msg)
			} else {
				return embs, nil
			}
		})
}
