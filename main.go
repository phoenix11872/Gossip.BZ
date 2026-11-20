package main

import (
	"log"
	"net/http"
	"os"

	"Gossipz.BZ/router"
)

func main() {
	// fmt.Print("Hello world printing some stuff up")
	// docs, err := docs.IngestFile("/home/phoenix/Work/openSource/Gossipz.BZ/agents.md")
	// if err != nil {
	// 	fmt.Print(err)
	// }
	// fmt.Print(docs)
	srv:= router.ServerRoutes()
	port:= os.Getenv("port");
	err:= http.ListenAndServe(":"+port, srv)
	if(err !=nil){
		log.Fatal(err);
	}
}
