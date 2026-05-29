package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swd-15/FHIR_Patient_API/internal/service"
)

//HTTPハンドラをまとめた構造体
type PatientHandler struct {
	svc *service.PatientService
}

//PatientHandlerを生成する
func NewPatientHandler(svc *service.PatientService) *PatientHandler {
	return &PatientHandler{svc: svc}
}

//Gin Engineにルートを登録する
func (h *PatientHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", h.Health)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/patients", h.ListPatients)
		v1.GET("/patients/:id", h.GetPatient)
	}
}

//APIの起動確認
func (h *PatientHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

//全患者の一覧を返す
func (h *PatientHandler) ListPatients(c *gin.Context) {
	patients := h.svc.ListPatients()
	c.JSON(http.StatusOK, gin.H{
		"count":    len(patients),
		"patients": patients,
	})
}

//指定IDの患者詳細を返す
func (h *PatientHandler) GetPatient(c *gin.Context) {
	id := c.Param("id")
	patient, ok := h.svc.GetPatient(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found", "id": id})
		return
	}
	c.JSON(http.StatusOK, patient)
}
