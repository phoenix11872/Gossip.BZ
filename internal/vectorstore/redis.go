package vectorstore

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)
var _ VectorStore = (*RedisStore)(nil)
type RedisStore struct {
	client    *redis.Client
	indexName string
	keyPrefix string
	dim       int
}

func NewRedisStore(client *redis.Client, dim int) (*RedisStore,error) {
	if client == nil{
		return nil, fmt.Errorf("client is nil")
	}
	if dim <= 0{
		return nil, fmt.Errorf("dim is less than 0")
	}
	return &RedisStore{
		client: client,
		indexName: "idx:chunks",
		keyPrefix: "chunk:",
		dim: dim,
	}, nil
}

func (s *RedisStore) EnsureIndex(ctx context.Context) error {
	_, err:= s.client.FTCreate(ctx, s.indexName, 
		&redis.FTCreateOptions{
			OnHash: true,
			Prefix: []interface{}{s.keyPrefix}, 
		}, 
		&redis.FieldSchema{
			FieldName: "document_uid",
			FieldType: redis.SearchFieldTypeTag,
		}, 
		&redis.FieldSchema{
			FieldName: "text",
			FieldType: redis.SearchFieldTypeText,
		}, 
		&redis.FieldSchema{
			FieldName: "embeddings",
			FieldType: redis.SearchFieldTypeVector
			VectorArgs: &redis.FTVectorArgs{
				HNSWOptions: &redis.FTHNSWOptions{
					Type: "FLOAT32",
					Dim: s.dim, 
					DistanceMetric: "COSINE",
				},
			},
		}, 
	).Result()
	if err!= nil && strings.Contains(err.Error(), "index already exists"){
		return nil
	}
	return err
}
func validateChunk(chunk ChunkVector)error{
	if chunk.DocumentUid == ""{
		return fmt.Errorf("document uid is required")
	}
	if chunk.Text == ""{
		return fmt.Errorf("text is required")
	}
	if len(chunk.Vector) != s.dim{
		return fmt.Errorf("vector dimension mismatch")
	}
	return nil
}
func chunkKey(chunkChunkVector) string{
	return fmt.Sprintf("%s:%d", chunk.DocumentUid, chunk.Index)
}
func float32ToBlob(values []float32)[]byte{
	blob:= make([]byte, 4*len(values))
	for i, value:= range values{
		binary.LittleEndian.PutUint32(blob[i*4:], math.Float32bits(value))
	}
	return blob
}

func (s *RedisStore)Upsert(ctx context.Context, chunks[]ChunkVector) error {
	if len(chunks) == 0{
		return nil
	}
	pipe:= s.client.Pipeline()
	for i, chunk:= range chunks{
		if err:= validateChunk(chunk); err!= nil{
			return fmt.Errorf("chunk %d: %w", i, err)
		}
		pipe.HSet(ctx, s.keyPrefix+chunkKey(chunk), 
		"document_uid", chunk.DocumentUid, 
		"text", chunk.Text,
		"embeddings", float32ToBlob(chunk.Vector))
	}
	_, err:= pipe.Exec(ctx);err!=nil{
		return fmt.Errorf("failed to upsert chunks: %w", err)
	}
	return nil
}

func (s *RedisStore) Search(ctx context.Context, query[]float32, k int)([]HIT, error){
	if k<=0{
		return nil, fmt.Errorf("k must be greater than 0")
	}
	if len(query) != s.dim{
		return nil, fmt.Errorf("query dimension mismatch")
	}
	//* is wildcard pre filter which matches any document with specified requirements
	res, err:= s.client.FTSearchWithArgs(ctx, s.indexName,
		fmt.Sprintf("*=>[KNN %d @embeddings] $vec AS distance]",k),
	    &redis.FTSearchOptions{
			Return: []redis.FTSearchReturn{
				{FieldName: "document_uid"},
				{FieldName: "text"}, 
				{FieldName: "distance"},
			},
			SortBy: []redis.FTSearchSortBy{{FieldName: "distance", Asc: true}},
			Limit: k, 
			DialectVersion: 2,
			Params: map[string]interface{}{
				"vec": float32ToBlob(query),
			},
		},
	).Result()
	if err!= nil{
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	hits:= make([]HIT, 0, len(res.Docs))
	for _, doc:= range res.Docs{
		distance, err:= strconv.
	}
	return hits, nil
}



