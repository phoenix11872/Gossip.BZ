package internalEmbeddings

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

const defaultBatchSize = 100

type Embedder interface {
	EmbedDocuments(ctx context.Context, chunks [][]rune) ([][]float32, error)
}

type GeminiEmbedder struct {
	client    *genai.Client
	model     string
	batchSize int
	outputDim *int32
}
func NewGeminiEmbedder(client *genai.Client, model string, batchSize int, outputDim *int32)(*GeminiEmbedder, error){
	if client == nil{
		return nil, fmt.Errorf("genai client is required")
	}
	if model == ""{
		return nil, fmt.Errorf("embediing model is required")
	}
	if batchSize <=0{
		batchSize = defaultBatchSize
	}
	return &GeminiEmbedder{
		client: client,
		model: model,
		batchSize: batchSize,
		outputDim: outputDim,
	},nil
}