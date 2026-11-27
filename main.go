package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"Gossipz.BZ/router"
	"github.com/redis/go-redis/v9"
	"google.golang.org/genai"
)

func main() {
	// fmt.Print("Hello world printing some stuff up")
	// docs, err := docs.IngestFile("/home/phoenix/Work/openSource/Gossipz.BZ/agents.md")
	// if err != nil {
	// 	fmt.Print(err)
	// }
	// fmt.Print(docs)
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	defer rdb.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to redis at %s: %v", redisAddr, err)
	}
	log.Println("Connected to Redis")
	geminiClient, err:= genai.NewClient(context.Background(),
		&genai.ClientConfig{
			APIKey: os.Getenv("GEMINI_API_KEY"),
			Backend: genai.BackendGeminiAPI,
		})
	if err !=nil{
		log.Fatalf("Could not connect to gemini api")
	}
	
	srv := router.ServerRoutes(rdb, geminiClient)
	port := os.Getenv("port")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, srv); err != nil {
		log.Fatal(err)
	}
}
