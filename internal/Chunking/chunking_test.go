package chunking

import (
	"reflect"
	"testing"
	docs "Gossipz.BZ/internal/Documents"
)

func chunkRanges(chunks []Chunks) [][2]int {
	var ranges [][2]int
	for _,chunk:= range chunks{
		ranges = append(ranges, [2]int{chunk.Start, chunk.End})
	}
	return ranges
}

	
	

func TestSubdivideIntoChunksBoundaryAndLengthCases(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		policy     ChunkingPolicies
		wantRanges [][2]int
	}{
		{
			name:    "empty document",
			content: "",
			policy:  ChunkingPolicies{ChunkSize: 4},
		},
		{
			name:       "shorter than chunk size",
			content:    "abc",
			policy:     ChunkingPolicies{ChunkSize: 4},
			wantRanges: [][2]int{{0, 2}},
		},
		{
			name:       "exactly chunk size",
			content:    "abcd",
			policy:     ChunkingPolicies{ChunkSize: 4},
			wantRanges: [][2]int{{0, 3}},
		},
		{
			name:       "one character over chunk size",
			content:    "abcde",
			policy:     ChunkingPolicies{ChunkSize: 4},
			wantRanges: [][2]int{{0, 3}, {4, 4}},
		},
		{
			name:       "overlapping sliding windows",
			content:    "abcdefghij",
			policy:     ChunkingPolicies{ChunkSize: 4, Overlapping: 2},
			wantRanges: [][2]int{{0, 3}, {2, 5}, {4, 7}, {6, 9}},
		},
		{
			name:       "empty boundary falls back to sliding windows",
			content:    "abcdef",
			policy:     ChunkingPolicies{ChunkSize: 3, BoundaryConditions: []string{""}},
			wantRanges: [][2]int{{0, 2}, {3, 5}},
		},
		{
			name:       "unmatched boundary falls back to sliding windows",
			content:    "abcdef",
			policy:     ChunkingPolicies{ChunkSize: 3, BoundaryConditions: []string{"|"}},
			wantRanges: [][2]int{{0, 2}, {3, 5}},
		},
		{
			name:    "boundary at the beginning and end",
			content: "\nabc\n",
			policy: ChunkingPolicies{
				ChunkSize:          2,
				BoundaryConditions: []string{"\n"},
			},
			wantRanges: [][2]int{{1, 2}, {3, 3}},
		},
		{
			name:    "boundary priority and multi-character separators",
			content: "aa\n\nbbbb\ncc\n\ndd",
			policy: ChunkingPolicies{
				ChunkSize:          4,
				BoundaryConditions: []string{"\n\n", "\n"},
			},
			wantRanges: [][2]int{{0, 1}, {4, 7}, {9, 10}, {13, 14}},
		},
		{
			name:    "unicode characters use rune positions",
			content: "🙂🙂|éé|中",
			policy: ChunkingPolicies{
				ChunkSize:          2,
				BoundaryConditions: []string{"|"},
			},
			wantRanges: [][2]int{{0, 1}, {3, 4}, {6, 6}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := SubdivideIntoChunks(docs.Document{
				DocumentUid: "test-document",
				Content:     test.content,
			}, test.policy)
			if err != nil {
				t.Fatalf("SubdivideIntoChunks() returned unexpected error: %v", err)
			}

			if gotRanges := chunkRanges(got); !reflect.DeepEqual(gotRanges, test.wantRanges) {
				t.Fatalf("chunk ranges = %v, want %v", gotRanges, test.wantRanges)
			}

			for _, chunk := range got {
				if chunk.DocumentUid != "test-document" {
					t.Errorf("chunk document UID = %q, want %q", chunk.DocumentUid, "test-document")
				}
			}
		})
	}
}

func TestSubdivideIntoChunksRejectsInvalidPolicies(t *testing.T) {
	tests := []struct {
		name   string
		policy ChunkingPolicies
	}{
		{
			name:   "zero chunk size",
			policy: ChunkingPolicies{ChunkSize: 0},
		},
		{
			name:   "negative chunk size",
			policy: ChunkingPolicies{ChunkSize: -1},
		},
		{
			name:   "overlap equal to chunk size",
			policy: ChunkingPolicies{ChunkSize: 4, Overlapping: 4},
		},
		{
			name:   "overlap larger than chunk size",
			policy: ChunkingPolicies{ChunkSize: 4, Overlapping: 5},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			chunks, err := SubdivideIntoChunks(docs.Document{Content: "abcdef"}, test.policy)
			if err == nil {
				t.Fatal("SubdivideIntoChunks() returned nil error for invalid policy")
			}
			if len(chunks) != 0 {
				t.Fatalf("chunks = %v, want no chunks", chunks)
			}
		})
	}
}
