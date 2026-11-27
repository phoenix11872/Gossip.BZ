package vectorstore

import "context"

type ChunkVector struct {
	DocumentUid string
	Text        string
	Vector      []float32
	Index       int
}
type HIT struct {
	DocumentUid string
	Text        string
	Distance    float32
}

type VectorStore interface {
	EnsureIndex(ctx context.Context) error
	Upsert(ctx context.Context, chunks []ChunkVector) error
	Search(ctx context.Context, query []float32, k int) ([]HIT, error)
}
