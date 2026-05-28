package tests

import (
	"testing"

	"github.com/swd-15/FHIR_Patient_API/internal/fhir"
)

var testBundleJSON = []byte(`{
  "resourceType": "Bundle",
  "type": "collection",
  "entry": [
    {
      "resource": {
        "resourceType": "Patient",
        "id": "p001",
        "name": [{"family": "山田", "given": ["太郎"]}],
        "gender": "male",
        "birthDate": "1980-04-15"
      }
    },
    {
      "resource": {
        "resourceType": "Patient",
        "id": "p002",
        "name": [{"family": "佐藤", "given": ["花子"]}],
        "gender": "female",
        "birthDate": "1995-11-30"
      }
    }
  ]
}`)

//Bundleが正しくパースされるか
func TestParseBundle(t *testing.T) {
	bundle, err := fhir.ParseBundle(testBundleJSON)
	if err != nil {
		t.Fatalf("ParseBundle failed: %v", err)
	}
	if bundle.ResourceType != "Bundle" {
		t.Errorf("expected resourceType=Bundle, got %s", bundle.ResourceType)
	}
	if len(bundle.Entry) != 2 {
		t.Errorf("expected 2 entries, got %d", len(bundle.Entry))
	}
}

//患者のみ正しく抽出されるか
func TestExtractPatients(t *testing.T) {
	bundle, _ := fhir.ParseBundle(testBundleJSON)
	patients := fhir.ExtractPatients(bundle)

	if len(patients) != 2 {
		t.Fatalf("expected 2 patients, got %d", len(patients))
	}
	if patients[0].ID != "p001" {
		t.Errorf("expected id=p001, got %s", patients[0].ID)
	}
	if patients[0].FullName != "山田 太郎" {
		t.Errorf("expected full_name=山田 太郎, got %s", patients[0].FullName)
	}
}

//存在する患者IDで取得できるか
func TestFindPatient_Found(t *testing.T) {
	bundle, _ := fhir.ParseBundle(testBundleJSON)
	patient, ok := fhir.FindPatient(bundle, "p001")

	if !ok {
		t.Fatal("expected patient p001 to be found")
	}
	if patient.BirthDate != "1980-04-15" {
		t.Errorf("expected birthDate=1980-04-15, got %s", patient.BirthDate)
	}
}

//存在しないIDでfalseが返るか
func TestFindPatient_NotFound(t *testing.T) {
	bundle, _ := fhir.ParseBundle(testBundleJSON)
	_, ok := fhir.FindPatient(bundle, "p999")

	if ok {
		t.Fatal("expected patient p999 NOT to be found")
	}
}
