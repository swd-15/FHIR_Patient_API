package tests

import (
	"testing"

	"github.com/swd-15/FHIR_Patient_API/internal/service"
)

const testBundlePath = "../sample/bundle.json"

//サービスが正しく初期化されるかチェック
func TestNewPatientService(t *testing.T) {
	svc, err := service.NewPatientService(testBundlePath)
	if err != nil {
		t.Fatalf("NewPatientService failed: %v", err)
	}
	if svc == nil {
		t.Fatal("expected service to be non-nil")
	}
}

//患者一覧が取得できるかチェック
func TestListPatients(t *testing.T) {
	svc, _ := service.NewPatientService(testBundlePath)
	patients := svc.ListPatients()

	if len(patients) == 0 {
		t.Fatal("expected at least 1 patient")
	}
}

//存在する患者IDで取得できるかチェック
func TestGetPatient_Found(t *testing.T) {
	svc, _ := service.NewPatientService(testBundlePath)
	patient, ok := svc.GetPatient("p001")

	if !ok {
		t.Fatal("expected patient p001 to be found")
	}
	if patient.ID != "p001" {
		t.Errorf("expected id=p001, got %s", patient.ID)
	}
}

//存在しないIDでfalseが返るかチェック
func TestGetPatient_NotFound(t *testing.T) {
	svc, _ := service.NewPatientService(testBundlePath)
	_, ok := svc.GetPatient("p999")

	if ok {
		t.Fatal("expected patient p999 NOT to be found")
	}
}
