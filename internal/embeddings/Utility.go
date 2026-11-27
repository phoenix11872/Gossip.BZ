package internalEmbeddings

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

func (e *GeminiEmbedder)EmbedDocuments(
ctx context.Context,
chunks [][]rune)([][]float32, error){
	return e.embedBatch(ctx, "RETRIEVAL_DOCUMENT", chunks)
}
func (e *GeminiEmbedder)EmbedQuery(ctx context.Context, query string)([]float32, error){
	if strings.TrimSpace(query) == ""{
		return nil, fmt.Errorf("query is empty")
	}
	vecs, err:= e.embedBatch(ctx, "RETRIEVAL_QUERY",[][]rune{[]rune(query)})
	if err!=nil{
		return nil, err
	}
	if len(vecs)!=1{
		return nil, fmt.Errorf("Expecpected 1 embeddings")
	}
	return vecs[0], nil
}



func (e *GeminiEmbedder) embedBatch(ctx context.Context,
	taskType string, 
	inputs [][]rune)([][]float32, error){
		if len(inputs) == 0{
			return nil, nil
		}
		results:= make([][]float32,0, len(inputs))
		cfg:= &genai.EmbedContentConfig{
			TaskType: taskType,
			OutputDimensionality: e.outputDim,
		}
		for start:=0;start<len(inputs);start+=e.batchSize{
			end := start + e.batchSize
			if end > len(inputs){
				end = len(inputs)
			}
			contents := make([]*genai.Content,0,end- start)
			for idx, chunk:= range inputs[start:end]{
				text:= string(chunk)
				if strings.TrimSpace(text)==""{
					return nil, fmt.Errorf("chunk at index %d is empty",idx);
				}
				contents = append(contents, genai.NewContentFromText(text, genai.RoleUser))
			}
			resp, err:= e.client.Models.EmbedContent(ctx, e.model, contents, cfg)
			if err!=nil{
				return nil, fmt.Errorf("embed batch [%d:%d]:%w",start, end, err)
				
			}
			if len(resp.Embeddings) != len(contents) {
				return nil, fmt.Errorf("embed batch [%d:%d]: got %d embeddings for %d inputs", start, end, len(resp.Embeddings), len(contents))
			}
			for idx, emb:=range resp.Embeddings{
				if emb == nil || len(emb.Values)==0{
					return nil, fmt.Errorf("empty embeding returned at index %d",start + idx)
				}
				results = append(results, emb.Values);
			}
		}
		return results, nil
	}