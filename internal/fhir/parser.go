package fhir

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// JSONファイルからFHIR Bundleを読み込んでパースする
func LoadBundle(path string) (*Bundle, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("bundle file read error: %w", err)
	}
	return ParseBundle(data)
}

// JSONバイト列をBundleにパースする
func ParseBundle(data []byte) (*Bundle, error) {
	var bundle Bundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return nil, fmt.Errorf("bundle parse error: %w", err)
	}
	return &bundle, nil
}

// BundleからPatientリソースのみを取り出す
func ExtractPatients(bundle *Bundle) []PatientSummary {
	var patients []PatientSummary
	for _, entry := range bundle.Entry {
		r := entry.Resource
		if r.ResourceType != "Patient" {
			continue
		}
		patients = append(patients, PatientSummary{
			ID:        r.ID,
			FullName:  buildFullName(r),
			Gender:    r.Gender,
			BirthDate: r.BirthDate,
		})
	}
	return patients
}

// 指定IDのPatientを返す
func FindPatient(bundle *Bundle, id string) (*PatientSummary, bool) {
	for _, entry := range bundle.Entry {
		r := entry.Resource
		if r.ResourceType != "Patient" || r.ID != id {
			continue
		}
		p := &PatientSummary{
			ID:        r.ID,
			FullName:  buildFullName(r),
			Gender:    r.Gender,
			BirthDate: r.BirthDate,
		}
		return p, true
	}
	return nil, false
}

// 指定患者IDに紐づくConditionリソースを返す
func FindConditions(bundle *Bundle, patientID string) []ConditionResponse {
	var conditions []ConditionResponse
	ref := "Patient/" + patientID
	for _, entry := range bundle.Entry {
		r := entry.Resource
		if r.ResourceType != "Condition" {
			continue
		}
		if r.Subject == nil || r.Subject.Reference != ref {
			continue
		}
		display, code := extractCode(r.Code)
		clinicalStatus := ""
		if r.ClinicalStatus != nil && len(r.ClinicalStatus.Coding) > 0 {
			clinicalStatus = r.ClinicalStatus.Coding[0].Code
		}
		conditions = append(conditions, ConditionResponse{
			PatientID:      patientID,
			Display:        display,
			Code:           code,
			ClinicalStatus: clinicalStatus,
			OnsetDate:      r.OnsetDateTime,
		})
	}
	return conditions
}

// 指定患者IDに紐づくAllergyIntoleranceリソースを返す
func FindAllergies(bundle *Bundle, patientID string) []AllergyIntoleranceResponse {
	var allergies []AllergyIntoleranceResponse
	ref := "Patient/" + patientID
	for _, entry := range bundle.Entry {
		r := entry.Resource
		if r.ResourceType != "AllergyIntolerance" {
			continue
		}
		if r.Patient == nil || r.Patient.Reference != ref {
			continue
		}
		display, code := extractCode(r.Code)
		clinicalStatus := ""
		if r.ClinicalStatus != nil && len(r.ClinicalStatus.Coding) > 0 {
			clinicalStatus = r.ClinicalStatus.Coding[0].Code
		}
		category := ""
		if len(r.AllergyCategory) > 0 {
    		category = r.AllergyCategory[0]
		}
		allergies = append(allergies, AllergyIntoleranceResponse{
			PatientID:      patientID,
			Display:        display,
			Code:           code,
			ClinicalStatus: clinicalStatus,
			Criticality:    r.Criticality,
			Category:       category,
		})
	}
	return allergies
}

func FindInfections(bundle *Bundle, patientID string) []InfectionResponse {
	var infections []InfectionResponse
	ref := "Patient/" + patientID
	for _, entry := range bundle.Entry {
		r := entry.Resource
		if r.ResourceType != "Observation" {
			continue
		}
		if r.Subject == nil || r.Subject.Reference != ref {
			continue
		}
		// LOINCシステムの感染症コードのみ絞り込む
		if r.Code == nil || len(r.Code.Coding) == 0 {
			continue
		}
		if r.Code.Coding[0].System != "http://loinc.org" {
			continue
		}
		display, code := extractCode(r.Code)
		result := ""
		if r.ValueCodeableConcept != nil && len(r.ValueCodeableConcept.Coding) > 0 {
			result = r.ValueCodeableConcept.Coding[0].Display
		}
		infections = append(infections, InfectionResponse{
			PatientID:     patientID,
			Display:       display,
			Code:          code,
			Result:        result,
			EffectiveDate: r.Effective,
			Status:        r.Status,
		})
	}
	return infections
}

func FindMedications(bundle *Bundle, patientID string) []MedicationResponse {
	var medications []MedicationResponse
	ref := "Patient/" + patientID
	for _, entry := range bundle.Entry {
		r := entry.Resource
		if r.ResourceType != "MedicationRequest" {
			continue
		}
		if r.Subject == nil || r.Subject.Reference != ref {
			continue
		}
		display, code := extractCode(r.MedicationCodeableConcept)
		dosage := ""
		if len(r.DosageInstruction) > 0 {
			dosage = r.DosageInstruction[0].Text
		}
		medications = append(medications, MedicationResponse{
			PatientID: patientID,
			Display:   display,
			Code:      code,
			Status:    r.Status,
			Dosage:    dosage,
		})
	}
	return medications
}

// 指定患者IDに紐づくObservationリソースを返す
func FindObservations(bundle *Bundle, patientID string) []ObservationResponse {
	var observations []ObservationResponse
	ref := "Patient/" + patientID
	for _, entry := range bundle.Entry {
		r := entry.Resource
		if r.ResourceType != "Observation" {
			continue
		}
		if r.Subject == nil || r.Subject.Reference != ref {
			continue
		}
		display, code := extractCode(r.Code)
		var value float64
		unit := ""
		if r.ValueQuantity != nil {
			value = r.ValueQuantity.Value
			unit = r.ValueQuantity.Unit
		}
		observations = append(observations, ObservationResponse{
			PatientID:     patientID,
			Display:       display,
			Code:          code,
			Value:         value,
			Unit:          unit,
			EffectiveDate: r.Effective,
			Status:        r.Status,
		})
	}
	return observations
}


// displayとcodeを取り出す
func extractCode(cc *CodeableConcept) (display, code string) {
	if cc == nil {
		return "", ""
	}
	if len(cc.Coding) > 0 {
		display = cc.Coding[0].Display
		code = cc.Coding[0].Code
	}
	if display == "" {
		display = cc.Text
	}
	return display, code
}

// 姓名を結合して返す
func buildFullName(r Resource) string {
	if len(r.Name) == 0 {
		return ""
	}
	n := r.Name[0]
	given := strings.Join(n.Given, " ")
	if given != "" {
		return n.Family + " " + given
	}
	return n.Family
}
