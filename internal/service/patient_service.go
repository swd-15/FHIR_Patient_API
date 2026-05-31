package service

import (
	"fmt"

	"github.com/swd-15/FHIR_Patient_API/internal/fhir"
)

//FHIR Bundleをインメモリに保持し、クエリを提供
type PatientService struct {
	bundle *fhir.Bundle
}

//BundleファイルをロードしてServiceを初期化
func NewPatientService(bundlePath string) (*PatientService, error) {
	bundle, err := fhir.LoadBundle(bundlePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load bundle: %w", err)
	}
	return &PatientService{bundle: bundle}, nil
}

//全患者のサマリーリストを返す
func (s *PatientService) ListPatients() []fhir.PatientSummary {
	return fhir.ExtractPatients(s.bundle)
}

//指定IDの患者を返す
func (s *PatientService) GetPatient(id string) (*fhir.PatientSummary, bool) {
	return fhir.FindPatient(s.bundle, id)
}

//指定患者の疾患情報リストを返す
func (s *PatientService) GetConditions(patientID string) []fhir.ConditionResponse {
	return fhir.FindConditions(s.bundle, patientID)
}

//指定患者のアレルギー情報リストを返す
func (s *PatientService) GetAllergies(PatientID string) []fhir.AllergyIntoleranceResponse{
	return fhir.FindAllergies(s.bundle, PatientID)
}

//指定患者の検査値リストを返す
func (s *PatientService) GetObservations(patientID string) []fhir.ObservationResponse {
	return fhir.FindObservations(s.bundle, patientID)
}
