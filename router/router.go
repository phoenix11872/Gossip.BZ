package router

import (
	"net/http"
	"os"

	"github.com/rs/cors"
)

func ServerRoutes() http.Handler{
	mux:= http.NewServeMux();
	mux.HandleFunc("POST /{userId}/document",);
	mux.HandleFunc("GET /api/v1/session");
	mux.HandleFunc("GET /api/v1/session");
	mux.HandleFunc("GET /api/v1/session/{sessionId}/message");
	mux.HandleFunc("POST /api/v1/session/{sessionId}/message");
	
	clientUrl:= os.Getenv("clientUrl");
	c:= cors.New(cors.Options{
		AllowedOrigins: []string{clientUrl},
		AllowCredentials: true,
		AllowedMethods: []string{"GET","POST"},
		AllowedHeaders: []string{"Content-Type","Authorization"},
	})
	return c.Handler(mux);
}