package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/swd-15/FHIR_Patient_API/internal/handler"
	"github.com/swd-15/FHIR_Patient_API/internal/service"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	svc, _ := service.NewPatientService("../sample/bundle.json")
	h := handler.NewPatientHandler(svc)
	r := gin.Default()
	h.RegisterRoutes(r)
	return r
}

///healthが200を返すかチェック
func TestHealthEndpoint(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

///api/v1/patientsが200を返すかチェック
func TestListPatientsEndpoint(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/patients", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	if res["count"].(float64) == 0 {
		t.Error("expected at least 1 patient")
	}
}

//存在する患者IDで200を返すかチェック
func TestGetPatientEndpoint_Found(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/patients/p001", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

//存在しない患者IDで404を返すかチェック
func TestGetPatientEndpoint_NotFound(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/patients/p999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}
