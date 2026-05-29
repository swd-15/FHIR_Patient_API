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

