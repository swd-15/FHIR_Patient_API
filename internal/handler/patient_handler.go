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
		v1.GET("/patients/:id/conditions", h.GetConditions)
		v1.GET("/patients/:id/allergies", h.GetAllergies)
		v1.GET("/patients/:id/observations", h.GetObservations)
		v1.GET("/patients/:id/medications", h.GetMedications)
		v1.GET("/patients/:id/infections", h.GetInfections)
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

//指定患者の疾患情報を返す
func (h *PatientHandler) GetConditions(c *gin.Context) {
	id := c.Param("id")
	if _, ok := h.svc.GetPatient(id); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found", "id": id})
		return
	}
	conditions := h.svc.GetConditions(id)
	c.JSON(http.StatusOK, gin.H{
		"patient_id": id,
		"count":      len(conditions),
		"conditions": conditions,
	})
}

//指定患者の処方情報を返す
func (h *PatientHandler) GetMedications(c *gin.Context) {
	id :=c.Param("id")
	if _, ok := h.svc.GetPatient(id); !ok{
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found","id": id})
		return
	}
	medications := h.svc.GetMedications(id)
	c.JSON(http.StatusOK, gin.H{
		"patient_id":  id,
		"count":       len(medications),
		"medications": medications,
	})
}

//指定患者のアレルギー情報を返す
func (h *PatientHandler) GetAllergies(c *gin.Context) {
	id := c.Param("id")
	if _, ok := h.svc.GetPatient(id); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found", "id": id})
		return
	}
	allergies := h.svc.GetAllergies(id)
	c.JSON(http.StatusOK, gin.H{
		"patient_id": id,
		"count":      len(allergies),
		"allergies":  allergies,
	})
}

//指定患者の検査値情報を返す
func (h *PatientHandler) GetObservations(c *gin.Context) {
	id := c.Param("id")
	if _, ok := h.svc.GetPatient(id); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found", "id": id})
		return
	}
	observations := h.svc.GetObservations(id)
	c.JSON(http.StatusOK, gin.H{
		"patient_id":   id,
		"count":        len(observations),
		"observations": observations,
	})
}

func (h *PatientHandler) GetInfections(c *gin.Context) {
	id :=c.Param("id")
	if _, ok := h.svc.GetPatient(id); !ok{
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found","id": id})
		return
	}
	infections := h.svc.GetInfections(id)
	c.JSON(http.StatusOK, gin.H{
		"patient_id":  id,
		"count":       len(infections),
		"infections": infections,
	})
}
