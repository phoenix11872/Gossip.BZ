package main

import (
	"fmt"

	chunking "Gossipz.BZ/internal"
)

func main() {
	fmt.Print("Hello world printing some stuff up")
	docs, err := chunking.IngestFile("/home/phoenix/Work/openSource/Gossipz.BZ/agents.md")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(docs)
}
