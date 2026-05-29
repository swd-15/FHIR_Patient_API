package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/swd-15/FHIR_Patient_API/internal/handler"
	"github.com/swd-15/FHIR_Patient_API/internal/service"
)

func main() {
	bundlePath := getEnv("BUNDLE_PATH", "sample/bundle.json")
	port := getEnv("PORT", "8080")

	// サービス初期化
	svc, err := service.NewPatientService(bundlePath)
	if err != nil {
		log.Fatalf("Failed to initialize patient service: %v", err)
	}

	// ルーター設定
	r := gin.Default()

	// ハンドラー登録
	h := handler.NewPatientHandler(svc)
	h.RegisterRoutes(r)

	log.Printf("Starting FHIR Patient API on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
