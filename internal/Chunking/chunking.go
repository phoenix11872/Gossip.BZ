package chunking

import (
	"fmt"
	docs "Gossipz.BZ/internal/Documents"
)


type Chunks struct {
	DocumentUid string
	Start       int
	End         int
	Embeddings  []float32
}

type ChunkingPolicies struct {
	Overlapping        int
	ChunkSize          int
	BoundaryConditions []string
}

func matchesBoundaryAt(content []rune, start int, boundary []rune) bool {
	if len(boundary) == 0 || start+len(boundary) > len(content) {
		return false
	}
	for i, r := range boundary {
		if content[start+i] != r {
			return false
		}
	}
	return true
}

func SimpleSlidingWindowChunker(p ChunkingPolicies, CurrentChunks Chunks) ([]Chunks, error){
	var ans []Chunks
	start := CurrentChunks.Start
	end := CurrentChunks.End
	step := p.ChunkSize - p.Overlapping
	if step <= 0 || p.ChunkSize <= 0 {
		return []Chunks{CurrentChunks}, fmt.Errorf("Chunk size not strictly bigger than overlap size")
	}
	for i := start; i <= end; i += step {
		chunkEnd := i + p.ChunkSize - 1
		if chunkEnd > end {
			chunkEnd = end
		}
		ans = append(ans, Chunks{DocumentUid: CurrentChunks.DocumentUid, Start: i, End: chunkEnd})
		if chunkEnd == end {
			break
		}
	}
	return ans, nil
}

func SplitIntoChunksByPolicy(content []rune, p ChunkingPolicies, bI int, CurrChunk Chunks) ([]Chunks, error) {
	currentSz := CurrChunk.End - CurrChunk.Start + 1
	if currentSz <= p.ChunkSize {
		return []Chunks{CurrChunk},nil
	}
	if bI >= len(p.BoundaryConditions) {
		return SimpleSlidingWindowChunker(p, CurrChunk)
	}

	boundary := []rune(p.BoundaryConditions[bI])
	if len(boundary) == 0 {
		return SplitIntoChunksByPolicy(content, p, bI+1, CurrChunk)
	}

	chunkStart := CurrChunk.Start
	end := CurrChunk.End
	var newChunks []Chunks

	for i := CurrChunk.Start; i <= end; i++ {
		if matchesBoundaryAt(content, i, boundary) {
			if chunkStart <= i-1 {
				newChunks = append(newChunks, Chunks{
					DocumentUid: CurrChunk.DocumentUid,
					Start:       chunkStart,
					End:         i - 1,
				})
			}
			chunkStart = i + len(boundary)
			i = chunkStart - 1
		}
	}

	if chunkStart <= end {
		newChunks = append(newChunks, Chunks{
			DocumentUid: CurrChunk.DocumentUid,
			Start:       chunkStart,
			End:         end,
		})
	}

	if len(newChunks) == 0 {
		return SplitIntoChunksByPolicy(content, p, bI+1, CurrChunk)
	}

	var finalAns []Chunks
	for _, value := range newChunks {
		if value.End < value.Start {
			continue
		}
		goodChunks , err:= SplitIntoChunksByPolicy(content, p, bI+1, value)
		if(err!=nil){
			return finalAns, err
		}
		finalAns = append(finalAns, goodChunks...)
	}
	return finalAns,nil

}

func SubDivideIntoChunksByParagraphs(doc docs.Document, p ChunkingPolicies) {
	_, _ = SubdivideIntoChunks(doc, p)
}

func SubdivideIntoChunks(doc docs.Document, p ChunkingPolicies) ([]Chunks, error) {
	chunkslist := []Chunks{}
	
	if p.ChunkSize <= 0 {
		return chunkslist, fmt.Errorf("chunk size must be greater than zero")
	}
	if p.Overlapping >= p.ChunkSize {
		return chunkslist, fmt.Errorf("Overlapping token window cannot be greater than chunks size")
	}
	content := []rune(doc.Content)
	if len(content) == 0 {
		return chunkslist, nil
	}
	initialChunk := Chunks{
		DocumentUid: doc.DocumentUid,
		Start:       0,
		End:         len(content) - 1,
	}
	return SplitIntoChunksByPolicy(content, p, 0, initialChunk)
}
