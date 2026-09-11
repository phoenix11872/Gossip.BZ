package chunking

import (
	"fmt"
	"unicode/utf8"
)

type Chunks struct {
	DocumentUid string
	ParagraphId int
	Start       int
	End         int
	Embeddings []float64
}

type ChunkingPolicies struct{
	Overlapping int
	ChunkSize int
	strategy string
}
func SubDivideIntoChunksByParagraphs(doc Document, p ChunkingPolicies){
	ChunkList := []Chunks{}
	overlappingWindow:= []rune
	szOfDoc:= utf8.RuneCountInString(doc.Content)
	content:= []rune(doc.Content)
	for i:=0;i<szOfDoc;{
		j:=i
		szOverlap := utf8.RuneCountInString(overlappingWindow)
		textToConsider:= overlappingWindow
		newOverlappingWindow:= []rune
		prevSlashNCharacterPos :=-1
		for ;j<min(szOfDoc, p.ChunkSize - szOverlap);j++{
			if ParagraphSwitchDetected(content[j],prevSlashNCharacterPos){
				break;
			}
			textToConsider.
			
		}
	}
}

func SubdivideIntoChunks(doc Document, p ChunkingPolicies)([]Chunks, error){
	chunkslist := []Chunks{}
	if(p.Overlapping > p.ChunkSize){
		return chunkslist, fmt.Errorf("Overlapping token window cannot be greater than chunks size")
	}
}



