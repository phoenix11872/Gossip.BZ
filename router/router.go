package router

import (
	"net/http"
	"os"

	"Gossipz.BZ/handlers"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"google.golang.org/genai"
)

func ServerRoutes(rdb *redis.Client,gemini *genai.Client) http.Handler{
	mux:= http.NewServeMux();
	h:= handlers.NewHandler(rdb, gemini);
	// mux.HandleFunc("POST /{userId}/document",);
	// mux.HandleFunc("GET /api/v1/session");
	// mux.HandleFunc("GET /api/v1/session");
	// mux.HandleFunc("GET /api/v1/session/{sessionId}/message");
	// mux.HandleFunc("POST /api/v1/session/{sessionId}/message");
	
	clientUrl:= os.Getenv("clientUrl");
	c:= cors.New(cors.Options{
		AllowedOrigins: []string{clientUrl},
		AllowCredentials: true,
		AllowedMethods: []string{"GET","POST"},
		AllowedHeaders: []string{"Content-Type","Authorization"},
	})
	return c.Handler(mux);
}