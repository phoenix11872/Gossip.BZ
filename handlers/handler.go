package handlers

import (
	"net/http"

	"github.com/redis/go-redis/v9"
	"google.golang.org/genai"
)

type AppHandler struct{
	RDB *redis.Client
	Gemini *genai.Client
}

func NewHandler(rdb *redis.Client, gemini *genai.Client)*AppHandler{
	return &AppHandler{
		RDB: rdb,
		Gemini: gemini,
	}
}

