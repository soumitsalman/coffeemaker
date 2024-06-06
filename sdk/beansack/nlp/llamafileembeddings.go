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

type EmbeddingsRequest struct {
	Inputs []string `json:"content"`
}

type embedddingResult struct {
	Embedding []float32 `json:"embedding,omitempty"`
}

type EmbeddingResponse struct {
	Results []embedddingResult `json:"results,omitempty"`
}

type EmbeddingServerError string

func (err EmbeddingServerError) Error() string {
	return string(err)
}

type EmbeddingsDriver struct {
	Url string
	Ctx int
}

func NewLlamaFileDriver(base_url string, ctx int) *EmbeddingsDriver {
	return &EmbeddingsDriver{
		Url: base_url + "/embedding",
		Ctx: ctx,
	}
}

func (driver *EmbeddingsDriver) CreateBatchTextEmbeddings(texts []string, task_type string) [][]float32 {
	// if the count is over the window size split in half and try
	if CountTokens(texts) > driver.Ctx {
		return append(
			driver.CreateBatchTextEmbeddings(texts[:len(texts)/2], task_type),
			driver.CreateBatchTextEmbeddings(texts[len(texts)/2:], task_type)...)
	}
	input_texts := datautils.Transform(texts, func(item *string) string { return driver.toEmbeddingInput(*item, task_type) })
	embs := driver.createEmbeddings(&EmbeddingsRequest{input_texts})
	// if the embeddings generation is failing insert duds
	if embs == nil {
		return make([][]float32, len(texts))
	}
	return embs
}

func (driver *EmbeddingsDriver) CreateTextEmbeddings(text string, task_type string) []float32 {
	output := driver.createEmbeddings(&EmbeddingsRequest{[]string{driver.toEmbeddingInput(text, task_type)}})
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

func (driver *EmbeddingsDriver) createEmbeddings(input *EmbeddingsRequest) [][]float32 {
	return retryT(
		func() ([][]float32, error) {
			if embs, err := postHTTPRequest[EmbeddingResponse](driver.Url, "", input); err != nil {
				log.Printf("[EmbeddingsDriver] Embedding generation failed. %v\n", err)
				return nil, err // return a dud
			} else if len(embs.Results) != len(input.Inputs) {
				err_msg := fmt.Sprintf("[EmbeddingsDriver] Embedding generation failed. Expected number of embeddings %d. Generated number of embeddings: %d", len(input.Inputs), len(embs.Results))
				log.Println(err_msg)
				return nil, EmbeddingServerError(err_msg) // return a dud
			} else {
				return datautils.Transform(embs.Results, func(item *embedddingResult) []float32 { return item.Embedding }), nil
			}
		})
}
