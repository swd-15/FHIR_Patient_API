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
    },
    {
      "resource": {
        "resourceType": "Condition",
        "subject": {"reference": "Patient/p001"},
        "code": {
          "coding": [{"system": "http://hl7.org/fhir/sid/icd-10", "code": "E11", "display": "2型糖尿病"}]
        },
        "clinicalStatus": {
          "coding": [{"code": "active"}]
        },
        "onsetDateTime": "2018-06-01"
      }
    },
    {
      "resource": {
        "resourceType": "Observation",
        "status": "final",
        "subject": {"reference": "Patient/p001"},
        "code": {
          "coding": [{"system": "http://jpfhir.jp/fhir/core/CodeSystem/JP_ObservationLabResultCode", "code": "3H010000002326101", "display": "HbA1c"}]
        },
        "effectiveDateTime": "2024-11-10",
        "valueQuantity": {"value": 7.8, "unit": "%"}
      }
    },
    {
      "resource": {
        "resourceType": "Observation",
        "status": "final",
        "subject": {"reference": "Patient/p001"},
        "category": [
          {
            "coding": [
              {
                "system": "http://terminology.hl7.org/CodeSystem/observation-category",
                "code": "laboratory"
              }
            ]
          }
        ],
        "code": {
          "coding": [{"system": "http://loinc.org", "code": "22322-2", "display": "HBs抗原（B型肝炎）"}]
        },
        "effectiveDateTime": "2024-11-10",
        "valueCodeableConcept": {
          "coding": [{"system": "http://snomed.info/sct", "code": "260385009", "display": "陰性"}]
        }
      }
    },
    {
      "resource": {
        "resourceType": "MedicationRequest",
        "status": "active",
        "subject": {"reference": "Patient/p001"},
        "medicationCodeableConcept": {
          "coding": [
            {
              "system": "http://www.nlm.nih.gov/research/umls/rxnorm",
              "code": "197844",
              "display": "メトホルミン500mg"
            }
          ]
        },
        "dosageInstruction": [
          {"text": "500mg 1日2回"}
        ]
      }
    },
    {
      "resource": {
        "resourceType": "AllergyIntolerance",
        "patient": {"reference": "Patient/p001"},
        "allergyCategory": ["medication"],
        "code": {
          "coding": [{"system": "http://jpfhir.jp/fhir/core/CodeSystem/JP_AllergyIntolerance", "code": "J8A3199", "display": "ペニシリン"}]
        },
        "clinicalStatus": {
          "coding": [{"code": "active"}]
        },
        "criticality": "high"
      }
    }
  ]
}`)

// Bundleが正しくパースされるか
func TestParseBundle(t *testing.T) {
	bundle, err := fhir.ParseBundle(testBundleJSON)
	if err != nil {
		t.Fatalf("ParseBundle failed: %v", err)
	}
	if bundle.ResourceType != "Bundle" {
		t.Errorf("expected resourceType=Bundle, got %s", bundle.ResourceType)
	}
	if len(bundle.Entry) != 7 {
		t.Errorf("expected 7 entries, got %d", len(bundle.Entry))
	}
}

// 患者のみ正しく抽出されるか
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

// 存在する患者IDで取得できるか
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

// 存在しないIDでfalseが返るか
func TestFindPatient_NotFound(t *testing.T) {
	bundle, _ := fhir.ParseBundle(testBundleJSON)
	_, ok := fhir.FindPatient(bundle, "p999")
	if ok {
		t.Fatal("expected patient p999 NOT to be found")
	}
}

// 疾患情報が正しく取得できるか
func TestFindConditions(t *testing.T) {
	bundle, _ := fhir.ParseBundle(testBundleJSON)
	conditions := fhir.FindConditions(bundle, "p001")
	if len(conditions) != 1 {
		t.Fatalf("expected 1 condition, got %d", len(conditions))
	}
	if conditions[0].Display != "2型糖尿病" {
		t.Errorf("expected display=2型糖尿病, got %s", conditions[0].Display)
	}
	if conditions[0].Code != "E11" {
		t.Errorf("expected code=E11, got %s", conditions[0].Code)
	}
}

// 検査値が正しく取得できるか
func TestFindObservations(t *testing.T) {
	bundle, _ := fhir.ParseBundle(testBundleJSON)
	observations := fhir.FindObservations(bundle, "p001")
	if len(observations) != 2 {
		t.Fatalf("expected 2 observations, got %d", len(observations))
	}
	if observations[0].Display != "HbA1c" {
		t.Errorf("expected display=HbA1c, got %s", observations[0].Display)
	}
}

// 処方情報が正しく取得できるか
func TestFindMedications(t *testing.T) {
	bundle, _ := fhir.ParseBundle(testBundleJSON)
	medications := fhir.FindMedications(bundle, "p001")
	if len(medications) != 1 {
		t.Fatalf("expected 1 medication, got %d", len(medications))
	}
	if medications[0].Display != "メトホルミン500mg" {
		t.Errorf("expected display=メトホルミン500mg, got %s", medications[0].Display)
	}
	if medications[0].Dosage != "500mg 1日2回" {
		t.Errorf("expected dosage=500mg 1日2回, got %s", medications[0].Dosage)
	}
}

// アレルギー情報が正しく取得できるか
func TestFindAllergies(t *testing.T) {
	bundle, _ := fhir.ParseBundle(testBundleJSON)
	allergies := fhir.FindAllergies(bundle, "p001")
	if len(allergies) != 1 {
		t.Fatalf("expected 1 allergy, got %d", len(allergies))
	}
	if allergies[0].Display != "ペニシリン" {
		t.Errorf("expected display=ペニシリン, got %s", allergies[0].Display)
	}
	if allergies[0].Category != "medication" {
		t.Errorf("expected category=medication, got %s", allergies[0].Category)
	}
}

// 感染症情報が正しく取得できるか
func TestFindInfections(t *testing.T) {
	bundle, _ := fhir.ParseBundle(testBundleJSON)
	infections := fhir.FindInfections(bundle, "p001")
	if len(infections) != 1 {
		t.Fatalf("expected 1 infection, got %d", len(infections))
	}
	if infections[0].Display != "HBs抗原（B型肝炎）" {
		t.Errorf("expected display=HBs抗原（B型肝炎）, got %s", infections[0].Display)
	}
	if infections[0].Result != "陰性" {
		t.Errorf("expected result=陰性, got %s", infections[0].Result)
	}
}

// 別患者のデータが混入しないか
func TestFindConditions_OtherPatient(t *testing.T) {
	bundle, _ := fhir.ParseBundle(testBundleJSON)
	conditions := fhir.FindConditions(bundle, "p002")
	if len(conditions) != 0 {
		t.Errorf("expected 0 conditions for p002, got %d", len(conditions))
	}
}
