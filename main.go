package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/swd-15/FHIR_Patient_API/internal/handler"
	"github.com/swd-15/FHIR_Patient_API/internal/service"
)

func main() {
	port := getEnv("PORT", "8080")

	var svc *service.PatientService
	var err error

	// FHIR_MODEがserverの場合にFHIRサーバーから取得
	if getEnv("FHIR_MODE", "file") == "server" {
	baseURL := getEnv("FHIR_BASE_URL", "https://hapi.fhir.org/baseR4")
	countStr := getEnv("FHIR_PATIENT_COUNT", "5")
	count := 5
	fmt.Sscanf(countStr, "%d", &count)
	log.Printf("Fetching %d patients from FHIR server: %s", count, baseURL)
	svc, err = service.NewPatientServiceFromFHIRMultiple(baseURL, count)
	} else {
	bundlePath := getEnv("BUNDLE_PATH", "sample/bundle.json")
	log.Printf("Loading from file: %s", bundlePath)
	svc, err = service.NewPatientService(bundlePath)
	}

	if err != nil {
		log.Fatalf("Failed to initialize patient service: %v", err)
	}

	r := gin.Default()
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
